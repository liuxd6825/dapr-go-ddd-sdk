package sql

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_sql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"gorm.io/gorm"
	dbschema "gorm.io/gorm/schema"
)

type MapEntity = ddd.MapEntity
type Dao struct {
	*db.DaoBase
	cfg      *db.DaoConfig
	db       *gorm.DB
	dao      *ddd_sql.Dao[MapEntity]
	dbSchema *dbschema.Schema
}

func NewDao(cfg *db.DaoConfig) db.Dao {
	cfg.Valid()
	var database *gorm.DB
	eb := ddd.NewMapEntityBuilder[ddd.MapEntity]()
	if cfg.Database != nil {
		if val, ok := cfg.Database.(*gorm.DB); ok {
			database = val
		} else {
			panic("database config error")
		}
	}

	if database == nil {
		item := restapp.GetDb(cfg.DbKey)
		if item == nil {
			panic("db item not found")
		}
		if val, ok := item.GetDB().(*gorm.DB); ok {
			database = val
		} else {
			panic("database config error")
		}
	}

	dao := ddd_sql.NewDao[MapEntity](database, cfg.DbKey, eb, cfg.Schema.GetTableName())
	daoBase := db.NewDaoBase(dao, cfg)
	dbSchema, err := NewDBSchema(cfg.Schema)
	if err != nil {
		panic(err)
	}
	return &Dao{
		DaoBase:  daoBase,
		dao:      dao,
		cfg:      cfg,
		db:       database,
		dbSchema: dbSchema,
	}
}

func (d *Dao) Table() db.Table {
	return newTable(d.db, d.cfg.Schema, d.dbSchema)
}

func (d *Dao) Create(ctx context.Context, entity ddd.MapEntity, opts ...*db.CallOptions) {
	d.DaoBase.Create(ctx, entity, opts...)
}
