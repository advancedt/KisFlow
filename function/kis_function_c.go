package function

import (
	"KisFlow/kis"
	"KisFlow/log"
	"context"
)

type KisFunctionC struct {
	BaseFunction
}

func (f *KisFunctionC) Call(ctx context.Context, flow kis.Flow) error {

	log.Logger().InfoF("KisFunctionC, flow = %+v\n", flow)

	// 通过KisPool 路由到具体的执行计算Funtion
	if err := kis.Pool().CallFunction(ctx, f.Config.FName, flow); err != nil {
		log.Logger().ErrorFX(ctx, "Function Called Error err = %s\n", err)
		return err
	}

	return nil
}
