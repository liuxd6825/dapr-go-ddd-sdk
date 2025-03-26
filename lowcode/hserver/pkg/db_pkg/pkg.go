package db_pkg

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	idao2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	mongodb2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/impl/mongodb"
	neo4j2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/impl/neo4j"
	sql2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/impl/sql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/jsonschema/v6"
	neo4jdriver "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"gorm.io/gorm"
)

type Pkg struct {
	server element.Server
	daoMap *types.CMap[any]
}

type NewDaoConfig struct {
	DbKey      string             `json:"dbKey"`
	IsPubEvent bool               `json:"isPubEvent"`
	AggField   string             `json:"aggField"`
	Schema     *jsonschema.Schema `json:"schema"`
}

func New(server element.Server) *Pkg {
	return &Pkg{
		server: server,
		daoMap: types.NewCMap[any](),
	}
}

func (p *Pkg) NewRSQLBuilder() *rsql.Builder {
	return rsql.NewBuilder()
}

func (p *Pkg) NewDao(opts *NewDaoConfig) idao2.Dao[map[string]any] {
	dbKey := opts.DbKey
	if dbKey == "" {
		dbKey = "default"
	}
	if opts.Schema == nil {
		panic("db.NewDao() args schema is required")
	}

	tableName := opts.Schema.Name()
	daoKey := getDaoKey(dbKey, tableName)
	if v, ok := p.daoMap.Get(daoKey); v != nil && ok {
		return v.(idao2.Dao[map[string]any])
	}

	cfg := &idao2.DaoConfig{
		DbKey:      opts.DbKey,
		IsPubEvent: opts.IsPubEvent,
		AggField:   opts.AggField,
		Env:        p.server.GetEnvCfg(),
		DBSchema:   dbschema.NewDBSchemaWithJsonSchema(opts.Schema),
	}

	var dao idao2.Dao[map[string]any]
	item := p.getDbItem(dbKey)

	switch item.GetDBType() {
	case restapp.DbType_MongoDB:
		dao = mongodb2.NewDao[map[string]any](cfg)
	case restapp.DbType_Sqlite,
		restapp.DbType_Oracle,
		restapp.DbType_Postgres,
		restapp.DbType_MySQL,
		restapp.DbType_MsSQL:
		dao = sql2.NewDao[map[string]any](cfg)
	case restapp.DbType_Redis:
		panic(errors.New(fmt.Sprintf("%s database nonsupport Redis", dbKey)))
	case restapp.DbType_Neo4j:
		dao = neo4j2.NewDao[map[string]any](cfg)
	default:
		panic(errors.New(fmt.Sprintf("%s database not exists", dbKey)))
	}

	p.daoMap.Add(daoKey, dao)
	return dao
}

func (p *Pkg) GetDao(dbKey, tableName string) idao2.Dao[map[string]any] {
	daoKey := getDaoKey(dbKey, tableName)
	if v, ok := p.daoMap.Get(daoKey); v != nil && ok {
		return v.(idao2.Dao[map[string]any])
	}
	panic(fmt.Errorf("dao not found %s", daoKey))
}

type NewTableOptions struct {
	DbKey     string             `json:"dbKey"`
	TableName string             `json:"tableName"`
	Schema    *jsonschema.Schema `json:"schema"`
}

func (p *Pkg) NewTable(opts *NewTableOptions) idao2.Table {
	dbKey := opts.DbKey
	if dbKey == "" {
		dbKey = "default"
	}

	dbSch := dbschema.NewDBSchemaWithJsonSchema(opts.Schema)

	var table idao2.Table
	item := p.getDbItem(opts.DbKey)

	switch item.GetDBType() {
	case restapp.DbType_MongoDB:
		table = mongodb2.NewTable(item.GetMongo(), dbSch)
	case restapp.DbType_Sqlite,
		restapp.DbType_Oracle,
		restapp.DbType_Postgres,
		restapp.DbType_MySQL,
		restapp.DbType_MsSQL:
		database, ok := item.GetDB().(*gorm.DB)
		if ok {
			panic(errors.New(fmt.Sprintf("%s database nonsupport gorm.DB", dbKey)))
		}
		table = sql2.NewTable(database, dbSch, map[string]any{})
	case restapp.DbType_Neo4j:
		if driver, ok := item.GetDB().(neo4jdriver.DriverWithContext); ok {
			table = neo4j2.NewTable(driver, dbSch)
		} else {
			panic(errors.New(fmt.Sprintf("%s database nonsupport neo4j.Driver", dbKey)))
		}
	case restapp.DbType_Redis:
		panic(errors.New(fmt.Sprintf("%s database nonsupport Redis", dbKey)))

	default:
		panic(errors.New(fmt.Sprintf("%s database not exists", dbKey)))
	}
	return table
}

func (p *Pkg) getDbItem(dbKey string) restapp.DBItem {
	item := restapp.GetDB(dbKey)
	if item == nil {
		panic(errors.New(" %s dbKey not exists", dbKey))
	}
	return item
}

func (p *Pkg) StartTx(ctx context.Context, dbKeys []string, txFunc store.TxFunc, options ...*store.SessionOptions) (err error) {
	return tx.StartTx(ctx, dbKeys, txFunc, options...)
}

func getDaoKey(dbKey, tableName string) string {
	return fmt.Sprintf("%s.%s", dbKey, tableName)
}
