package neo4j

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/impl"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Dao struct {
	*impl.DaoBase
	cfg    *db.DaoConfig
	dao    ddd_repository.Dao[map[string]any]
	driver neo4j.DriverWithContext
}

func NewDao(cfg *db.DaoConfig) db.Dao {
	cfg.Valid()
	var driver neo4j.DriverWithContext
	//eb := ddd.NewMapEntityBuilder[map[string]any]()
	if cfg.Database != nil {
		if val, ok := cfg.Database.(neo4j.DriverWithContext); ok {
			driver = val
		} else {
			panic("database config error neo4j.DriverWithContext")
		}
	}

	if driver == nil {
		item := restapp.GetDb(cfg.DbKey)
		if item == nil {
			panic(fmt.Sprintf("dbKey %s not found", cfg.DbKey))
		}
		if val, ok := item.GetDB().(neo4j.DriverWithContext); ok {
			driver = val
		} else {
			panic(fmt.Sprintf("dbKey %s is not neo4j.DriverWithContext", cfg.DbKey))
		}
	}

	dao := ddd_neo4j.NewMapNodeDao(driver, []string{cfg.Schema.Name}, cfg.Schema)
	daoBase := impl.NewDaoBase(dao, cfg)

	return &Dao{
		DaoBase: daoBase,
		dao:     dao,
		cfg:     cfg,
		driver:  driver,
	}
}

func (d *Dao) Table() db.Table {
	return newTable(d.driver, d.cfg.Schema)
}

func (d *Dao) GetTableName() string {
	return d.cfg.Schema.Name
}
