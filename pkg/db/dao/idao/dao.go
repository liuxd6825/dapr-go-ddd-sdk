package idao

import (
	"context"

	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

type Dao[T any] interface {
	GetSchema() *store2.DBSchema
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
	UpdateMapByRSQL(ctx context.Context, rsql string, data map[string]any, opts ...CallOptions) *Result

	Delete(ctx context.Context, entity T, opts ...CallOptions) *Result
	DeleteMany(ctx context.Context, entity []T, opts ...CallOptions) *Result
	DeleteById(ctx context.Context, id string, opts ...CallOptions) *Result
	DeleteByIds(ctx context.Context, ids []string, opts ...CallOptions) *Result
	DeleteAll(ctx context.Context, opts ...CallOptions) *Result
	DeleteByRSQL(ctx context.Context, rsql string, opts ...CallOptions) *Result

	FindById(ctx context.Context, id string, opts ...CallOptions) (T, error)
	FindByIds(ctx context.Context, ids []string, opts ...CallOptions) ([]T, error)
	FindByRSQL(ctx context.Context, rsql string, opts ...CallOptions) ([]T, error)
	FindOneByRSQL(ctx context.Context, rsql string, opts ...CallOptions) (T, error)
	FindAll(ctx context.Context, opts ...CallOptions) *store2.FindListResult[T]
	FindPaging(ctx context.Context, qry store2.FindPagingQuery, opts ...CallOptions) store2.FindPagingResult[T]
	FindAutoComplete(ctx context.Context, qry store2.FindAutoCompleteQuery, opts ...CallOptions) store2.FindPagingResult[T]
	FindDistinct(ctx context.Context, qry store2.FindDistinctQuery, opts ...CallOptions) store2.FindPagingResult[T]

	SumByRSQL(ctx context.Context, rSql string, valueCols []*store2.ValueCol, opts ...CallOptions) (map[string]any, error)
	SumEntity(ctx context.Context, qry store2.FindPagingQuery, opts ...CallOptions) ([]T, error)
	SumByQuery(ctx context.Context, qry store2.FindPagingQuery, opts ...CallOptions) (map[string]any, error)

	CountByRSQL(ctx context.Context, rsql string, opts ...CallOptions) (int64, error)
	Table() Table

	GetStore() store2.IStore[T]
	GetDbType() string
	//GetFilterMap(tenantId string, rSql string) map[string]any
}

type CallOptions = store2.Options

func NewCallOptions(opts ...CallOptions) CallOptions {
	return store2.NewOptions(opts...)
}
