package dao

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/impl/mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/impl/neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/impl/sql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
)

type NewConfig struct {
	DBKey              string          `json:"dbKey"`              // 可选 DBKey 与 DB
	DB                 any             `json:"db"`                 // 可选 DBKey 与 DB
	IsPubEvent         *bool           `json:"isPubEvent"`         // 可选
	AggField           string          `json:"aggField"`           // 可选
	AggType            string          `json:"aggType"`            // 可选
	DBSchema           *store.DBSchema `json:"dbSchema"`           // 可选
	TableName          string          `json:"tableName"`          // 可选
	GraphType          idao.GraphType  `json:"graphType"`          // 可选
	GraphLabels        []string        `json:"graphLabels"`        // 可选
	IsCancelModified   bool            `json:"isCancelModified"`   // 可选
	IsCancelSoftDelete bool            `json:"isCancelSoftDelete"` // 可选
	Env                *env.Env        `json:"-"`                  // 可选
}

var cache = types.NewCMap[any]()

var isFalse bool = false
var isTrue bool = true

func IsFalse() *bool {
	return &isFalse
}

func IsTrue() *bool {
	return &isTrue
}

func NewDao[T any](newCfg *NewConfig) idao.Dao[T] {
	if newCfg == nil {
		panic("new dao must have a non-nil pointer")
	}

	dbKey := newCfg.DBKey
	tableName := newCfg.TableName
	envInst := newCfg.Env
	if envInst == nil {
		envInst = env.GetEnv()
	}
	if newCfg.DBKey == "" {
		db := env.GetDBDefault()
		newCfg.DBKey = db.GetDBKey()
	}
	if tableName == "" && newCfg.DBSchema != nil {
		tableName = newCfg.DBSchema.TableName
	}

	// dao缓存 取得daoKey
	className := reflectutils.GetClassName[T]()
	daoKey := getDaoKey(dbKey, tableName, className)
	if v, ok := cache.Get(daoKey); v != nil && ok {
		return v.(idao.Dao[T])
	}

	// 是struct类型
	if newCfg.DBSchema == nil && reflectutils.IsStruct[T]() {
		entity := reflectutils.NewInstance[T]()
		if tableName == "" {
			tableName = reflectutils.GetStructName(entity)
			tableName = stringutils.SnakeString(tableName)
		}
		dbSchema := dbschema.NewDBSchemaWithStruct(tableName, entity, tableName)
		newCfg.DBSchema = dbSchema
	}

	if newCfg.DBSchema == nil {
		panic("db.NewDao() args schema is required")
	}

	item := getDBItem(envInst, dbKey)
	if item == nil && newCfg.DB != nil {
		newCfg.DB = newCfg.DB
	}

	isPubEvent := false
	if newCfg.IsPubEvent != nil {
		isPubEvent = *newCfg.IsPubEvent
	} else {
		isPubEvent = env.GetEnv().App.IsPubEvent
	}

	daoCfg := &idao.DaoConfig{
		DBKey:              newCfg.DBKey,
		DB:                 newCfg.DB,
		IsPubEvent:         isPubEvent,
		AggField:           newCfg.AggField,
		AggType:            newCfg.AggType,
		Env:                envInst,
		DBSchema:           newCfg.DBSchema,
		GraphType:          newCfg.GraphType,
		GraphLabels:        newCfg.GraphLabels,
		IsCancelModified:   newCfg.IsCancelModified,
		IsCancelSoftDelete: newCfg.IsCancelSoftDelete,
	}

	if isPubEvent {
		daoCfg.OutboxDao = NewOutboxDao(newCfg.DBKey)
	}
	var dao idao.Dao[T]

	if item != nil {
		switch item.GetDBType() {
		case env.DBType_MongoDB:
			dao = mongodb.NewDao[T](daoCfg)
		case env.DBType_Sqlite,
			env.DBType_Oracle,
			env.DBType_Postgres,
			env.DBType_MySQL,
			env.DBType_MsSQL:
			dao = sql.NewDao[T](daoCfg)
		case env.DBType_Neo4j:
			dao = neo4j.NewDao[T](daoCfg)
		default:
			panic(errors.New(fmt.Sprintf("%s database not exists", dbKey)))
		}
	}
	cache.Add(daoKey, any(dao))
	return dao
}

func getDaoKey(dbKey, tableName string, className string) string {
	return fmt.Sprintf("%s.%s.%s", dbKey, tableName, className)
}

func getDBItem(env *env.Env, dbKey string) env.DBItem {
	item := env.GetDB(dbKey)
	return item
}
