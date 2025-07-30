package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"strings"
	"time"
)

type CDCRecord struct {
	DB           string             `json:"db"`     // 数据库
	Table        string             `json:"table"`  // 数据表
	Before       map[string]any     `json:"before"` // 之前数据
	After        map[string]any     `json:"after"`  // 之后数据
	OpType       OpType             `json:"opType"` // 操作状态 "r" for read/backfill, "c" for create, "u" for update, "d" for delete
	CdcTimestamp time.Time          `json:"cdcTimestamp"`
	DBSchema     *dbschema.DBSchema `json:"-"`
}

type OpType string

const (
	OpTypeRead   OpType = "r"
	OpTypeCreate OpType = "c"
	OpTypeUpdate OpType = "u"
	OpTypeDelete OpType = "d"
)

// IsMaster 是主数据
func (r *CDCRecord) IsMaster() bool {
	return !strings.Contains(r.Table, "_")
}

// IsRelation 是关系数据
func (r *CDCRecord) IsRelation() bool {
	return strings.Contains(r.Table, "_")
}

// IsRename 是否数据更新
func (r *CDCRecord) IsRename() bool {
	if r.OpType == "u" {
		newName, _ := maputils.GetString(r.After, "name", "")
		oldName, _ := maputils.GetString(r.Before, "name", "")
		if newName != oldName {
			return true
		}
	}
	return false
}

// IsChangedRelType 是否数据更新
func (r *CDCRecord) IsChangedRelType() bool {
	if r.OpType == "u" {
		newType, _ := maputils.GetString(r.After, "relation_type", "")
		oldType, _ := maputils.GetString(r.Before, "relation_type", "")
		if newType != oldType {
			return true
		}
	}
	return false
}

// IsChangedBusFields 是否更新业务字段
func (r *CDCRecord) IsChangedBusFields() bool {
	if r.OpType == "u" {
		for k, _ := range r.After {
			if k != "relation_type" && k != "name" {
				return true
			}
		}
	}
	return false
}

func (r *CDCRecord) AfterMap() map[string]any {
	return r.After
}

func (r *CDCRecord) BeforeMap() map[string]any {
	return r.Before
}

func (r *CDCRecord) AfterName() string {
	val, _ := maputils.GetString(r.After, "name", "")
	return val
}

func (r *CDCRecord) BeforeName() string {
	val, _ := maputils.GetString(r.Before, "name", "")
	return val
}

func (r *CDCRecord) newMap(vals map[string]any) map[string]any {
	mapData := make(map[string]any)
	for key, value := range vals {
		propName := stringutils.FirstLowerCamelString(key)
		mapData[propName] = value
	}
	return mapData
}

func (r *CDCRecord) SetAfterValue(key string, val any) *CDCRecord {
	r.After[key] = val
	return r
}

func (r *CDCRecord) SetBeforeValue(key string, val any) *CDCRecord {
	r.Before[key] = val
	return r
}
