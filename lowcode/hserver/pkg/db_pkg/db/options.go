package db

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/common"
	"github.com/liuxd6825/jsonschema/v6"
)

type DaoConfig struct {
	DbKey              string             `json:"dbKey"`              // 数据库Key
	IsPubEvent         bool               `json:"isPubEvent"`         // 是否发布消息
	AggField           string             `json:"aggField"`           // 聚合根字段
	Env                common.IEnvConfig  `json:"env"`                // 环境配置
	Schema             *jsonschema.Schema `json:"schema"`             // 数据结构
	Database           any                `json:"database"`           // 数据库连接对象
	IsCancelModified   bool               `json:"isCancelModified"`   // 取消创建者与更新都信息
	IsCancelSoftDelete bool               `json:"isCancelSoftDelete"` // 取消软删除
}

func (c *DaoConfig) Valid() {
	if c.DbKey == "" {
		panic("DaoConfig empty dbKey")
	}
	if c.Env == nil {
		panic("DaoConfig empty env")
	}
	if c.Schema == nil {
		panic("DaoConfig empty schema")
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

func (c *DaoConfig) GetEnv() common.IEnvConfig {
	return c.Env
}

func (c *DaoConfig) GetSchema() *jsonschema.Schema {
	return c.Schema
}
