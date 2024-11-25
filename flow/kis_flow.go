package flow

import (
	"KisFlow/common"
	"KisFlow/config"
	"KisFlow/conn"
	"KisFlow/function"
	"KisFlow/id"
	"KisFlow/kis"
	"KisFlow/log"
	"context"
	"errors"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

type KisFlow struct {
	// Flow的分布式实例ID(用于KisFlow内部区分不同实例)
	Id string
	// Flow的可读名称
	Name string
	// FLow配置策略
	Conf *config.KisFlowConfig
	//Function 列表
	Funcs          map[string]kis.Function // 当前flow拥有的全部管理的全部Function对象, key: FunctionName
	FlowHead       kis.Function            // 当前Flow所拥有的Function列表表头
	FlowTail       kis.Function            // 当前Flow所拥有的Function列表表尾
	flock          sync.RWMutex            // 管理链表插入读写的锁
	ThisFunction   kis.Function            // Flow当前正在执行的KisFunction对象
	ThisFunctionId string                  // 当前执行到的Function ID (策略配置ID)
	PrevFunctionId string                  // 当前执行到的Function 上一层FunctionID(策略配置ID)
	// Function 列表参数
	funcParams map[string]config.FParam // flow在当前Function的自定义固定配置参数,Key:function的实例NsID, value:FParam

	fplock sync.RWMutex      // 管理funcParams的读写锁
	buffer common.KisRowArr  // 用来临时存放输入字节数据的内部Buf, 一条数据为interface{}, 多条数据为[]interface{} 也就是KisBatch
	data   common.KisDataMap // 流式计算各个层级的数据源
	inPut  common.KisRowArr  // 当前Function的计算输入数据

	action kis.Action // 当前Flow所携带的Action动作
	abort  bool       // 是否中断Flow

	cache *cache.Cache //Flow流的临时缓存上下文环境

	metaData map[string]interface{} // Flow自订的临时数据
	mLock    sync.RWMutex           //管理metaData的读写锁
}

// TODO for test
// NewKisFlow 创建一个KisFlow
func NewKisFlow(conf *config.KisFlowConfig) *KisFlow {
	flow := new(KisFlow)

	// 基础信息
	flow.Id = id.KisID(common.KisIdTypeFlow)
	flow.Name = conf.FlowName
	flow.Conf = conf

	// Function 列表
	flow.Funcs = make(map[string]kis.Function)
	flow.funcParams = make(map[string]config.FParam)

	flow.data = make(common.KisDataMap)

	flow.cache = cache.New(cache.NoExpiration, common.DeFaultFlowCacheCleanUp*time.Minute)

	flow.metaData = make(map[string]interface{})

	return flow
}

// Link 将Function链接到Flow中
// fConf: 当前Function策略
// fParams: 当前Flow携带的Function动态参数
func (flow *KisFlow) Link(fConf *config.KisFuncConfig, fParams config.FParam) error {
	// 创建Function
	f := function.NewKisFunction(flow, fConf)

	// 当前Function有Connector关联，需要初始化Connector实例
	if fConf.Option.CName != "" {

		// 获取Connector配置
		connConfig, err := fConf.GetConnConfig()
		if err != nil {
			panic(err)
		}

		// 创建Connector对象
		connector := conn.NewKisConnector(connConfig)

		// 初始化Connector, 执行ConnectorInit方法
		if err = connector.Init(); err != nil {
			panic(err)
		}

		// 关联Function实例和Connector实例关系
		_ = f.AddConnector(connector)
	}

	// Flow 添加 Function
	if err := flow.appendFunc(f, fParams); err != nil {
		return err
	}

	return nil
}

// appendFunc 将Function添加到Flow中, 链表操作
func (flow *KisFlow) appendFunc(function kis.Function, fParam config.FParam) error {

	if function == nil {
		return errors.New("AppendFunc append nil to List")
	}

	flow.flock.Lock()
	defer flow.flock.Unlock()

	if flow.FlowHead == nil {
		// 首次添加节点
		flow.FlowHead = function
		flow.FlowTail = function

		function.SetN(nil)
		function.SetP(nil)
	} else {
		// 将Function插入链表尾部
		function.SetP(flow.FlowTail)
		function.SetN(nil)

		flow.FlowTail.SetN(function)
		flow.FlowTail = function
	}

	// 将Function Name 详细Hash对应关系添加到flow对象中
	flow.Funcs[function.GetConfig().FName] = function

	// 先添加function默认携带的Params参数
	params := make(config.FParam)

	for key, value := range function.GetConfig().Option.Params {
		params[key] = value
	}

	// 将得到的FParams存留在flow结构体中，用来function业务直接通过Hash获取
	// key 为当前Function的KisId，不用Fid的原因是为了防止一个Flow添加两个相同策略Id的Function

	flow.funcParams[function.GetId()] = params

	return nil
}

// Run 启动KisFlow的流式计算, 从起始Function开始执行流
func (flow *KisFlow) Run(ctx context.Context) error {

	var fn kis.Function

	fn = flow.FlowHead
	// 重置abort
	flow.abort = false // 每次进入调度，要重置abort状态

	if flow.Conf.Status == int(common.FlowDisable) {
		// flow被配置关闭
		return nil
	}

	// 因为此时还没有执行任何Function, 所以PrevFunctionId为FirstVirtual 因为没有上一层Function
	flow.PrevFunctionId = common.FunctionIdFirstVirtual

	// 提交数据流原始数据
	if err := flow.commitSrcData(ctx); err != nil {
		return err
	}

	// 流式链式调用
	// 如果设置abort则不进入下次循环调度
	for fn != nil && flow.abort == false {
		// flow 当前执行到的Function标记
		fid := fn.GetId()
		flow.ThisFunction = fn
		flow.ThisFunctionId = fid

		// 得到当前Function 要处理的源数据
		if inputData, err := flow.getCurData(); err != nil {
			log.Logger().ErrorFX(ctx, "flow.Run(): getCurData err = %s\n", err.Error())
			return err
		} else {
			flow.inPut = inputData
		}

		if err := fn.Call(ctx, flow); err != nil {
			// Error
			return err
		} else {
			// Success
			fn, err = flow.dealAction(ctx, fn)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (flow *KisFlow) GetName() string {

	return flow.Name

}

func (flow *KisFlow) GetThisFunction() kis.Function {

	return flow.ThisFunction

}

func (flow *KisFlow) GetThisFuncConf() *config.KisFuncConfig {
	return flow.ThisFunction.GetConfig()
}

// GetConnector 得到当前正在执行的Function的Connector
func (flow *KisFlow) GetConnector() (kis.Connector, error) {
	if conn := flow.ThisFunction.GetConnector(); conn != nil {
		return conn, nil
	} else {
		return nil, errors.New("GetConnector(): Connector is nil")
	}
}

// GetConnConf 得到当前正在执行的Function的Connector的配置
func (flow *KisFlow) GetConnConf() (*config.KisConnConfig, error) {
	if conn := flow.ThisFunction.GetConnector(); conn != nil {
		return conn.GetConfig(), nil
	} else {
		return nil, errors.New("GetConnConf(): Connector is nil")
	}
}

func (flow *KisFlow) GetConfig() *config.KisFlowConfig {
	return flow.Conf
}

func (flow *KisFlow) GetFuncConfigByName(funcName string) *config.KisFuncConfig {
	if f, ok := flow.Funcs[funcName]; ok {
		return f.GetConfig()
	} else {
		log.Logger().ErrorF("GetFuncConfigByName(): Function %s not found", funcName)
		return nil
	}
}

// Next 当前Flow执行到的Function进入下一层Function所携带的Action动作
func (flow *KisFlow) Next(acts ...kis.ActionFunc) error {

	// 加载Function FaaS传递的Action动作
	flow.action = kis.LoadActions(acts)

	return nil
}
