package kis

import (
	"KisFlow/common"
	"KisFlow/config"
	"context"
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
}
