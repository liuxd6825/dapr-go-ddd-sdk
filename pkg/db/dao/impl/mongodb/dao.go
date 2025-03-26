package mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/store_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	idao2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/impl"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetCollectionCallback func(ctx context.Context) (*store_mongodb.MongoDB, *mongo.Collection)

type Dao[T any] struct {
	*impl.DaoBase[T]
	db  *store_mongodb.MongoDB
	dao store.IStore[T] // 数据访问对象
	cfg *idao2.DaoConfig
}

var _mongodb *store_mongodb.MongoDB

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

func NewDao[T any](cfg *idao2.DaoConfig, tableNames ...string) idao2.Dao[T] {
	cfg.Valid()
	var mongoDb *store_mongodb.MongoDB
	if v, ok := cfg.Database.(*store_mongodb.MongoDB); ok {
		mongoDb = v
	} else {
		item := restapp.GetDB(cfg.DbKey)
		if item == nil {
			panic(errors.New(" %s database not found", cfg.DbKey))
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

	var mongodb *store_mongodb.MongoDB
	var coll *mongo.Collection

	getCollCallback := func(ctx context.Context) (*store_mongodb.MongoDB, *mongo.Collection) {
		if mongodb == nil || coll == nil {
			mongodb = opt.MongoDB
			coll = opt.MongoDB.GetCollection(tableName)
		}
		return mongodb, coll
	}

	if opt.GetCollCallback != nil {
		getCollCallback = opt.GetCollCallback
	}
	eb := store.NewAnyEntityBuilder[T](cfg.DBSchema)
	daoOpts := store_mongodb.NewOptions[T]().SetAutoCreateCollection(true).SetAutoCreateIndex(true).SetEntityBuilder(eb)

	dao := store_mongodb.NewDao[T](cfg.DBSchema, getCollCallback, daoOpts)
	res := &Dao[T]{
		dao:     dao,
		DaoBase: impl.NewDaoBase[T](dao, cfg),
		cfg:     cfg,
		db:      mongoDb,
	}
	return res
}

func (d *Dao[T]) Table() idao2.Table {
	return NewTable(d.db, d.cfg.DBSchema)
}

func (d *Dao[T]) GetTableName() string {
	return d.cfg.DBSchema.TableName
}
