package kis

// Action KisFlow执行流程Actions
type Action struct {

	// DataReuse 是否复用上层
	DataReuse bool
	// Abort中指Flow执行
	Abort bool
}

// ActionFunc kisFlow Functional Option 类型
type ActionFunc func(ops *Action)

// LoadActions 加载Action，依执行ActionFunc操作函数
func LoadActions(acts []ActionFunc) Action {
	action := Action{}
	if acts == nil {
		return action
	}

	for _, act := range acts {
		act(&action)
	}
	return action
}

// ActionAbort 中止Flow的执行
func ActionAbort(action *Action) {
	action.Abort = true
}

// ActionDataReuse Next复用上层Function数据Option
func ActionDataReuse(act *Action) {
	act.DataReuse = true
}
