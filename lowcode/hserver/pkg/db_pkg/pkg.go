package db_pkg

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/impl/mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/impl/sql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/jsonschema/v6"
	"gorm.io/gorm"
)

type Pkg struct {
	server element.Server
	daoMap *types.CMap[db.Dao]
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
		daoMap: types.NewCMap[db.Dao](),
	}
}

func (p *Pkg) GetDao(daoKey string) db.Dao {
	if v, ok := p.daoMap.Get(daoKey); v != nil && ok {
		return v.(db.Dao)
	}
	panic(fmt.Errorf("dao not found %s", daoKey))
}

func (p *Pkg) NewDao(opts *NewDaoConfig) db.Dao {
	dbKey := opts.DbKey
	if dbKey == "" {
		dbKey = "default"
	}
	if opts.Schema == nil {
		panic("db.NewDao() args schema is required")
	}

	tableName := opts.Schema.Name
	daoKey := fmt.Sprintf("%s.%s", dbKey, tableName)
	if v, ok := p.daoMap.Get(daoKey); v != nil && ok {
		return v.(db.Dao)
	}

	cfg := &db.DaoConfig{
		DbKey:      opts.DbKey,
		IsPubEvent: opts.IsPubEvent,
		AggField:   opts.AggField,
		Env:        p.server.GetEnvCfg(),
		Schema:     opts.Schema,
	}

	var dao db.Dao
	item := p.getDbItem(dbKey)

	switch item.GetDBType() {
	case restapp.DbType_MongoDB:
		dao = mongodb.NewDao(cfg)
	case restapp.DbType_Sqlite,
		restapp.DbType_Oracle,
		restapp.DbType_Postgres,
		restapp.DbType_MySQL,
		restapp.DbType_MsSQL:
		dao = sql.NewDao(cfg)
	case restapp.DbType_Redis:
		panic(errors.New(fmt.Sprintf("%s database nonsupport Redis", dbKey)))
	case restapp.DbType_Neo4j:
		panic(errors.New(fmt.Sprintf("%s database nonsupport Neo4j", dbKey)))
	default:
		panic(errors.New(fmt.Sprintf("%s database not exists", dbKey)))
	}

	p.daoMap.Add(daoKey, dao)
	return dao
}

type NewTableOptions struct {
	DbKey     string             `json:"dbKey"`
	TableName string             `json:"tableName"`
	Schema    *jsonschema.Schema `json:"schema"`
}

func (p *Pkg) NewTable(opts *NewTableOptions) db.Table {
	dbKey := opts.DbKey
	if dbKey == "" {
		dbKey = "default"
	}

	var table db.Table
	item := p.getDbItem(opts.DbKey)

	switch item.GetDBType() {
	case restapp.DbType_MongoDB:
		table = mongodb.NewTable(item.GetMongo(), opts.Schema)
	case restapp.DbType_Sqlite,
		restapp.DbType_Oracle,
		restapp.DbType_Postgres,
		restapp.DbType_MySQL,
		restapp.DbType_MsSQL:
		database, ok := item.GetDB().(*gorm.DB)
		if ok {
			panic(errors.New(fmt.Sprintf("%s database nonsupport gorm.DB", dbKey)))
		}
		table = sql.NewTable(database, opts.Schema)
	case restapp.DbType_Redis:
		panic(errors.New(fmt.Sprintf("%s database nonsupport Redis", dbKey)))
	case restapp.DbType_Neo4j:
		panic(errors.New(fmt.Sprintf("%s database nonsupport Neo4j", dbKey)))
	default:
		panic(errors.New(fmt.Sprintf("%s database not exists", dbKey)))
	}
	return table
}

func (p *Pkg) getDbItem(dbKey string) restapp.DBItem {
	item := restapp.GetDb(dbKey)
	if item == nil {
		panic(errors.New(" %s dbKey not exists", dbKey))
	}
	return item
}

func (p *Pkg) StartTx(ctx context.Context, dbKeys []string, txFunc ddd_repository.TxFunc, options ...*ddd_repository.SessionOptions) (err error) {
	return tx.StartTx(ctx, dbKeys, txFunc, options...)
}
