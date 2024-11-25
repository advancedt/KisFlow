package faas

import (
	"KisFlow/kis"
	"context"
	"fmt"
)

func JumpFuncHandler(ctx context.Context, flow kis.Flow) error {
	fmt.Println("---> Call JumpFuncHandler ----")

	for _, row := range flow.Input() {
		str := fmt.Sprintf("In FuncName = %s, FuncId = %s, row = %s", flow.GetThisFuncConf().FName, flow.GetThisFunction().GetId(), row)
		fmt.Println(str)
	}

	return flow.Next(kis.ActionJumpFunc("funcName1"))
}
