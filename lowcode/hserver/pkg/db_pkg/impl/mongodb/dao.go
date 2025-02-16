package mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/dao/mongo_dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"go.mongodb.org/mongo-driver/mongo"
)

type Dao struct {
	*db.DaoBase
	db  *ddd_mongodb.MongoDB
	dao *ddd_mongodb.Dao[ddd.MapEntity] // 数据访问对象
	cfg *db.DaoConfig
}

type ModelOptions = mongo_dao.RepositoryOptions

type DaoOptions struct {
	DbKey      string
	TableName  string
	IsPubEvent bool
	AggField   string
	// mongo
	MongoDB         *ddd_mongodb.MongoDB
	GetCollCallback mongo_dao.GetCollectionCallback
	Server          element.Server
}

func NewDao(cfg *db.DaoConfig) *Dao {
	cfg.Valid()

	tableName := cfg.TableName
	item := restapp.GetDb(cfg.DbKey)
	if item == nil {
		panic(errors.New(" %s database not found", cfg.DbKey))
	}
	mongoDB := item.GetMongo()

	opt := mongo_dao.NewRepositoryOptions(&mongo_dao.RepositoryOptions{
		MongoDB: mongoDB,
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
	entBuilder := ddd.NewMapEntityBuilder[ddd.MapEntity]()
	daoOpts := ddd_mongodb.NewOptions[ddd.MapEntity]().SetAutoCreateCollection(true).SetAutoCreateIndex(true).SetEntityBuilder(entBuilder)

	dao := ddd_mongodb.NewDao[ddd.MapEntity](getCollCallback, daoOpts)
	res := &Dao{
		dao:     dao,
		DaoBase: db.NewDaoBase(dao, cfg),
		cfg:     cfg,
		db:      mongoDB,
	}
	return res
}

func (d *Dao) Table() db.Table {
	return NewTable(d.db, d.cfg.GetTableName(), d.cfg.Schema)
}
