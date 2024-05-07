package db

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/schema"
	"go.mongodb.org/mongo-driver/mongo"
)

type IDao[T IEntity] interface {
	Save(ctx context.Context, setData *ddd.SetData[T], opts ...ddd_repository.Options) error
	Insert(ctx context.Context, entity T, opts ...ddd_repository.Options) error
	InsertByMap(ctx context.Context, tenantId string, data map[string]any, opts ...ddd_repository.Options) error
	InsertMany(ctx context.Context, entity []T, opts ...ddd_repository.Options) error
	Update(ctx context.Context, entity T, opts ...ddd_repository.Options) error
	UpdateByMap(ctx context.Context, tenantId string, filterMap map[string]any, data any, opts ...ddd_repository.Options) error
	UpdateMany(ctx context.Context, entities []T, opts ...ddd_repository.Options) error
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...ddd_repository.Options) (*ddd_repository.BulkWriteResult, error)
	UpdateManyByFilter(ctx context.Context, tenantId, filter string, data interface{}, opts ...ddd_repository.Options) error
	DeleteById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) error
	DeleteByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) error
	DeleteAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) error
	DeleteByFilter(ctx context.Context, tenantId string, filter string, opts ...ddd_repository.Options) error
	DeleteByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...ddd_repository.Options) error
	FindById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) (T, bool, error)
	FindByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) ([]T, bool, error)
	FindAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T]
	FindListByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T]
	FindPaging(ctx context.Context, query ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T]
	FindAutoComplete(ctx context.Context, qry ddd_repository.FindAutoCompleteQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T]
	FindDistinct(ctx context.Context, qry ddd_repository.FindDistinctQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T]
	AggregateByPipeline(ctx context.Context, pipeline mongo.Pipeline, data interface{}) error
	SumEntity(ctx context.Context, qry ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) ([]T, bool, error)
	SumMap(ctx context.Context, qry ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) ([]map[string]any, bool, error)
	Sum(ctx context.Context, qry ddd_repository.FindPagingQuery, data any, opts ...ddd_repository.Options) (any, bool, error)
	GetFilterMap(tenantId string, rsqlStr string) (map[string]interface{}, error)
}

type Dao struct {
	schema *schema.Schema
}

func NewDao(schema *schema.Schema) *Dao {
	return &Dao{schema: schema}
}
