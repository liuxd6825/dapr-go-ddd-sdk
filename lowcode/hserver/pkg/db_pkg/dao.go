package db_pkg

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
)

type Dao interface {
	GetAggField() string
	Create(ctx context.Context, entity ddd.MapEntity, opts ...*DaoOptions)
	Update(ctx context.Context, entity ddd.MapEntity, opts ...*DaoOptions)
	DeleteById(ctx context.Context, id string, opts ...*DaoOptions)
	CreateMany(ctx context.Context, entity []ddd.MapEntity, opts ...*DaoOptions)
	DeleteByIds(ctx context.Context, ids []string, opts ...*DaoOptions)
	UpdateByMap(ctx context.Context, filterMap map[string]any, data map[string]any, opts ...*DaoOptions)
	UpdateMany(ctx context.Context, entities []ddd.MapEntity, opts ...*DaoOptions)
	UpdateManyByFilter(ctx context.Context, filter string, data interface{}, opts ...*DaoOptions)
	DeleteAll(ctx context.Context, opts ...*DaoOptions)
	DeleteByFilter(ctx context.Context, filter string, opts ...*DaoOptions)
	DeleteByMap(ctx context.Context, filterMap map[string]interface{}, opts ...*DaoOptions)

	FindById(ctx context.Context, id string, opts ...*DaoOptions) ddd.MapEntity
	FindByIds(ctx context.Context, ids []string, opts ...*DaoOptions) []ddd.MapEntity
	FindAll(ctx context.Context, opts ...*DaoOptions) *ddd_repository.FindListResult[ddd.MapEntity]
	FindListByMap(ctx context.Context, filterMap map[string]interface{}, opts ...*DaoOptions) *ddd_repository.FindListResult[ddd.MapEntity]
	FindPaging(ctx context.Context, findPagingMap any, opts ...*DaoOptions) *ddd_repository.FindPagingResult[ddd.MapEntity]
	FindAutoComplete(ctx context.Context, qry *ddd_repository.FindAutoCompleteQueryRequest, opts ...*DaoOptions) *ddd_repository.FindPagingResult[ddd.MapEntity]
	FindDistinct(ctx context.Context, qry *ddd_repository.FindDistinctQueryRequest, opts ...*DaoOptions) *ddd_repository.FindPagingResult[ddd.MapEntity]

	SumEntity(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, opts ...*DaoOptions) []ddd.MapEntity
	SumMap(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, opts ...*DaoOptions) []map[string]any
	Sum(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, data any, opts ...*DaoOptions) any

	GetFilterMap(tenantId string, rSql string) ddd.MapEntity
}

type NewDaoOptions struct {
	DbName     string `json:"dbName"`
	TableName  string `json:"tableName"`
	IsPubEvent bool   `json:"isPubEvent"`
	AggField   string `json:"aggField"`
}

type DaoOptions struct {
	AggId        *string
	EventType    *string
	EventVersion *string
	CommandId    *string
	ddd_repository.RepositoryOptions
}

func NewRepositoryOptions(opts []*DaoOptions) []ddd_repository.Options {
	var res []ddd_repository.Options
	for _, o := range opts {
		res = append(res, o)
	}
	return res
}

func NewOptions(opts ...*DaoOptions) *DaoOptions {
	o := new(DaoOptions)
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

func (e *DaoOptions) GetAggId(defVal string) string {
	if e != nil && e.AggId == nil {
		return defVal
	}
	return *e.AggId
}

func (e *DaoOptions) GetEventVersion(defVal string) string {
	if e != nil && e.EventVersion == nil {
		return defVal
	}
	return *e.EventVersion
}

func (e *DaoOptions) GetCommandId(defVal string) string {
	if e != nil && e.CommandId != nil {
		return *e.CommandId
	}
	return defVal
}
