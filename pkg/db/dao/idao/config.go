package idao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbevent"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"strings"
)

type GraphType string
type FieldStyle string

type DaoConfig struct {
	DBKey string `json:"dbKey"` // 数据库Key
	//IsPubEvent         bool            `json:"isPubEvent"`         // 是否发布消息
	AggField           string          `json:"aggField"`           // 聚合根字段
	AggType            string          `json:"aggType"`            // 聚合根类型
	Env                *env.Env        `json:"env"`                // 环境配置
	DBSchema           *store.DBSchema `json:"schema"`             // 数据结构
	DB                 any             `json:"db"`                 // 数据库连接对象
	IsCancelModified   bool            `json:"isCancelModified"`   // 取消创建者与更新都信息
	IsCancelSoftDelete bool            `json:"isCancelSoftDelete"` // 取消软删除
	GraphType          GraphType       `json:"graphType"`          // 节点类型 在neo4j: node, rel
	GraphLabels        []string        `json:"graphLabels"`
	OutboxDao          Dao[*dbevent.OutboxEvent]
}

const (
	GraphType_None GraphType = ""
	GraphType_Node GraphType = "node"
	GraphType_Rel  GraphType = "rel"
)

func GetRelType(val string) (GraphType, error) {
	val = strings.ToLower(val)
	switch val {
	case "":
		return GraphType_None, nil
	case "node":
		return GraphType_Node, nil
	case "rel":
		return GraphType_Rel, nil
	default:
		return "", errors.New("invalid DaoType " + val)
	}
}

func IsGraphType(val []string, graphType GraphType) bool {
	for _, v := range val {
		if v == string(graphType) {
			return true
		}
	}
	return false
}

func (c *DaoConfig) Valid() {
	if c.DBKey == "" {
		panic("DaoConfig empty DBKey")
	}
	if c.Env == nil {
		panic("DaoConfig empty Env")
	}
	if c.DBSchema == nil {
		panic("DaoConfig empty DBSchema")
	}
}

func (c *DaoConfig) GetDBKey() string {
	return c.DBKey
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
