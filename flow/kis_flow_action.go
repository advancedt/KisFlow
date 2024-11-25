package flow

import (
	"KisFlow/kis"
	"context"
)

// dealAction 处理Action,决定接下来Flow的流程走向
func (flow *KisFlow) dealAction(ctx context.Context, fn kis.Function) (kis.Function, error) {
	//if err := flow.commitCurData(ctx); err != nil {
	//	return nil, err
	//}
	if flow.action.DataReuse {
		if err := flow.commitReuseData(ctx); err != nil {
			return nil, err
		}
	} else {
		if err := flow.commitCurData(ctx); err != nil {
			return nil, err
		}
	}

	// 更新上一层FunctionId游标
	flow.PrevFunctionId = flow.ThisFunctionId
	fn = fn.Next()

	// Abort Action 强制中止
	if flow.action.Abort {
		flow.abort = true
	}

	// 清空Action
	flow.action = kis.Action{}

	return fn, nil
}
