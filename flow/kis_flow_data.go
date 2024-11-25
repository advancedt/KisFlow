package flow

import (
	"KisFlow/common"
	"KisFlow/log"
	"context"
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

func (flow *KisFlow) CommitRow(row interface{}) error {
	// 所有提交的数据都会暂存在flow.Buffer成员中，作为缓冲区
	flow.buffer = append(flow.buffer, row)

	return nil
}

// commitSrcData 提交当前Flow的数据源数据, 表示首次提交当前Flow的原始数据源
// 将flow的临时数据buffer，提交到flow的data中,(data为各个Function层级的源数据备份)
// 会清空之前所有的flow数据

func (flow *KisFlow) commitSrcData(ctx context.Context) error {
	/*
		在整个Flow的运行周期只会运行一次，作为当前Flow的原始数据
	*/

	// 制作批量数据batch
	dataCnt := len(flow.buffer)
	bacth := make(common.KisRowArr, 0, dataCnt)

	for _, row := range flow.buffer {
		bacth = append(bacth, row)
	}
	// 清空之前所有数据
	flow.clearData(flow.data)

	// 首次提交，记录flow原始数据
	flow.data[common.FunctionIdFirstVirtual] = bacth

	// 清空缓冲Buf
	flow.buffer = flow.buffer[0:0]

	log.Logger().DebugFX(ctx, "====> After CommitSrcData, flow_name = %s, flow_id = %s\nAll Level Data =\n %+v\n", flow.Name, flow.Id, flow.data)

	return nil
}

// 清空flow所有数据
func (flow *KisFlow) clearData(data common.KisDataMap) {
	for k := range data {
		delete(data, k)
	}
}

// 中间层数据的提交
func (flow KisFlow) commitCurData(ctx context.Context) error {

	// 判断本层计算是否有结果数据，如果没有则退出本次Flow Run循环
	if len(flow.buffer) == 0 {

		flow.abort = true
		return nil
	}

	// 批量制作batch
	batch := make(common.KisRowArr, 0, len(flow.buffer))

	// //如果strBuf为空，则没有添加任何数据
	for _, row := range flow.buffer {
		batch = append(batch, row)
	}

	//将本层计算的缓冲数据提交到本层的计算结果中
	flow.data[flow.ThisFunctionId] = batch

	// 清空缓冲Buf
	flow.buffer = flow.buffer[0:0]

	log.Logger().DebugFX(ctx, " ====> After commitCurData, flow_name = %s, flow_id = %s\nAll Level Data =\n %+v\n", flow.Name, flow.Id, flow.data)

	return nil
}

// 获取正在执行的Function的源数据
func (flow *KisFlow) getCurData() (arr common.KisRowArr, err error) {
	if flow.PrevFunctionId == "" {
		return nil, errors.New(fmt.Sprintf("flow.PrevFunctionId is not set"))
	}

	if _, ok := flow.data[flow.PrevFunctionId]; !ok {
		return nil, errors.New(fmt.Sprintf("[%s] is not in flow.data", flow.PrevFunctionId))
	}

	return flow.data[flow.PrevFunctionId], nil
}

// Input 得到flow当前执行Function的输入源数据
func (flow *KisFlow) Input() common.KisRowArr {
	return flow.inPut
}

func (flow *KisFlow) commitReuseData(ctx context.Context) error {
	// 判断上层是否有数据，如果没有则退出本次Flow Run循环
	if len(flow.data[flow.PrevFunctionId]) == 0 {
		flow.abort = true
		return nil
	}

	// 本层结果数据等于上层结果数据 (复用上层结果数据到本层)
	flow.data[flow.ThisFunctionId] = flow.data[flow.PrevFunctionId]

	// 清空缓冲Buf  (如果是ReuseData选项，那么提交的全部数据，都将不会携带到下一层)
	flow.buffer = flow.buffer[0:0]

	log.Logger().DebugFX(ctx, " ====> After commitReuseData, flow_name = %s, flow_id = %s\nAll Level Data =\n %+v\n", flow.Name, flow.Id, flow.data)

	return nil
}

func (flow *KisFlow) commitVoidData(ctx context.Context) error {

	if len(flow.buffer) != 0 {
		return nil
	}

	// 制作空数据
	batch := make(common.KisRowArr, 0)

	// 将本层计算的缓冲数据提交到本层结果数据中
	flow.data[flow.ThisFunctionId] = batch

	log.Logger().DebugFX(ctx, " ====> After commitVoidData, flow_name = %s, flow_id = %s\nAll Level Data =\n %+v\n", flow.Name, flow.Id, flow.data)

	return nil
}

func (flow *KisFlow) GetCacheData(key string) interface{} {
	if data, found := flow.cache.Get(key); found {
		return data
	}
	return nil
}

func (flow *KisFlow) SetCacheData(key string, value interface{}, Exp time.Duration) {
	if Exp == common.DefaultExpiration {
		flow.cache.Set(key, value, cache.DefaultExpiration)
	} else {
		flow.cache.Set(key, value, Exp)
	}
}

// GetMetaData 得到当前Flow对象的临时数据
func (flow *KisFlow) GetMetaData(key string) interface{} {
	flow.mLock.RLock()
	defer flow.mLock.RUnlock()

	data, ok := flow.metaData[key]
	if !ok {
		return nil
	}
	return data
}

func (flow *KisFlow) SetMetaData(key string, value interface{}) {
	flow.mLock.Lock()
	defer flow.mLock.Unlock()

	flow.metaData[key] = value
}
