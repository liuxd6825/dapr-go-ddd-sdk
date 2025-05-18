package mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/store_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/impl"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetCollectionCallback func(ctx context.Context) (store_mongodb.IMongoDB, *mongo.Collection)

type Dao[T any] struct {
	*impl.DaoBase[T]
	db    store_mongodb.IMongoDB
	store store.IStore[T] // 数据访问对象
	cfg   *idao.DaoConfig
}

var _mongodb store_mongodb.IMongoDB

func init() {
	//设置bson使用自定义的日期json格式
	primitive.UseCustomTimeFormat = true
}

type ModelOptions = RepositoryOptions

type DaoOptions struct {
	DbKey      string
	TableName  string
	IsPubEvent bool
	AggField   string
	// mongo
	MongoDB         store.IStore[map[string]any]
	GetCollCallback GetCollectionCallback
	Server          element.Server
}

func NewDao[T any](cfg *idao.DaoConfig, tableNames ...string) idao.Dao[T] {
	cfg.Valid()
	var mongoDb store_mongodb.IMongoDB
	if v, ok := cfg.DB.(store_mongodb.IMongoDB); ok {
		mongoDb = v
	} else {
		item := env.GetDB(cfg.DBKey)
		if item == nil {
			panic(errors.New(" %s database not found", cfg.DBKey))
		}
		mongoDb = item.GetMongo()
	}

	tableName := cfg.DBSchema.Name
	if len(tableNames) > 0 {
		tableName = tableNames[0]
	}
	opt := NewRepositoryOptions(&RepositoryOptions{
		MongoDB: mongoDb,
		//GetCollCallback: opts.GetCollCallback,
	})

	var mongodb store_mongodb.IMongoDB
	var coll *mongo.Collection

	getCollCallback := func(ctx context.Context) (store_mongodb.IMongoDB, *mongo.Collection) {
		if mongodb == nil || coll == nil {
			mongodb = opt.MongoDB
			coll = opt.MongoDB.GetDatabase().Collection(tableName)
		}
		return mongodb, coll
	}

	if opt.GetCollCallback != nil {
		getCollCallback = opt.GetCollCallback
	}
	eb := store.NewAnyEntityBuilder[T](cfg.DBSchema)
	daoOpts := store_mongodb.NewOptions[T]().SetAutoCreateCollection(true).SetAutoCreateIndex(true).SetEntityBuilder(eb)
	mongoStore := store_mongodb.NewDao[T](cfg.DBSchema, getCollCallback, daoOpts)

	res := &Dao[T]{
		DaoBase: impl.NewDaoBase[T](mongoStore, cfg),
		store:   mongoStore,
		cfg:     cfg,
		db:      mongoDb,
	}
	return res
}

func (d *Dao[T]) Table() idao.Table {
	return NewTable(d.db, d.cfg.DBSchema)
}

func (d *Dao[T]) GetTableName() string {
	return d.cfg.DBSchema.TableName
}

func (d *Dao[T]) GetConfig() *idao.DaoConfig {
	return d.cfg
}
