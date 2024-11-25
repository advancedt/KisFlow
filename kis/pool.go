package kis

import (
	"KisFlow/common"
	"KisFlow/log"
	"context"
	"errors"
	"fmt"
	"sync"
)

// 同步控制变量，确保代码段仅执行一次
// 提供一个Do 方法
var _poolOnce sync.Once

// kisPool 用于管理全部的Function和Flow配置的池子
type kisPool struct {
	fnRouter funcRouter   // 全部function的管理路由
	fnLock   sync.RWMutex // fnRouter锁

	flowRouter flowRouter   // 全部的flow对象
	flowLock   sync.RWMutex // flowRouter 锁

	cInitRouter ConnInitRouter // 全部的Connector初始化路由
	ciLock      sync.RWMutex   //cInitRouter 锁

	cTree      connTree             // 全部的Connector管理路由
	connectors map[string]Connector // 全部的connector对象
	clock      sync.RWMutex         // cTree锁
}

// 单例
var _pool *kisPool

// Pool 单例构造
func Pool() *kisPool {
	_poolOnce.Do(func() {
		// 创建kisPool对象
		_pool = new(kisPool)

		//fnRouter 初始化
		_pool.fnRouter = make(funcRouter)

		//flowRouter 初始化
		_pool.flowRouter = make(flowRouter)

		// connTree初始化
		_pool.cTree = make(connTree)
		_pool.cInitRouter = make(ConnInitRouter)
		_pool.connectors = make(map[string]Connector)
	})

	return _pool
}

// 注册以及获取Pool
func (pool *kisPool) AddFlow(name string, flow Flow) {
	pool.flowLock.Lock()
	defer pool.flowLock.Unlock()

	// 对Flow进行查验，相同的Flow无法注册多次
	if _, ok := pool.flowRouter[name]; !ok {
		pool.flowRouter[name] = flow
	} else {
		errString := fmt.Sprintf("Pool AddFlow Repeat FlowName=%s\n", name)
		panic(errString)
	}

	log.Logger().InfoF("Add FlowRouter FlowName=%s\n", name)
}

func (pool *kisPool) GetFlow(name string) Flow {
	pool.flowLock.RLock()
	defer pool.flowLock.RUnlock()

	if flow, ok := pool.flowRouter[name]; ok {
		return flow
	} else {
		return nil
	}
}

// 注册及调度Function

// FaaS注册Function计算业务逻辑，通过Function Name索引及注册
func (pool *kisPool) FaaS(fnName string, f FaaS) {
	pool.fnLock.Lock()
	defer pool.fnLock.Unlock()

	if _, ok := pool.fnRouter[fnName]; !ok {
		pool.fnRouter[fnName] = f
	} else {
		errString := fmt.Sprintf("KisPoll FaaS Repeat FuncName=%s", fnName)
		panic(errString)
	}

	log.Logger().InfoF("Add KisPool FuncName=%s", fnName)
}

// CallFunction 调度 Function
func (pool *kisPool) CallFunction(ctx context.Context, fnName string, flow Flow) error {
	if f, ok := pool.fnRouter[fnName]; ok {
		return f(ctx, flow)
	}
	log.Logger().ErrorFX(ctx, "FuncName: %s Can not find in KisPool, Not Added.\n", fnName)

	return errors.New("FuncName: " + fnName + " Can not find in NsPool, Not Added.")
}

// CaaSInit对象
func (pool *kisPool) CasSInit(cname string, c ConnInit) {
	pool.ciLock.Lock()
	defer pool.ciLock.Unlock()

	if _, ok := pool.cInitRouter[cname]; !ok {
		pool.cInitRouter[cname] = c
	} else {
		errString := fmt.Sprintf("KisPool Reg CaaSInit Repeat CName=%s\n", cname)
		panic(errString)
	}

	log.Logger().InfoF("Add KisPool CaaSInit CName=%s", cname)
}

func (pool *kisPool) CallConnInit(conn Connector) error {
	pool.ciLock.RLock()
	defer pool.ciLock.RUnlock()

	init, ok := pool.cInitRouter[conn.GetName()]

	if !ok {
		panic(errors.New(fmt.Sprintf("init connector cname = %s not reg..", conn.GetName())))
	}

	return init(conn)
}

// CaaS注册Connector Call 业务

func (pool *kisPool) CaaS(cname string, fname string, mode common.KisMode, c CaaS) {
	pool.clock.Lock()
	defer pool.clock.Unlock()

	if _, ok := pool.cTree[cname]; !ok {
		// cid首次注册
		pool.cTree[cname] = make(connSL)

		// 初始化各类型FunctionMode
		pool.cTree[cname][common.S] = make(connFuncRouter)
		pool.cTree[cname][common.L] = make(connFuncRouter)
	}

	if _, ok := pool.cTree[cname][mode][fname]; !ok {
		pool.cTree[cname][mode][fname] = c
	} else {
		errString := fmt.Sprintf("CaaS Repeat CName=%s, FName=%s, Mode =%s\n", cname, fname, mode)
		panic(errString)
	}
	log.Logger().InfoF("Add KisPool CaaS CName=%s, FName=%s, Mode =%s", cname, fname, mode)
}

// CallConnector 调度Connector
func (pool *kisPool) CallConnector(ctx context.Context, flow Flow, conn Connector, args interface{}) error {
	fn := flow.GetThisFunction()
	fnConf := flow.GetThisFuncConf()
	mode := common.KisMode(fnConf.FMode)

	if callback, ok := pool.cTree[conn.GetName()][mode][fnConf.FName]; ok {
		return callback(ctx, conn, fn, flow, args)
	}

	log.Logger().ErrorFX(ctx, "CName:%s FName:%s mode:%s Can not find in KisPool, Not Added.\n", conn.GetName(), fnConf.FName, mode)

	return errors.New(fmt.Sprintf("CName:%s FName:%s mode:%s Can not find in KisPool, Not Added.", conn.GetName(), fnConf.FName, mode))
}

func (pool *kisPool) GetFlows() []Flow {
	pool.flowLock.RLock()
	defer pool.flowLock.RUnlock()

	var flows []Flow

	for _, flow := range pool.flowRouter {
		flows = append(flows, flow)
	}

	return flows
}
