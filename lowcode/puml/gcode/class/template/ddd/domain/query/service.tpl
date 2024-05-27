package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_neo4j"
	"github.com/liuxu6825/duxm-graph-service/pkg/query-service/domain/graph/query"
	"github.com/liuxu6825/duxm-graph-service/pkg/query-service/domain/graph/view"
)

type AccountNodeQueryDomainService interface {
	ImportCsvFile(ctx context.Context, tenantId, caseId, importFile string, callback ddd_neo4j.ImportSaveFileCallback, data []*view.AccountView, opts ...Options) error
	Save(ctx context.Context, data *ddd.SetData[*view.AccountView], opts ...Options) error

	Create(ctx context.Context, view *view.AccountView, opts ...Options) error
	CreateMany(ctx context.Context, views []*view.AccountView, opts ...Options) error

	Update(ctx context.Context, view *view.AccountView, opts ...Options) error
	UpdateMany(ctx context.Context, views []*view.AccountView, opts ...Options) error

	Delete(ctx context.Context, view *view.AccountView, opts ...Options) error
	DeleteMany(ctx context.Context, tenantId string, views []*view.AccountView, opts ...Options) error
	DeleteById(ctx context.Context, tenantId string, id string, opts ...Options) error
	DeleteByIds(ctx context.Context, tenantId string, ids []string, opts ...Options) error
	DeleteByFilter(ctx context.Context, tenantId, filter string, opts ...Options) error
	DeleteAll(ctx context.Context, tenantId string, opts ...Options) error
	DeleteByCaseId(ctx context.Context, tenantId string, caseId string, opts ...Options) error

	FindById(ctx context.Context, qry *query.AccountNodeFindByIdQuery, opts ...Options) (*view.AccountView, bool, error)
	FindByIds(ctx context.Context, qry *query.AccountNodeFindByIdsQuery, opts ...Options) ([]*view.AccountView, bool, error)
	FindAll(ctx context.Context, qry *query.AccountNodeFindAllQuery, opts ...Options) ([]*view.AccountView, bool, error)
	FindPaging(ctx context.Context, qry *query.AccountNodeFindPagingQuery, opts ...Options) (*query.AccountNodeFindPagingResult, bool, error)
	FindByCaseId(ctx context.Context, qry *query.AccountNodeFindByCaseIdQuery, opts ...Options) ([]*view.AccountView, bool, error)
}
