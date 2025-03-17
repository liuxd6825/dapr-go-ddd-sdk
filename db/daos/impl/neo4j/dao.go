package neo4j

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos/impl"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Dao[T any] struct {
	*impl.DaoBase[T]
	cfg    *idao.DaoConfig
	dao    ddd_repository.Dao[T]
	driver neo4j.DriverWithContext
}

func NewDao[T any](cfg *idao.DaoConfig, tableName ...string) idao.Dao[T] {
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

	labels := []string{cfg.Schema.Name}
	if len(tableName) > 0 {
		labels = tableName
	}
	dbSch := cfg.Schema
	var dao ddd_repository.Dao[T]
	if cfg.DaoType == "node" {
		dao = ddd_neo4j.NewNodeDao[T](driver, dbSch, labels)
	} else {
		dao = ddd_neo4j.NewRelationDao[T](driver, dbSch, labels)
	}
	daoBase := impl.NewDaoBase[T](dao, cfg)
	return &Dao[T]{
		DaoBase: daoBase,
		dao:     dao,
		cfg:     cfg,
		driver:  driver,
	}
}

func (d *Dao[T]) Table() idao.Table {
	return newTable(d.driver, d.cfg.Schema)
}

func (d *Dao[T]) GetTableName() string {
	return d.cfg.Schema.Name
}
