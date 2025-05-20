package neo4j

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/store_neo4j"
	idao2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/impl"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Dao[T any] struct {
	*impl.DaoBase[T]
	cfg    *idao2.DaoConfig
	dao    store.IStore[T]
	driver neo4j.DriverWithContext
}

func NewDao[T any](cfg *idao2.DaoConfig) idao2.Dao[T] {
	cfg.Valid()
	var driver neo4j.DriverWithContext
	//eb := ddd.NewMapEntityBuilder[map[string]any]()
	if cfg.DB != nil {
		if val, ok := cfg.DB.(neo4j.DriverWithContext); ok {
			driver = val
		} else {
			panic("database config error neo4j.DriverWithContext")
		}
	}
	envInst := cfg.Env
	if envInst == nil {
		envInst = env.GetEnv()
	}

	if driver == nil {
		item := envInst.GetDB(cfg.DBKey)
		if item == nil {
			panic(fmt.Sprintf("dbKey %s not found", cfg.DBKey))
		}
		if val, ok := item.GetDB().(neo4j.DriverWithContext); ok {
			driver = val
		} else {
			panic(fmt.Sprintf("dbKey %s is not neo4j.DriverWithContext", cfg.DBKey))
		}
	}

	dbSch := cfg.DBSchema
	config := &store_neo4j.Config[T]{
		DBSchema: dbSch,
		Labels:   cfg.GraphLabels,
	}
	if config.Labels == nil && cfg.GraphType == idao2.GraphType_Node {
		panic(fmt.Sprintf("neo4j dao %s config.Labels is be must nil", dbSch.TableName))
	}

	var storeImp store.IStore[T]
	if cfg.GraphType == "node" {
		storeImp = store_neo4j.NewNodeDao[T](driver, config)
	} else if cfg.GraphType == "rel" {
		storeImp = store_neo4j.NewRelationDao[T](driver, config)
	} else if cfg.GraphType != "" {
		panic(fmt.Sprintf("NewDao() error : invalid graphType %s", cfg.GraphType))
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

func (d *Dao[T]) GetConfig() *idao2.DaoConfig {
	return d.cfg
}
