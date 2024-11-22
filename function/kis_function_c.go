package function

import (
	"KisFlow/kis"
	"KisFlow/log"
	"context"
	"fmt"
)

type KisFunctionC struct {
	BaseFunction
}

func (f *KisFunctionC) Call(ctx context.Context, flow kis.Flow) error {

	log.Logger().InfoF("KisFunctionC, flow = %+v\n", flow)

	// 处理业务数据

	for i, row := range flow.Input() {
		fmt.Printf("In KisFunctionC, row = %+v\n", row)

		// 提交本层计算结果数据
		_ = flow.CommitRow("Data From KisFunctionC, index " + " " + fmt.Sprintf("%d", i))
	}
	return nil
}
