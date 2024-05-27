package repository

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_neo4j"
	"github.com/liuxu6825/duxm-graph-service/pkg/query-service/domain/graph/service"
	"github.com/liuxu6825/duxm-graph-service/pkg/query-service/domain/graph/view"
)

type AccountNodeViewRepository interface {
	ImportCsvFile(ctx context.Context, tenantId string, caseId string, importFile string, saveCallback ddd_neo4j.ImportSaveFileCallback, data []*view.AccountView, opts ...service.Options) error
	Create(ctx context.Context, view *view.AccountView, opts ...service.Options) error
	CreateMany(ctx context.Context, views []*view.AccountView, opts ...service.Options) error

	Update(ctx context.Context, view *view.AccountView, opts ...service.Options) error
	UpdateMany(ctx context.Context, views []*view.AccountView, opts ...service.Options) error

	Save(ctx context.Context, setData *ddd.SetData[*view.AccountView], opts ...service.Options) error

	Delete(ctx context.Context, view *view.AccountView, opts ...service.Options) error
	DeleteMany(ctx context.Context, tenantId string, views []*view.AccountView, opts ...service.Options) error
	DeleteById(ctx context.Context, tenantId string, id string, opts ...service.Options) error
	DeleteByIds(ctx context.Context, tenantId string, ids []string, opts ...service.Options) error
	DeleteByFilter(ctx context.Context, tenantId, filter string, opts ...service.Options) error
	DeleteAll(ctx context.Context, tenantId string, opts ...service.Options) error
	DeleteByGraphId(ctx context.Context, tenantId string, graphId string, opts ...service.Options) error
	DeleteByTenantId(ctx context.Context, tenantId string, opts ...service.Options) error

	FindById(ctx context.Context, tenantId string, id string, opts ...service.Options) (*view.AccountView, bool, error)
	FindByIds(ctx context.Context, tenantId string, ids []string, opts ...service.Options) ([]*view.AccountView, bool, error)
	FindAll(ctx context.Context, tenantId string, opts ...service.Options) ([]*view.AccountView, bool, error)
	FindPaging(ctx context.Context, query ddd_repository.FindPagingQuery, opts ...service.Options) (*ddd_repository.FindPagingResult[*view.AccountView], bool, error)
	FindByGraphId(ctx context.Context, tenantId string, graphId string, opts ...service.Options) ([]*view.AccountView, bool, error)
}
