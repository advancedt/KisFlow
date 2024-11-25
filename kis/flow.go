package kis

import (
	"KisFlow/common"
	"KisFlow/config"
	"context"
	"time"
)

type Flow interface {
	// Run 调度Flow，一次调度Flow中的Function并且执行
	Run(ctx context.Context) error
	// Link 将 Flow 中的 Function 按照配置文件中的配置进行连接
	Link(fConf *config.KisFuncConfig, fParams config.FParam) error
	// 提交Flow数据到即将执行的Function层
	CommitRow(row interface{}) error
	// Input 得到flow当前执行Function的输入源数据
	Input() common.KisRowArr

	// 得到Flow名称
	GetName() string

	// GetThisFunction得到当前正在执行的Function
	GetThisFunction() Function
	// GetThisFunctionConf 得到当前正在执行的Function的配置
	GetThisFuncConf() *config.KisFuncConfig

	GetConnector() (Connector, error)
	// GetConnConf 得到当前正在执行的Function的Connector的配置
	GetConnConf() (*config.KisConnConfig, error)

	// GetConfig 得到当前Flow的配置
	GetConfig() *config.KisFlowConfig
	// GetFuncConfigByName 得到当前Flow的配置
	GetFuncConfigByName(funcName string) *config.KisFuncConfig
	// Next 当前Flow执行到的Funtion进入下一层Function锁携带的Action动作
	Next(acts ...ActionFunc) error
	// 得到当前Flow的缓存数据
	GetCacheData(key string) interface{}
	// 设置当前Flow的缓存数据
	SetCacheData(key string, value interface{}, Exp time.Duration)
	//GetMetaData得到当前Flow的临时数据
	GetMetaData(key string) interface{}
	//SetMetaData 设置当前Flow的临时数据
	SetMetaData(key string, value interface{})
	// GetFuncParam 得到Flow的当前正在执行的Function的配置默认参数，取出一对key-value
	GetFuncParam(key string) string
	// GetFuncParamAll 得到Flow的当前正在执行的Function的配置默认参数，取出全部Key-Value
	GetFuncParamAll() config.FParam
}
