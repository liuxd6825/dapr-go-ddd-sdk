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
	InsertMap(ctx context.Context, tenantId string, data map[string]interface{}, opts ...Options) error
	InsertMany(ctx context.Context, entities []T, opts ...Options) *SetManyResult[T]

	// 更新

	Update(ctx context.Context, entity T, opts ...Options) *SetResult[T]
	UpdateManyByFilter(ctx context.Context, tenantId, filter string, data any, opts ...Options) *SetManyCountResult
	UpdateManyById(ctx context.Context, entities []T, opts ...Options) *SetManyResult[T]
	UpdateManyMaskById(ctx context.Context, entities []T, mask []string, opts ...Options) *SetManyResult[T]
	UpdateMapById(ctx context.Context, tenantId string, id string, data map[string]any, opts ...Options) error
	FindOneAndUpdateById(ctx context.Context, tenantId string, id string, data map[string]any, opts ...Options) (T, error)
	UpdateMap(ctx context.Context, tenantId string, filter any, data any, opts ...Options) error
	UpdateMapAndGetCount(ctx context.Context, tenantId string, filter any, data any, opts ...Options) (int64, error)

	// 删除

	Delete(ctx context.Context, entity T, opts ...Options) *SetResult[T]
	DeleteByFilter(ctx context.Context, tenantId, filter string, opts ...Options) error
	DeleteById(ctx context.Context, tenantId string, id string, opts ...Options) *SetResult[T]
	DeleteByIds(ctx context.Context, tenantId string, ids []string, opts ...Options) error
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

	// 汇总

	SumEntity(ctx context.Context, qry FindPagingQuery, opts ...Options) ([]T, bool, error)
	SumMap(ctx context.Context, qry FindPagingQuery, opts ...Options) ([]map[string]any, bool, error)
	Sum(ctx context.Context, qry FindPagingQuery, resData any, opts ...Options) (any, bool, error)

	// 总数

	CountByMap(ctx context.Context, tenantId string, filterData any, opts ...Options) (int64, error)
	CountByRSQL(ctx context.Context, tenantId string, rsql string, opts ...Options) (int64, error)

	StartTx(ctx context.Context, fun TxFunc, options ...*SessionOptions) error

	SetMetadata(metadata map[string]any)
	GetMetadata() map[string]any
	AddMetadata(key string, val any)
}
