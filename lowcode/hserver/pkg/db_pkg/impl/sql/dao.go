package sql

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_sql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/impl"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"gorm.io/gorm"
	dbschema "gorm.io/gorm/schema"
)

type Dao struct {
	*impl.DaoBase
	cfg      *db.DaoConfig
	db       *gorm.DB
	dao      *ddd_sql.MapDao
	dbSchema *dbschema.Schema
}

func NewDao(cfg *db.DaoConfig) db.Dao {
	cfg.Valid()
	var database *gorm.DB
	//eb := ddd.NewMapEntityBuilder[map[string]any]()
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
	dbSchema, err := NewDBSchema(cfg.Schema)
	if err != nil {
		panic(err)
	}

	sqlDao := ddd_sql.NewMapDao(database, cfg.DbKey, cfg.Schema.GetTableName())
	sqlDao.GetMetadata()["dbSchema"] = dbSchema
	sqlDao.GetMetadata()["schema"] = cfg.Schema

	daoBase := impl.NewDaoBase(sqlDao, cfg)

	return &Dao{
		DaoBase:  daoBase,
		dao:      sqlDao,
		cfg:      cfg,
		db:       database,
		dbSchema: dbSchema,
	}
}

func (d *Dao) Table() db.Table {
	return newTable(d.db, d.cfg.Schema, d.dbSchema)
}
