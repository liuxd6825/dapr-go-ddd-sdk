package db

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
)

type DaoConfig struct {
	DbKey      string            `json:"dbKey"`
	IsPubEvent *bool             `json:"isPubEvent"`
	AggField   string            `json:"aggField"`
	Env        common.IEnvConfig `json:"env"`
	Schema     *schema.Schema    `json:"schema"`
}

func (c *DaoConfig) Valid() {
	if c.DbKey == "" {
		panic("empty dbKey")
	}
	if c.Env == nil {
		panic("empty env")
	}
	if c.Schema == nil {
		panic("empty schema")
	}
}

func (c *DaoConfig) GetDbKey() string {
	return c.DbKey
}

func (c *DaoConfig) GetIsPubEvent() bool {
	if c.IsPubEvent == nil {
		return true
	}
	return *c.IsPubEvent
}
func (c *DaoConfig) GetAggField() string {
	return c.AggField
}

func (c *DaoConfig) GetEnv() common.IEnvConfig {
	return c.Env
}

func (c *DaoConfig) GetSchema() *schema.Schema {
	return c.Schema
}
