package ddd_repository

import "context"

type TxFunc func(ctx context.Context, options ...*SessionOptions) error

type Dao[T any] interface {
	NewEntity() T
	NewEntityList() []T

	GetTenantId(entity T) string
	SetTenantId(entity T, tenantId string)

	GetId(entity T) string
	SetId(entity T, id string)
	// 添加

	Insert(ctx context.Context, entity T, opts ...Options) *SetResult[T]
	InsertMap(ctx context.Context, tenantId string, data map[string]interface{}, opts ...Options) (res *SetResult[T])
	InsertMany(ctx context.Context, entities []T, opts ...Options) *SetManyResult[T]

	// 更新

	Update(ctx context.Context, entity T, opts ...Options) *SetResult[T]
	UpdateByRSQL(ctx context.Context, tenantId, filterRSQL string, data map[string]any, opts ...Options) *SetManyCountResult
	UpdateMany(ctx context.Context, entities []T, opts ...Options) *SetManyResult[T]
	UpdateMap(ctx context.Context, tenantId string, id string, data map[string]any, opts ...Options) *SetResult[T]
	UpdateMapAndGetCount(ctx context.Context, tenantId string, filter any, data any, opts ...Options) *SetResult[T]

	// 删除

	Delete(ctx context.Context, entity T, opts ...Options) *SetResult[T]
	DeleteByRSQL(ctx context.Context, tenantId, filter string, opts ...Options) *SetResult[T]
	DeleteById(ctx context.Context, tenantId string, id string, opts ...Options) *SetResult[T]
	DeleteByIds(ctx context.Context, tenantId string, ids []string, opts ...Options) *SetResult[T]
	DeleteAll(ctx context.Context, tenantId string, opts ...Options) *SetResult[T]
	DeleteByMap(ctx context.Context, tenantId string, filterMap map[string]any, opts ...Options) *SetResult[T]

	// 查询

	FindById(ctx context.Context, tenantId string, id string, opts ...Options) *FindOneResult[T]
	FindByIds(ctx context.Context, tenantId string, ids []string, opts ...Options) *FindListResult[T]
	FindOneByMap(ctx context.Context, tenantId string, filterMap map[string]any, opts ...Options) *FindOneResult[T]
	FindListByMap(ctx context.Context, tenantId string, filterMap map[string]any, opts ...Options) *FindListResult[T]
	FindByRSQL(ctx context.Context, tenantId string, rsql string, opts ...Options) *FindListResult[T]
	FindAll(ctx context.Context, tenantId string, opts ...Options) *FindListResult[T]
	FindPaging(ctx context.Context, qry FindPagingQuery, opts ...Options) (result *FindPagingResult[T])
	FindAutoComplete(ctx context.Context, qry FindAutoCompleteQuery, opts ...Options) *FindPagingResult[T]
	FindDistinct(ctx context.Context, qry FindDistinctQuery, opts ...Options) *FindPagingResult[T]
	FindOneAndUpdateById(ctx context.Context, tenantId string, id string, data map[string]any, opts ...Options) (T, error)
	// 汇总

	SumEntity(ctx context.Context, qry FindPagingQuery, opts ...Options) ([]T, bool, error)
	SumMap(ctx context.Context, qry FindPagingQuery, opts ...Options) ([]map[string]any, bool, error)
	Sum(ctx context.Context, qry FindPagingQuery, resData any, opts ...Options) (any, bool, error)
	SumByRSQL(ctx context.Context, tenantId, rSql string, valueCols []*ValueCol, opts ...Options) map[string]any
	// 总数

	CountByMap(ctx context.Context, tenantId string, filterData any, opts ...Options) (int64, error)
	CountByRSQL(ctx context.Context, tenantId string, rsql string, opts ...Options) (int64, error)

	StartTx(ctx context.Context, fun TxFunc, options ...*SessionOptions) error

	SetMetadata(metadata map[string]any)
	GetMetadata() map[string]any
	AddMetadata(key string, val any)
}
