package caas

import (
	"KisFlow/kis"
	"context"
	"fmt"
)

func CaasDemoHandler1(ctx context.Context, conn kis.Connector, fn kis.Function, flow kis.Flow, args interface{}) error {
	fmt.Printf("===> In CaasDemoHanler1: flowName: %s, cName:%s, fnName:%s, mode:%s\n",
		flow.GetName(), conn.GetName(), fn.GetConfig().FName, fn.GetConfig().FMode)
	fmt.Printf("===> In CaasDemoHanler1: flowName: %s, cName:%s, fnName:%s, mode:%s\n",
		flow.GetName(), conn.GetName(), fn.GetConfig().FName, fn.GetConfig().FMode)
	return nil
}