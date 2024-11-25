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

	// KisConnector的自定义临时数据
	metaData map[string]interface{}
	// 管理metaData的读写锁
	mLock sync.RWMutex
}

// NewKisConnector 根据配置策略新创建一个KisConnector
func NewKisConnector(config *config.KisConnConfig) *KisConnector {
	conn := new(KisConnector)
	conn.CID = id.KisID()
	conn.CName = config.CName
	conn.Conf = config

	conn.metaData = make(map[string]interface{})

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

// GetMetaData 得到当前Connector的临时数据
func (conn *KisConnector) GetMetaData(key string) interface{} {
	conn.mLock.RLock()
	defer conn.mLock.RUnlock()

	data, ok := conn.metaData[key]
	if !ok {
		return nil
	}

	return data
}

// SetMetaData 设置当前Connector的临时数据
func (conn *KisConnector) SetMetaData(key string, value interface{}) {
	conn.mLock.Lock()
	defer conn.mLock.Unlock()

	conn.metaData[key] = value
}
