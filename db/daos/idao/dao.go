package idao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
)

type Dao[T any] interface {
	GetSchema() *dbschema.DBSchema
	GetAggField() string

	Create(ctx context.Context, entity T, opts ...*CallOptions) *Result
	CreateMany(ctx context.Context, entity []T, opts ...*CallOptions) *Result

	Update(ctx context.Context, entity T, opts ...*CallOptions) *Result
	UpdateMap(ctx context.Context, id string, entity map[string]any, opts ...*CallOptions) *Result
	UpdateMany(ctx context.Context, entities []T, opts ...*CallOptions) *Result
	UpdateByRSQL(ctx context.Context, rsql string, entity T, opts ...*CallOptions) *Result

	Delete(ctx context.Context, entity T, opts ...*CallOptions) *Result
	DeleteById(ctx context.Context, id string, opts ...*CallOptions) *Result
	DeleteByIds(ctx context.Context, ids []string, opts ...*CallOptions) *Result
	DeleteAll(ctx context.Context, opts ...*CallOptions) *Result
	DeleteByRSQL(ctx context.Context, rsql string, opts ...*CallOptions) *Result

	FindById(ctx context.Context, id string, opts ...*CallOptions) T
	FindByIds(ctx context.Context, ids []string, opts ...*CallOptions) []T
	FindByRSQL(ctx context.Context, rsql string, opts ...*CallOptions) []T
	FindAll(ctx context.Context, opts ...*CallOptions) *ddd_repository.FindListResult[T]
	FindPaging(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, opts ...*CallOptions) *ddd_repository.FindPagingResult[T]
	FindAutoComplete(ctx context.Context, qry *ddd_repository.FindAutoCompleteQueryRequest, opts ...*CallOptions) *ddd_repository.FindPagingResult[T]
	FindDistinct(ctx context.Context, qry *ddd_repository.FindDistinctQueryRequest, opts ...*CallOptions) *ddd_repository.FindPagingResult[T]

	//SumByRSQL(ctx context.Context, rSql string, valueCols []*ddd_repository.ValueCol, opts ...*CallOptions) T
	SumByRSQL(ctx context.Context, rSql string, valueCols []*ddd_repository.ValueCol, opts ...*CallOptions) map[string]any
	SumEntity(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, opts ...*CallOptions) []T
	SumByQuery(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, opts ...*CallOptions) map[string]any

	CountByRSQL(ctx context.Context, rsql string, opts ...*CallOptions) int64
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
