package idao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbevent"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"strings"
)

type RelType string
type FieldStyle string

type DaoConfig struct {
	DbKey              string          `json:"dbKey"`              // 数据库Key
	IsPubEvent         bool            `json:"isPubEvent"`         // 是否发布消息
	AggField           string          `json:"aggField"`           // 聚合根字段
	AggType            string          `json:"aggType"`            // 聚合根类型
	Env                *env.Env        `json:"env"`                // 环境配置
	DBSchema           *store.DBSchema `json:"schema"`             // 数据结构
	DB                 any             `json:"db"`                 // 数据库连接对象
	IsCancelModified   bool            `json:"isCancelModified"`   // 取消创建者与更新都信息
	IsCancelSoftDelete bool            `json:"isCancelSoftDelete"` // 取消软删除
	RefType            RelType         `json:"RefType"`            // 节点类型 在neo4j: node, rel
	PropertyIsField    bool            `json:"propertyIsField"`    // 属性名称即字段名称
	OutboxDao          Dao[*dbevent.Outbox]
}

const (
	RelType_None RelType = ""
	RelType_Node RelType = "node"
	RelType_Rel  RelType = "rel"
)

func GetRefType(val string) (RelType, error) {
	val = strings.ToLower(val)
	switch val {
	case "":
		return RelType_None, nil
	case "node":
		return RelType_Node, nil
	case "rel":
		return RelType_Rel, nil
	default:
		return "", errors.New("invalid DaoType " + val)
	}
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

func (c *DaoConfig) GetEnv() *env.Env {
	return c.Env
}

func (c *DaoConfig) GetSchema() *store.DBSchema {
	return c.DBSchema
}
