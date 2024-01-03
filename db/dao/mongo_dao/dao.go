package mongo_dao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"go.mongodb.org/mongo-driver/mongo"
)

type Dao[T ddd.Entity] struct {
	dao *ddd_mongodb.Dao[T]
}

type RepositoryOptions struct {
	mongoDB         *ddd_mongodb.MongoDB
	getCollCallback GetCollectionCallback
}

type GetCollectionCallback func(ctx context.Context) (*ddd_mongodb.MongoDB, *mongo.Collection)

var _mongodb *ddd_mongodb.MongoDB

func NewSession(isWrite bool) ddd_repository.Session {
	return ddd_mongodb.NewSession(isWrite, GetDB())
}

func NewRepositoryOptions(opts ...*RepositoryOptions) *RepositoryOptions {
	o := &RepositoryOptions{}
	for _, item := range opts {
		if item.mongoDB != nil {
			o.mongoDB = item.mongoDB
		}
		if item.getCollCallback != nil {
			o.getCollCallback = item.getCollCallback
		}
	}
	if o.mongoDB == nil {
		o.mongoDB = _mongodb
	}
	if o.mongoDB == nil {
		o.mongoDB = GetDB()
	}
	return o
}

func NewDao[T ddd.Entity](collectionName string, opts ...*RepositoryOptions) *Dao[T] {
	initCollName := collectionName
	opt := NewRepositoryOptions(opts...)
	var mongodb *ddd_mongodb.MongoDB
	var coll *mongo.Collection

	getCollCallback := func(ctx context.Context) (*ddd_mongodb.MongoDB, *mongo.Collection) {
		if mongodb == nil || coll == nil {
			mongodb = opt.mongoDB
			coll = opt.mongoDB.GetCollection(initCollName)
		}
		return mongodb, coll
	}

	if opt.getCollCallback != nil {
		getCollCallback = opt.getCollCallback
	}

	return &Dao[T]{
		dao: ddd_mongodb.NewDao[T](getCollCallback, ddd_mongodb.NewOptions().SetAutoCreateCollection(true).SetAutoCreateIndex(true)),
	}
}

func (d *Dao[T]) Save(ctx context.Context, setData *ddd.SetData[T], opts ...ddd_repository.Options) error {
	return d.dao.Save(ctx, setData, opts...).GetError()
}

func (d *Dao[T]) Insert(ctx context.Context, entity T, opts ...ddd_repository.Options) error {
	return d.dao.Insert(ctx, entity, opts...).GetError()
}

func (d *Dao[T]) InsertByMap(ctx context.Context, tenantId string, data map[string]any, opts ...ddd_repository.Options) error {
	return d.dao.InsertMap(ctx, tenantId, data, opts...)
}

func (d *Dao[T]) InsertMany(ctx context.Context, entity []T, opts ...ddd_repository.Options) error {
	return d.dao.InsertMany(ctx, entity, opts...).GetError()
}

func (d *Dao[T]) Update(ctx context.Context, entity T, opts ...ddd_repository.Options) error {
	return d.dao.Update(ctx, entity, opts...).GetError()
}

func (d *Dao[T]) UpdateByMap(ctx context.Context, tenantId string, filterMap map[string]any, data any, opts ...ddd_repository.Options) error {
	return d.dao.UpdateMap(ctx, tenantId, filterMap, data, opts...)
}

func (d *Dao[T]) UpdateMany(ctx context.Context, entities []T, opts ...ddd_repository.Options) error {
	return d.dao.UpdateManyById(ctx, entities, opts...).GetError()
}

func (d *Dao[T]) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...ddd_repository.Options) (*ddd_repository.BulkWriteResult, error) {
	return d.dao.BulkWrite(ctx, models, opts...)
}

func (d *Dao[T]) UpdateManyByFilter(ctx context.Context, tenantId, filter string, data interface{}, opts ...ddd_repository.Options) error {
	return d.dao.UpdateManyByFilter(ctx, tenantId, filter, data, opts...).GetError()
}

func (d *Dao[T]) DeleteById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) error {
	return d.dao.DeleteById(ctx, tenantId, id, opts...).GetError()
}

func (d *Dao[T]) DeleteByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) error {
	return d.dao.DeleteByIds(ctx, tenantId, ids, opts...)
}

func (d *Dao[T]) DeleteAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) error {
	return d.dao.DeleteAll(ctx, tenantId, opts...).GetError()
}

func (d *Dao[T]) DeleteByFilter(ctx context.Context, tenantId string, filter string, opts ...ddd_repository.Options) error {
	return d.dao.DeleteByFilter(ctx, tenantId, filter, opts...)
}

func (d *Dao[T]) DeleteByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...ddd_repository.Options) error {
	return d.dao.DeleteByMap(ctx, tenantId, filterMap, opts...).GetError()
}

func (d *Dao[T]) FindById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) (T, bool, error) {
	return d.dao.FindById(ctx, tenantId, id, opts...).Result()
}

func (d *Dao[T]) FindByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) ([]T, bool, error) {
	return d.dao.FindByIds(ctx, tenantId, ids, opts...).Result()
}

func (d *Dao[T]) FindAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	return d.dao.FindAll(ctx, tenantId, opts...)
}

func (d *Dao[T]) FindListByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	return d.dao.FindListByMap(ctx, tenantId, filterMap, opts...)
}

func (d *Dao[T]) FindPaging(ctx context.Context, query ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	return d.dao.FindPaging(ctx, query, opts...)
}

func (d *Dao[T]) FindAutoComplete(ctx context.Context, qry ddd_repository.FindAutoCompleteQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	return d.dao.FindAutoComplete(ctx, qry, opts...)
}

func (d *Dao[T]) FindDistinct(ctx context.Context, qry ddd_repository.FindDistinctQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	return d.dao.FindDistinct(ctx, qry, opts...)
}

func (d *Dao[T]) AggregateByPipeline(ctx context.Context, pipeline mongo.Pipeline, data interface{}) error {
	return d.dao.AggregateByPipeline(ctx, pipeline, data)
}

func (d *Dao[T]) SumEntity(ctx context.Context, qry ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) ([]T, bool, error) {
	return d.dao.SumEntity(ctx, qry, opts...)
}

func (d *Dao[T]) SumMap(ctx context.Context, qry ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) ([]map[string]any, bool, error) {
	return d.dao.SumMap(ctx, qry, opts...)
}

func (d *Dao[T]) Sum(ctx context.Context, qry ddd_repository.FindPagingQuery, data any, opts ...ddd_repository.Options) (any, bool, error) {
	return d.dao.Sum(ctx, qry, data, opts...)
}

func (d *Dao[T]) GetFilterMap(tenantId string, rsqlstr string) (map[string]interface{}, error) {
	return d.dao.GetFilterMap(tenantId, rsqlstr)
}

func (o *RepositoryOptions) SetMongoDB(mongoDB *ddd_mongodb.MongoDB) *RepositoryOptions {
	o.mongoDB = mongoDB
	return o
}

func (o *RepositoryOptions) SetGetCollectionCallback(v GetCollectionCallback) *RepositoryOptions {
	o.getCollCallback = v
	return o
}

func GetDB() *ddd_mongodb.MongoDB {
	if _mongodb != nil {
		return _mongodb
	}
	return restapp.GetMongoDB()
}

func SetDB(mongodb *ddd_mongodb.MongoDB) {
	_mongodb = mongodb
}
