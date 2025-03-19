package mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos/impl"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetCollectionCallback func(ctx context.Context) (*ddd_mongodb.MongoDB, *mongo.Collection)

type Dao[T any] struct {
	*impl.DaoBase[T]
	db  *ddd_mongodb.MongoDB
	dao ddd_repository.Dao[T] // 数据访问对象
	cfg *idao.DaoConfig
}

var _mongodb *ddd_mongodb.MongoDB

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
	MongoDB         ddd_repository.Dao[map[string]any]
	GetCollCallback GetCollectionCallback
	Server          element.Server
}

func NewDao[T any](cfg *idao.DaoConfig, tableNames ...string) idao.Dao[T] {
	cfg.Valid()
	var mongoDb *ddd_mongodb.MongoDB
	if v, ok := cfg.Database.(*ddd_mongodb.MongoDB); ok {
		mongoDb = v
	} else {
		item := restapp.GetDB(cfg.DbKey)
		if item == nil {
			panic(errors.New(" %s database not found", cfg.DbKey))
		}
		mongoDb = item.GetMongo()
	}

	tableName := cfg.Schema.Name
	if len(tableNames) > 0 {
		tableName = tableNames[0]
	}
	opt := NewRepositoryOptions(&RepositoryOptions{
		MongoDB: mongoDb,
		//GetCollCallback: opts.GetCollCallback,
	})

	var mongodb *ddd_mongodb.MongoDB
	var coll *mongo.Collection

	getCollCallback := func(ctx context.Context) (*ddd_mongodb.MongoDB, *mongo.Collection) {
		if mongodb == nil || coll == nil {
			mongodb = opt.MongoDB
			coll = opt.MongoDB.GetCollection(tableName)
		}
		return mongodb, coll
	}

	if opt.GetCollCallback != nil {
		getCollCallback = opt.GetCollCallback
	}
	eb := ddd.NewAnyEntityBuilder[T](cfg.Schema)
	daoOpts := ddd_mongodb.NewOptions[T]().SetAutoCreateCollection(true).SetAutoCreateIndex(true).SetEntityBuilder(eb)

	dao := ddd_mongodb.NewDao[T](cfg.Schema, getCollCallback, daoOpts)
	res := &Dao[T]{
		dao:     dao,
		DaoBase: impl.NewDaoBase[T](dao, cfg),
		cfg:     cfg,
		db:      mongoDb,
	}
	return res
}

func (d *Dao[T]) Table() idao.Table {
	return NewTable(d.db, d.cfg.Schema)
}

func (d *Dao[T]) GetTableName() string {
	return d.cfg.Schema.TableName
}
