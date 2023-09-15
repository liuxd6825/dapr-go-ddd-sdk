package ddd_neo4j

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
)

type Cypher interface {
	Insert(ctx context.Context, data interface{}) (CypherResult, error)
	InsertMany(ctx context.Context, list interface{}) (CypherResult, error)

	InsertOrUpdate(ctx context.Context, data interface{}) (CypherResult, error)

	Update(ctx context.Context, data interface{}, setFields ...string) (CypherResult, error)
	UpdateMany(ctx context.Context, list interface{}) (CypherResult, error)

	UpdateLabelById(ctx context.Context, tenantId string, id string, label string) (CypherResult, error)
	UpdateLabelByFilter(ctx context.Context, tenantId string, filter string, labels ...string) (CypherResult, error)

	DeleteById(ctx context.Context, tenantId string, id string) (CypherResult, error)
	DeleteByIds(ctx context.Context, tenantId string, ids []string) (CypherResult, error)
	DeleteAll(ctx context.Context, tenantId string) (CypherResult, error)
	DeleteByFilter(ctx context.Context, tenantId string, filter string) (CypherResult, error)
	DeleteByTenantId(ctx context.Context, tenantId string) (CypherResult, error)
	DeleteByLabels(ctx context.Context, tenantId string, label ...string) (CypherResult, error)
	
	DeleteLabelById(ctx context.Context, tenantId string, id string, label string) (CypherResult, error)
	DeleteLabelByFilter(ctx context.Context, tenantId string, filter string, labels ...string) (CypherResult, error)

	GetFilter(ctx context.Context, tenantId, filter string) (CypherResult, error)
	FindById(ctx context.Context, tenantId, id string) (CypherResult, error)
	FindByIds(ctx context.Context, tenantId string, ids []string) (CypherResult, error)
	FindByAggregateId(ctx context.Context, tenantId, aggregateName, aggregateId string) (result CypherResult, err error)
	FindByCaseId(ctx context.Context, tenantId string, caseId string) (result CypherResult, err error)
	FindAll(ctx context.Context, tenantId string) (CypherResult, error)
	FindPaging(ctx context.Context, query ddd_repository.FindPagingQuery) (CypherResult, error)
	Count(ctx context.Context, tenantId, filter string) (CypherResult, error)
}
