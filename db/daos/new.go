package daos

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos/impl/mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos/impl/neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos/impl/sql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
)

type NewConfig struct {
	DbKey      string           `json:"dbKey"`
	IsPubEvent bool             `json:"isPubEvent"`
	AggField   string           `json:"aggField"`
	DbSchema   *dbschema.Schema `json:"dbSchema"`
}

var daoMap = types.NewCMap[any]()

func NewDao[T any](newCfg *NewConfig) idao.Dao[T] {
	dbKey := newCfg.DbKey
	if dbKey == "" {
		dbKey = "default"
	}
	if newCfg.DbSchema == nil {
		panic("db.NewDao() args schema is required")
	}

	tableName := newCfg.DbSchema.Name
	daoKey := getDaoKey(dbKey, tableName)

	if v, ok := daoMap.Get(daoKey); v != nil && ok {
		return v.(idao.Dao[T])
	}

	daoCfg := &idao.DaoConfig{
		DbKey:      newCfg.DbKey,
		IsPubEvent: newCfg.IsPubEvent,
		AggField:   newCfg.AggField,
		Env:        restapp.GetEnvConfig(),
		Schema:     newCfg.DbSchema,
	}

	var dao idao.Dao[T]
	item := getDbItem(dbKey)

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

	daoMap.Add(daoKey, any(dao))
	return dao
}

func getDaoKey(dbKey, tableName string) string {
	return fmt.Sprintf("%s.%s", dbKey, tableName)
}

func getDbItem(dbKey string) restapp.DBItem {
	item := restapp.GetDb(dbKey)
	if item == nil {
		panic(errors.New(" %s dbKey not exists", dbKey))
	}
	return item
}
