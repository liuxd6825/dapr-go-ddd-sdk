package mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/impl"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetCollectionCallback func(ctx context.Context) (*ddd_mongodb.MongoDB, *mongo.Collection)

type Dao struct {
	*impl.DaoBase[map[string]any]
	db  *ddd_mongodb.MongoDB
	dao ddd_repository.Dao[map[string]any] // 数据访问对象
	cfg *db.DaoConfig
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

func NewDao(cfg *db.DaoConfig, tableNames ...string) *Dao {
	cfg.Valid()
	var mongoDb *ddd_mongodb.MongoDB
	if v, ok := cfg.Database.(*ddd_mongodb.MongoDB); ok {
		mongoDb = v
	} else {
		item := restapp.GetDb(cfg.DbKey)
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
	eb := ddd.NewAnyEntityBuilder[map[string]any]()
	daoOpts := ddd_mongodb.NewOptions[map[string]any]().SetAutoCreateCollection(true).SetAutoCreateIndex(true).SetEntityBuilder(eb)

	dao := ddd_mongodb.NewDao[map[string]any](getCollCallback, daoOpts)
	res := &Dao{
		dao:     dao,
		DaoBase: impl.NewDaoBase(dao, cfg),
		cfg:     cfg,
		db:      mongoDb,
	}
	return res
}

func (d *Dao) Table() db.Table {
	return NewTable(d.db, d.cfg.Schema)
}

func (d *Dao) GetTableName() string {
	return d.cfg.Schema.Name
}
