package dao

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/impl/mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/impl/neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/impl/sql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
)

type NewConfig struct {
	DBKey        string          `json:"dbKey"`
	DB           any             `json:"db"`
	EventPublish *bool           `json:"eventPublish"`
	AggField     string          `json:"aggField"`
	DBSchema     *store.DBSchema `json:"dbSchema"`
	TableName    string          `json:"tableName"`
}

var cache = types.NewCMap[any]()

func NewDao[T any](newCfg *NewConfig) idao.Dao[T] {
	if newCfg == nil {
		panic("new dao must have a non-nil pointer")
	}
	if newCfg.DBKey == "" {
		db := restapp.GetDBDefault()
		newCfg.DBKey = db.GetDBKey()
	}
	dbKey := newCfg.DBKey
	tableName := newCfg.TableName

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

	if tableName == "" && newCfg.DBSchema != nil {
		tableName = newCfg.DBSchema.TableName
	}

	// dao缓存 取得daoKey
	className := reflectutils.GetClassName[T]()
	daoKey := getDaoKey(dbKey, tableName, className)
	if v, ok := cache.Get(daoKey); v != nil && ok {
		return v.(idao.Dao[T])
	}

	item := getDBItem(dbKey)
	eventPublish := false
	if newCfg.EventPublish != nil {
		eventPublish = *newCfg.EventPublish
	} else if ep, ok := item.GetConfig().(restapp.EventPublish); ok {
		eventPublish = ep.GetEventPublish()
	}

	daoCfg := &idao.DaoConfig{
		DbKey:      newCfg.DBKey,
		Database:   newCfg.DB,
		IsPubEvent: eventPublish,
		AggField:   newCfg.AggField,
		Env:        restapp.GetEnvConfig(),
		DBSchema:   newCfg.DBSchema,
	}

	var dao idao.Dao[T]

	switch item.GetDBType() {
	case restapp.DbType_MongoDB:
		dao = mongodb.NewDao[T](daoCfg)
	case restapp.DbType_Sqlite,
		restapp.DbType_Oracle,
		restapp.DbType_Postgres,
		restapp.DbType_MySQL,
		restapp.DbType_MsSQL:
		dao = sql.NewDao[T](daoCfg)
	case restapp.DbType_Redis:
		panic(errors.New(fmt.Sprintf("%s database nonsupport Redis", dbKey)))
	case restapp.DbType_Neo4j:
		dao = neo4j.NewDao[T](daoCfg)
	default:
		panic(errors.New(fmt.Sprintf("%s database not exists", dbKey)))
	}

	cache.Add(daoKey, any(dao))
	return dao
}

func getDaoKey(dbKey, tableName string, className string) string {
	return fmt.Sprintf("%s.%s.%s", dbKey, tableName, className)
}

func getDBItem(dbKey string) restapp.DBItem {
	item := restapp.GetDB(dbKey)
	if item == nil {
		panic(errors.New(" %s dbKey not exists", dbKey))
	}
	return item
}
