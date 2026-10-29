package huge

import (
	"fmt"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
	idao2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/impl"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/store_huge"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
)

type Dao[T any] struct {
	*impl.DaoBase[T]
	cfg    *idao2.DaoConfig
	dao    store.IStore[T]
	client *hugegraph.CommonClient
}

func NewDao[T any](cfg *idao2.DaoConfig) idao2.Dao[T] {
	cfg.Valid()
	var client *hugegraph.CommonClient
	if cfg.DB != nil {
		if val, ok := cfg.DB.(*hugegraph.CommonClient); ok {
			client = val
		} else {
			panic("database config error hugegraph.CommonClient")
		}
	}
	envInst := cfg.Env
	if envInst == nil {
		envInst = envs.GetEnv()
	}

	if client == nil {
		item := envInst.GetDB(cfg.DBKey)
		if item == nil {
			panic(fmt.Sprintf("dbKey %s not found", cfg.DBKey))
		}
		if val := item.GetHugeClient(); val != nil {
			client = val
		} else {
			panic(fmt.Sprintf("dbKey %s is not hugegraph.CommonClient", cfg.DBKey))
		}
	}

	dbSch := cfg.DBSchema
	config := &store_huge.Config[T]{
		DBSchema: dbSch,
		Labels:   cfg.GraphLabels,
	}
	if config.Labels == nil && cfg.GraphType == idao2.GraphType_Node {
		panic(fmt.Sprintf("hugegraph dao %s config.Labels is required", dbSch.TableName))
	}

	var storeImp store.IStore[T] = nil
	if cfg.GraphType != "" {
		storeImp = store_huge.NewBaseDao(client, config)
	}
	daoBase := impl.NewDaoBase[T](storeImp, cfg)
	return &Dao[T]{
		DaoBase: daoBase,
		dao:     storeImp,
		cfg:     cfg,
		client:  client,
	}
}

func (d *Dao[T]) GetClient() *hugegraph.CommonClient {
	return d.client
}

func (d *Dao[T]) GetConfig() *idao2.DaoConfig {
	return d.cfg
}

func (d *Dao[T]) Table() idao2.Table {
	return newTable(d.client, d.cfg.DBSchema)
}

func (d *Dao[T]) GetHugeStore() *store_huge.BaseDao[T] {
	if d.dao == nil {
		return nil
	}
	v, ok := d.dao.(*store_huge.BaseDao[T])
	if !ok {
		panic("hugegraph store does not implement store_huge.Dao")
	}
	return v
}
