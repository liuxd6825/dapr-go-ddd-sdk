package idao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
)

type DaoConfig struct {
	DbKey              string          `json:"dbKey"`              // 数据库Key
	IsPubEvent         bool            `json:"isPubEvent"`         // 是否发布消息
	AggField           string          `json:"aggField"`           // 聚合根字段
	Env                IEnvConfig      `json:"env"`                // 环境配置
	DBSchema           *store.DBSchema `json:"schema"`             // 数据结构
	Database           any             `json:"database"`           // 数据库连接对象
	IsCancelModified   bool            `json:"isCancelModified"`   // 取消创建者与更新都信息
	IsCancelSoftDelete bool            `json:"isCancelSoftDelete"` // 取消软删除
	DaoType            string          `json:"daoType"`            // 节点类型 在neo4j: node, rel
}

func (c *DaoConfig) Valid() {
	if c.DbKey == "" {
		panic("DaoConfig empty DbKey")
	}
	if c.Env == nil {
		panic("DaoConfig empty Env")
	}
	if c.DBSchema == nil {
		panic("DaoConfig empty DBSchema")
	}
}

func (c *DaoConfig) GetDbKey() string {
	return c.DbKey
}

func (c *DaoConfig) GetIsPubEvent() bool {
	return c.IsPubEvent
}
func (c *DaoConfig) GetAggField() string {
	return c.AggField
}

func (c *DaoConfig) GetEnv() restapp.IEnvConfig {
	return c.Env
}

func (c *DaoConfig) GetSchema() *store.DBSchema {
	return c.DBSchema
}
