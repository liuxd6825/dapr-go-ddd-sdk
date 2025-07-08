package idao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
)

type Dao[T any] interface {
	GetSchema() *store.DBSchema
	GetAggField() string
	GetConfig() *DaoConfig

	Create(ctx context.Context, entity T, opts ...CallOptions) *Result
	CreateMany(ctx context.Context, entity []T, opts ...CallOptions) *Result
	CreateUpdate(ctx context.Context, entity T, opts ...CallOptions) *Result

	Merge(ctx context.Context, entity T, fields map[string]string, opts ...CallOptions) *Result

	Update(ctx context.Context, entity T, opts ...CallOptions) *Result
	UpdateMap(ctx context.Context, id string, entity map[string]any, opts ...CallOptions) *Result
	UpdateMany(ctx context.Context, entities []T, opts ...CallOptions) *Result
	UpdateByRSQL(ctx context.Context, rsql string, entity T, opts ...CallOptions) *Result

	Delete(ctx context.Context, entity T, opts ...CallOptions) *Result
	DeleteById(ctx context.Context, id string, opts ...CallOptions) *Result
	DeleteByIds(ctx context.Context, ids []string, opts ...CallOptions) *Result
	DeleteAll(ctx context.Context, opts ...CallOptions) *Result
	DeleteByRSQL(ctx context.Context, rsql string, opts ...CallOptions) *Result

	FindById(ctx context.Context, id string, opts ...CallOptions) (T, error)
	FindByIds(ctx context.Context, ids []string, opts ...CallOptions) ([]T, error)
	FindByRSQL(ctx context.Context, rsql string, opts ...CallOptions) ([]T, error)
	FindOneByRSQL(ctx context.Context, rsql string, opts ...CallOptions) (T, error)
	FindAll(ctx context.Context, opts ...CallOptions) *store.FindListResult[T]
	FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...CallOptions) store.FindPagingResult[T]
	FindAutoComplete(ctx context.Context, qry store.FindAutoCompleteQuery, opts ...CallOptions) store.FindPagingResult[T]
	FindDistinct(ctx context.Context, qry store.FindDistinctQuery, opts ...CallOptions) store.FindPagingResult[T]

	//SumByRSQL(ctx context.Context, rSql string, valueCols []*ddd_repository.ValueCol, opts ...CallOptions) T
	SumByRSQL(ctx context.Context, rSql string, valueCols []*store.ValueCol, opts ...CallOptions) (map[string]any, error)
	SumEntity(ctx context.Context, qry store.FindPagingQuery, opts ...CallOptions) ([]T, error)
	SumByQuery(ctx context.Context, qry store.FindPagingQuery, opts ...CallOptions) (map[string]any, error)

	CountByRSQL(ctx context.Context, rsql string, opts ...CallOptions) (int64, error)
	Table() Table

	GetStore() store.IStore[T]
	//GetFilterMap(tenantId string, rSql string) map[string]any
}

type CallOptions = store.Options

func NewCallOptions(opts ...CallOptions) CallOptions {
	return store.NewOptions(opts...)
}
