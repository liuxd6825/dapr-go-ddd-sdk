package db

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/impl/mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/impl/neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/impl/sql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/jsonschema/v6"
)

type NewDaoConfig struct {
	DbKey      string             `json:"dbKey"`
	IsPubEvent bool               `json:"isPubEvent"`
	AggField   string             `json:"aggField"`
	Schema     *jsonschema.Schema `json:"schema"`
}

var daoMap = types.NewCMap[db.Dao]()

func NewDao[T any](opts *NewDaoConfig) db.Dao[T] {
	dbKey := opts.DbKey
	if dbKey == "" {
		dbKey = "default"
	}
	if opts.Schema == nil {
		panic("db.NewDao() args schema is required")
	}

	tableName := opts.Schema.Name
	daoKey := getDaoKey(dbKey, tableName)
	if v, ok := daoMap.Get(daoKey); v != nil && ok {
		return v.(db.Dao[T])
	}

	cfg := &db.DaoConfig{
		DbKey:      opts.DbKey,
		IsPubEvent: opts.IsPubEvent,
		AggField:   opts.AggField,
		Env:        restapp.GetEnvConfig(),
		Schema:     opts.Schema,
	}

	var dao db.Dao[T]
	item := getDbItem(dbKey)

	switch item.GetDBType() {
	case restapp.DbType_MongoDB:
		dao = mongodb.NewDao[T](cfg)
	case restapp.DbType_Sqlite,
		restapp.DbType_Oracle,
		restapp.DbType_Postgres,
		restapp.DbType_MySQL,
		restapp.DbType_MsSQL:
		dao = sql.NewDao[T](cfg)
	case restapp.DbType_Redis:
		panic(errors.New(fmt.Sprintf("%s database nonsupport Redis", dbKey)))
	case restapp.DbType_Neo4j:
		dao = neo4j.NewDao[T](cfg)
	default:
		panic(errors.New(fmt.Sprintf("%s database not exists", dbKey)))
	}

	daoMap.Add(daoKey, dao)
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
