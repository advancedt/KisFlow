package conn

import (
	"KisFlow/config"
	"KisFlow/id"
	"KisFlow/kis"
	"context"
	"sync"
)

type KisConnector struct {
	// Connector ID
	CID string
	// Connector Name
	CName string
	// Connector Config
	Conf *config.KisConnConfig

	//Connector Init
	onceInit sync.Once
}

// NewKisConnector 根据配置策略新创建一个KisConnector
func NewKisConnector(config *config.KisConnConfig) *KisConnector {
	conn := new(KisConnector)
	conn.CID = id.KisID()
	conn.CName = config.CName
	conn.Conf = config

	return conn
}

// 实现Connector接口

func (conn *KisConnector) Init() error {
	var err error

	// 一个Connector接口只能初始化一次
	conn.onceInit.Do(func() {
		err = kis.Pool().CallConnInit(conn)
	})
	return err
}

// Call 调用Connector外挂存储逻辑的读写操作
func (conn *KisConnector) Call(ctx context.Context, flow kis.Flow, args interface{}) error {
	// 通过KisPool进行调度
	if err := kis.Pool().CallConnector(ctx, flow, conn, args); err != nil {
		return err
	}
	return nil
}

func (conn *KisConnector) GetId() string {
	return conn.CID
}

func (conn *KisConnector) GetName() string {
	return conn.CName
}

func (conn *KisConnector) GetConfig() *config.KisConnConfig {
	return conn.Conf
}
