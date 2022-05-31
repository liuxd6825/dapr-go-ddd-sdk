package ddd_repository

import (
	"context"
	"github.com/dapr/dapr-go-ddd-sdk/ddd"
)

//
// IRepository
// @Description: 仓储类接口
//
type IRepository[T ddd.Entity] interface {
	Insert(ctx context.Context, entity T, opts ...*SetOptions) *SetResult[T]
	Update(ctx context.Context, entity T, opts ...*SetOptions) *SetResult[T]

	Delete(ctx context.Context, entity T, opts ...*SetOptions) *SetResult[T]
	DeleteById(ctx context.Context, tenantId string, id string, opts ...*SetOptions) *SetResult[T]

	FindById(ctx context.Context, tenantId string, id string, opts ...*FindOptions) *FindOneResult[T]
	FindAll(ctx context.Context, tenantId string, opts ...*FindOptions) *FindListResult[T]
	FindPaging(ctx context.Context, search *FindPagingQuery, opts ...*FindOptions) *FindPagingResult[T]

	DoFindOne(fun func() (T, bool, error)) *FindOneResult[T]
	DoFindList(fun func() (*[]T, bool, error)) *FindListResult[T]
	DoSet(fun func() (T, error)) *SetResult[T]
	DoFilter(tenantId, filter string, fun func(filter map[string]interface{}) (*PagingData[T], bool, error)) *FindPagingResult[T]
}

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

type FindPagingQuery struct {
	TenantId string `json:"tenantId"`
	Fields   string `json:"fields"`
	// - name=="Kill Bill";year=gt=2003
	// - name=="Kill Bill" and year>2003
	// - genres=in=(sci-fi,action);(director=='Christopher Nolan',actor==*Bale);year=ge=2000
	// - genres=in=(sci-fi,action) and (director=='Christopher Nolan' or actor==*Bale) and year>=2000
	// - director.lastName==Nolan;year=ge=2000;year=lt=2010
	// - director.lastName==Nolan and year>=2000 and year<2010
	// - genres=in=(sci-fi,action);genres=out=(romance,animated,horror),director==Que*Tarantino
	// - genres=in=(sci-fi,action) and genres=out=(romance,animated,horror) or director==Que*Tarantino
	// or         : and ('OR' | 'or' and) *
	// and        : constraint ('AND' | 'and' constraint)*
	// constraint : group | comparison
	// group      : '(' or ')'
	// comparison : identifier comparator arguments
	// identifier : [a-zA-Z0-9]+('.'[a-zA-Z0-9]+)*
	// comparator : '==' | '!=' | '==~' | '!=~' | '>' | '>=' | '<' | '<=' | '=in=' | '=out='
	// arguments  : '(' listValue ')' | value
	// value      : int | double | string | date | datetime | boolean
	// listValue  : value(','value)*
	// int        : [0-9]+
	// double     : [0-9]+'.'[0-9]*
	// string     : '"'.*'"' | '\''.*'\''
	// date       : [0-9]{4}'-'[0-9]{2}'-'\[0-9]{2}
	// datetime   : date'T'[0-9]{2}':'[0-9]{2}':'[0-9]{2}('Z' | (('+'|'-')[0-9]{2}(':')?[0-9]{2}))?
	// boolean    : 'true' | 'false'
	Filter   string `json:"filter"`
	Sort     string `json:"sort"` // name:desc,id:asc
	PageNum  int64  `json:"pageNum"`
	PageSize int64  `json:"pageSize"`
}

type Option struct {
}
