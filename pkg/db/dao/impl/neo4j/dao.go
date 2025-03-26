package neo4j

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/store_neo4j"
	idao2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/impl"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Dao[T any] struct {
	*impl.DaoBase[T]
	cfg    *idao2.DaoConfig
	dao    store.IStore[T]
	driver neo4j.DriverWithContext
}

func NewDao[T any](cfg *idao2.DaoConfig, tableName ...string) idao2.Dao[T] {
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
		item := restapp.GetDB(cfg.DbKey)
		if item == nil {
			panic(fmt.Sprintf("dbKey %s not found", cfg.DbKey))
		}
		if val, ok := item.GetDB().(neo4j.DriverWithContext); ok {
			driver = val
		} else {
			panic(fmt.Sprintf("dbKey %s is not neo4j.DriverWithContext", cfg.DbKey))
		}
	}

	labels := []string{cfg.DBSchema.Name}
	if len(tableName) > 0 {
		labels = tableName
	}
	dbSch := cfg.DBSchema
	var storeImp store.IStore[T]
	if cfg.DaoType == "node" {
		storeImp = store_neo4j.NewNodeDao[T](driver, dbSch, labels)
	} else {
		storeImp = store_neo4j.NewRelationDao[T](driver, dbSch, labels)
	}
	daoBase := impl.NewDaoBase[T](storeImp, cfg)
	return &Dao[T]{
		DaoBase: daoBase,
		dao:     storeImp,
		cfg:     cfg,
		driver:  driver,
	}
}

func (d *Dao[T]) Table() idao2.Table {
	return newTable(d.driver, d.cfg.DBSchema)
}

func (d *Dao[T]) GetTableName() string {
	return d.cfg.DBSchema.Name
}
