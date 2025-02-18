package db

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
)

type Dao interface {
	GetSchema() *schema.Schema
	GetAggField() string

	Create(ctx context.Context, entity map[string]any, opts ...*CallOptions)
	Update(ctx context.Context, entity map[string]any, opts ...*CallOptions)
	DeleteById(ctx context.Context, id string, opts ...*CallOptions)
	CreateMany(ctx context.Context, entity []map[string]any, opts ...*CallOptions)
	DeleteByIds(ctx context.Context, ids []string, opts ...*CallOptions)
	UpdateByMap(ctx context.Context, filterMap map[string]any, data map[string]any, opts ...*CallOptions)
	UpdateMany(ctx context.Context, entities []map[string]any, opts ...*CallOptions)
	UpdateManyByFilter(ctx context.Context, filter string, data interface{}, opts ...*CallOptions)
	DeleteAll(ctx context.Context, opts ...*CallOptions)
	DeleteByFilter(ctx context.Context, filter string, opts ...*CallOptions)
	DeleteByMap(ctx context.Context, filterMap map[string]any, opts ...*CallOptions)

	FindById(ctx context.Context, id string, opts ...*CallOptions) map[string]any
	FindByIds(ctx context.Context, ids []string, opts ...*CallOptions) []map[string]any
	FindAll(ctx context.Context, opts ...*CallOptions) *ddd_repository.FindListResult[map[string]any]
	FindListByMap(ctx context.Context, filterMap map[string]any, opts ...*CallOptions) *ddd_repository.FindListResult[map[string]any]
	FindPaging(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, opts ...*CallOptions) *ddd_repository.FindPagingResult[map[string]any]
	FindAutoComplete(ctx context.Context, qry *ddd_repository.FindAutoCompleteQueryRequest, opts ...*CallOptions) *ddd_repository.FindPagingResult[map[string]any]
	FindDistinct(ctx context.Context, qry *ddd_repository.FindDistinctQueryRequest, opts ...*CallOptions) *ddd_repository.FindPagingResult[map[string]any]

	SumEntity(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, opts ...*CallOptions) []map[string]any
	SumMap(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, opts ...*CallOptions) []map[string]any
	Sum(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, data any, opts ...*CallOptions) any

	Table() Table
	//GetFilterMap(tenantId string, rSql string) map[string]any
}

type CallOptions struct {
	AggId        *string
	EventType    *string
	EventVersion *string
	CommandId    *string
	ddd_repository.RepositoryOptions
}

func NewRepositoryOptions(opts []*CallOptions) []ddd_repository.Options {
	var res []ddd_repository.Options
	for _, o := range opts {
		res = append(res, o)
	}
	return res
}

func NewCallOptions(opts ...*CallOptions) *CallOptions {
	o := new(CallOptions)
	for _, i := range opts {
		if i.EventVersion != nil {
			o.EventVersion = i.EventVersion
		}
		if i.EventType != nil {
			o.EventType = i.EventType
		}
		if i.AggId != nil {
			o.AggId = i.AggId
		}
	}
	return o
}

func (e *CallOptions) GetAggId(defVal string) string {
	if e != nil && e.AggId == nil {
		return defVal
	}
	return *e.AggId
}

func (e *CallOptions) GetEventVersion(defVal string) string {
	if e != nil && e.EventVersion == nil {
		return defVal
	}
	return *e.EventVersion
}

func (e *CallOptions) GetCommandId(defVal string) string {
	if e != nil && e.CommandId != nil {
		return *e.CommandId
	}
	return defVal
}
