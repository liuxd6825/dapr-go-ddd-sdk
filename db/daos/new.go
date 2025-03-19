package daos

import (
	"fmt"
	restapp2 "github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos/impl/mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos/impl/neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos/impl/sql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
)

type NewConfig struct {
	DbKey        string           `json:"dbKey"`
	Db           any              `json:"db"`
	EventPublish *bool            `json:"eventPublish"`
	AggField     string           `json:"aggField"`
	DbSchema     *dbschema.Schema `json:"dbSchema"`
	TableName    string           `json:"tableName"`
}

var daoMap = types.NewCMap[any]()

func NewDao[T any](newCfg *NewConfig) idao.Dao[T] {
	if newCfg == nil {
		panic("new dao must have a non-nil pointer")
	}
	if newCfg.DbKey == "" {
		db := restapp2.GetDBDefault()
		newCfg.DbKey = db.GetDBKey()
	}
	dbKey := newCfg.DbKey
	tableName := newCfg.TableName

	if newCfg.DbSchema == nil && !reflectutils.IsMap[T]() {
		entity := reflectutils.NewInstance[T]()
		if tableName == "" {
			tableName = reflectutils.GetStructName(entity)
			tableName = stringutils.SnakeString(tableName)
		}
		dbSchema := dbschema.NewSchemaWithStruct(tableName, entity, tableName)
		newCfg.DbSchema = dbSchema
	}

	if newCfg.DbSchema == nil {
		panic("db.NewDao() args schema is required")
	}

	if tableName == "" && newCfg.DbSchema != nil {
		tableName = newCfg.DbSchema.TableName
	}
	daoKey := getDaoKey(dbKey, tableName)

	if v, ok := daoMap.Get(daoKey); v != nil && ok {
		return v.(idao.Dao[T])
	}

	item := getDBItem(dbKey)

	eventPublish := false
	if newCfg.EventPublish != nil {
		eventPublish = *newCfg.EventPublish
	} else if ep, ok := item.GetConfig().(restapp2.EventPublish); ok {
		eventPublish = ep.GetEventPublish()
	}

	daoCfg := &idao.DaoConfig{
		DbKey:      newCfg.DbKey,
		Database:   newCfg.Db,
		IsPubEvent: eventPublish,
		AggField:   newCfg.AggField,
		Env:        restapp2.GetEnvConfig(),
		Schema:     newCfg.DbSchema,
	}

	var dao idao.Dao[T]

	switch item.GetDBType() {
	case restapp2.DbType_MongoDB:
		dao = mongodb.NewDao[T](daoCfg)
	case restapp2.DbType_Sqlite,
		restapp2.DbType_Oracle,
		restapp2.DbType_Postgres,
		restapp2.DbType_MySQL,
		restapp2.DbType_MsSQL:
		dao = sql.NewDao[T](daoCfg)
	case restapp2.DbType_Redis:
		panic(errors.New(fmt.Sprintf("%s database nonsupport Redis", dbKey)))
	case restapp2.DbType_Neo4j:
		dao = neo4j.NewDao[T](daoCfg)
	default:
		panic(errors.New(fmt.Sprintf("%s database not exists", dbKey)))
	}

	daoMap.Add(daoKey, any(dao))
	return dao
}

func getDaoKey(dbKey, tableName string) string {
	return fmt.Sprintf("%s.%s", dbKey, tableName)
}

func getDBItem(dbKey string) restapp2.DBItem {
	item := restapp2.GetDB(dbKey)
	if item == nil {
		panic(errors.New(" %s dbKey not exists", dbKey))
	}
	return item
}
