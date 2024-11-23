package caas

import (
	"KisFlow/kis"
	"fmt"
)

func InitConnDemo1(connector kis.Connector) error {
	fmt.Println("===> Call Connector InitDemo1")

	connConf := connector.GetConfig()
	fmt.Println(connConf)

	// init connector , 如 初始化数据库链接等

	return nil
}
