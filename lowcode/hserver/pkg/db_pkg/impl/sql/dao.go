package sql

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_sql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"gorm.io/gorm"
)

type MapEntity = ddd.MapEntity
type Dao struct {
	*db.DaoBase
	cfg *db.DaoConfig
	db  *gorm.DB
	dao *ddd_sql.Dao[MapEntity]
}

func NewDao(cfg *db.DaoConfig) db.Dao {
	cfg.Valid()
	
	entBuilder := ddd.NewMapEntityBuilder[ddd.MapEntity]()
	item := restapp.GetDb(cfg.DbKey)
	if item == nil {
		panic("db item not found")
	}
	database, ok := item.GetDB().(*gorm.DB)
	if !ok {
		panic("db item not found")
	}
	dao := ddd_sql.NewDao[MapEntity](database, cfg.DbKey, entBuilder, cfg.TableName)
	daoBase := db.NewDaoBase(dao, cfg)
	return &Dao{
		DaoBase: daoBase,
		dao:     dao,
		cfg:     cfg,
		db:      database,
	}
}

func (d *Dao) Table() db.Table {
	return NewTable(d.db, d.cfg.GetTableName(), d.cfg.Schema)
}
