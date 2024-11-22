package function

import (
	"KisFlow/kis"
	"context"
	"fmt"
)

type KisFunctionS struct {
	BaseFunction
}

func (f *KisFunctionS) Call(ctx context.Context, flow kis.Flow) error {
	fmt.Printf("KisFunctionS, flow = %+v\n", flow)

	// TODO 调用具体的Function执行方法

	return nil
}
