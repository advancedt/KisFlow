package function

import (
	"KisFlow/kis"
	"KisFlow/log"
	"context"
	"fmt"
)

type KisFunctionE struct {
	BaseFunction
}

func (f *KisFunctionE) Call(ctx context.Context, flow kis.Flow) error {
	//fmt.Printf("KisFunctionE, flow = %+v\n", flow)

	log.Logger().InfoF("KisFunctionE, flow = %+v\n", flow)

	for _, row := range flow.Input() {
		fmt.Printf("In KisFunctionE, row = %+v\n", row)
	}

	return nil
}
