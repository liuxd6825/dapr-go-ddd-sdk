package ddd_repository

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"go.mongodb.org/mongo-driver/bson"
)

type Repository[T ddd.Entity] interface {
	Insert(ctx context.Context, entity T, opts ...*SetOptions) *SetResult[T]
	Update(ctx context.Context, entity T, opts ...*SetOptions) *SetResult[T]
	Delete(ctx context.Context, entity ddd.Entity, opts ...*SetOptions) *SetResult[T]
	DeleteById(ctx context.Context, tenantId string, id string, opts ...*SetOptions) *SetResult[T]
	DeleteAll(ctx context.Context, tenantId string, opts ...*SetOptions) *SetResult[T]
	DeleteByMap(ctx context.Context, tenantId string, data map[string]interface{}, opts ...*SetOptions) *SetResult[T]
	NewFilter(tenantId string, filterMap map[string]interface{}) bson.D
	FindById(ctx context.Context, tenantId string, id string, opts ...*FindOptions) *FindOneResult[T]
	FindOneByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...*FindOptions) *FindOneResult[T]
	FindListByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...*FindOptions) *FindListResult[T]
	FindAll(ctx context.Context, tenantId string, opts ...*FindOptions) *FindListResult[T]
	FindPaging(ctx context.Context, query FindPagingQuery, opts ...*FindOptions) *FindPagingResult[T]
	DoFilter(tenantId string, filter string, fun func(filter map[string]interface{}) (*FindPagingResult[T], bool, error)) *FindPagingResult[T]
	DoFindList(fun func() (*[]T, bool, error)) *FindListResult[T]
	DoFindOne(fun func() (T, bool, error)) *FindOneResult[T]
	DoSet(fun func() (T, error)) *SetResult[T]
}

//
// IRepository
// @Description: 仓储类接口
//
/*
type IRepository[T ddd.Entity] interface {
	Insert(ctx context.Context, entity T, opts ...*SetOptions) *SetResult[T]
	Update(ctx context.Context, entity T, opts ...*SetOptions) *SetResult[T]

	Delete(ctx context.Context, entity T, opts ...*SetOptions) *SetResult[T]
	DeleteById(ctx context.Context, tenantId string, id string, opts ...*SetOptions) *SetResult[T]

	FindById(ctx context.Context, tenantId string, id string, opts ...*FindOptions) *FindOneResult[T]
	FindAll(ctx context.Context, tenantId string, opts ...*FindOptions) *FindListResult[T]
	FindPaging(ctx context.Context, search FindPagingQuery, opts ...*FindOptions) *FindPagingResult[T]

	DoFindOne(fun func() (T, bool, error)) *FindOneResult[T]
	DoFindList(fun func() (*[]T, bool, error)) *FindListResult[T]
	DoSet(fun func() (T, error)) *SetResult[T]
	DoFilter(tenantId, filter string, fun func(filter map[string]interface{}) (*PagingData[T], bool, error)) *FindPagingResult[T]
}
*/
//
// Repository
// @Description: 仓储类，
//
/*type Repository[T ddd.Entity] struct {
	mongo *ddd_mongodb.Repository[T]
}

func NewRepositoryWithMongo[T ddd.Entity](entityBuilder *EntityBuilder[T], mongodb *ddd_mongodb.MongoDB, collection *mongo.Collection) *Repository[T] {
	return &Repository[T]{
		entityBuilder: entityBuilder,
		collection:    collection,
		mongodb:       mongodb,
	}
}

func (r *Repository[T]) Insert(ctx context.Context, entity T, opts ...*SetOptions) *SetResult[T] {
	return r.mongo.Insert(ctx, entity, opts...)
}

func (r *Repository[T]) Update(ctx context.Context, entity T, opts ...*SetOptions) *SetResult[T] {
	return r.mongo.Update(ctx, entity, opts...)
}

func (r *Repository[T]) Delete(ctx context.Context, entity T, opts ...*SetOptions) *SetResult[T] {
	return r.mongo.Delete(ctx, entity, opts...)
}

func (r *Repository[T]) DeleteById(ctx context.Context, tenantId string, id string, opts ...*SetOptions) *SetResult[T] {
	return r.mongo.DeleteById(ctx, tenantId, id, opts...)
}

func (r *Repository[T]) FindById(ctx context.Context, tenantId string, id string, opts ...*FindOptions) *FindOneResult[T] {
	return r.mongo.FindById(ctx, tenantId, id, opts...)
}

func (r *Repository[T]) FindAll(ctx context.Context, tenantId string, opts ...*FindOptions) *FindOneResult[*[]T] {
	return r.mongo.FindAll(ctx, tenantId, opts...)
}

func (r *Repository[T]) FindPaging(ctx context.Context, search *PagingQuery, opts ...*FindOptions) *FindPagingResult[T] {
	return r.mongo.FindPaging(ctx, search, opts...)
}

func (r *Repository[T]) DoFindOne(fun func() (T, bool, error)) *FindOneResult[T] {
	return r.mongo.DoFindOne(fun)
}

func (r *Repository[T]) DoFindList(fun func() (*[]T, bool, error)) *FindOneResult[*[]T] {
	return r.mongo.DoFindList(fun)
}

func (r *Repository[T]) DoSet(fun func() (T, error)) *SetResult[T] {
	return r.mongo.DoSet(fun)
}

func (r *Repository[T]) DoFilter(tenantId, filter string, fun func(filter map[string]interface{}) (*[]T, bool, error)) *FindPagingResult[T] {
	return r.mongo.DoFilter(tenantId, filter, fun)
}
*/
/*type IRepos interface {
	Insert(interface{})
}

type Repos[T interface{}] struct {
}

func NewRepos[T]() IRepos {
	return &Repos{}
}

func (r Repos[T]) Insert(t T) {
	panic("implement me")
}*/

type Pageable struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

type Option struct {
}
