package service

import (
	"context"
	"gitee.com/liuxu6825/duxm-master-service/pkg/master/cmd-service/infrastructure/base/v2/domain/service"
	"gitee.com/liuxu6825/duxm-master-service/pkg/master/import-service/domain/excelview/model"
	"gitee.com/liuxu6825/duxm-master-service/pkg/master/import-service/domain/excelview/query"
	"gitee.com/liuxu6825/duxm-master-service/pkg/master/import-service/domain/excelview/repository"
	"gitee.com/liuxu6825/duxm-master-service/pkg/master/import-service/domain/excelview/repository/mongo"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/singleutils"
)

type SheetDomainService interface {
	Create(ctx context.Context, m *model.Sheet, opts ...service.Options) error
	Update(ctx context.Context, m *model.Sheet, opts ...service.Options) error
	DeleteById(ctx context.Context, tenantId string, id string, opts ...service.Options) error
	FindById(ctx context.Context, tenantId string, id string, opts ...service.Options) (*model.Sheet, bool, error)
	FindPaging(ctx context.Context, qry query.FindPagingQuery, opts ...service.Options) (*ddd_repository.FindPagingResult[*model.Sheet], bool, error)
	FindByName(ctx context.Context, qry *query.FindSheetByNameQuery) ([]*model.Sheet, bool, error)
}

type sheetDomainService struct {
	repos repository.SheetRepository
}

func NewSheetDomainService() SheetDomainService {
	return singleutils.CreateObj[*sheetDomainService](func() *sheetDomainService {
		return &sheetDomainService{
			repos: mongo.NewSheetRepository(),
		}
	})
}

func (f *sheetDomainService) Create(ctx context.Context, m *model.Sheet, opts ...service.Options) error {
	return f.repos.Create(ctx, m, opts...)
}

func (f *sheetDomainService) CreateMany(ctx context.Context, m []*model.Sheet, opts ...service.Options) error {
	return f.repos.CreateMany(ctx, m, opts...)
}

func (f *sheetDomainService) Update(ctx context.Context, m *model.Sheet, opts ...service.Options) error {
	return f.repos.Update(ctx, m, opts...)
}

func (f *sheetDomainService) DeleteById(ctx context.Context, tenantId string, id string, opts ...service.Options) error {
	return f.repos.DeleteById(ctx, tenantId, id, opts...)
}

func (f *sheetDomainService) FindById(ctx context.Context, tenantId string, id string, opts ...service.Options) (*model.Sheet, bool, error) {
	return f.repos.FindById(ctx, tenantId, id, opts...)
}

func (f *sheetDomainService) FindPaging(ctx context.Context, qry query.FindPagingQuery, opts ...service.Options) (*ddd_repository.FindPagingResult[*model.Sheet], bool, error) {
	return f.repos.FindPaging(ctx, qry, opts...)
}

func (f *sheetDomainService) FindByName(ctx context.Context, qry *query.FindSheetByNameQuery) ([]*model.Sheet, bool, error) {
	return f.repos.FindByName(ctx, qry)
}
