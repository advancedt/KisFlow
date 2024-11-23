package faas

import (
	"KisFlow/kis"
	"context"
	"fmt"
)

func FuncDemo3Handler(ctx context.Context, flow kis.Flow) error {

	fmt.Println("---> Call funcName3Handler ----")

	for _, row := range flow.Input() {
		str := fmt.Sprintf("In FuncName = %s, FuncId = %s, row = %s", flow.GetThisFuncConf().FName, flow.GetThisFunction().GetId(), row)
		fmt.Println(str)
	}

	return nil
}
