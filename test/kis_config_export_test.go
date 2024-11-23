package test

import (
	"KisFlow/common"
	"KisFlow/file"
	"KisFlow/kis"
	"KisFlow/test/caas"
	"KisFlow/test/faas"
	"testing"
)

func TestConfigExportYmal(t *testing.T) {
	// 0. 注册Function 回调业务
	kis.Pool().FaaS("funcName1", faas.FuncDemo1Handler)
	kis.Pool().FaaS("funcName2", faas.FuncDemo2Handler)
	kis.Pool().FaaS("funcName3", faas.FuncDemo3Handler)

	// 0. 注册ConnectorInit 和 Connector 回调业务
	kis.Pool().CasSInit("ConnName1", caas.InitConnDemo1)
	kis.Pool().CaaS("ConnName1", "funcName2", common.S, caas.CaasDemoHandler1)

	// 1. 加载配置文件并构建Flow
	if err := file.ConfigImportYaml("D:\\GoLandProject\\KisFlow\\test\\load_conf"); err != nil {
		panic(err)
	}

	// 2. 讲构建的内存KisFlow结构配置导出的文件当中
	flows := kis.Pool().GetFlows()
	for _, flow := range flows {
		if err := file.ConfigExportYaml(flow, "D:\\GoLandProject\\KisFlow\\test\\export_conf"); err != nil {
			panic(err)
		}
	}
}
