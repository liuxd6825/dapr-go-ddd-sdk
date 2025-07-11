package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/singleutils"
)

type SheetDomainService struct {
	repos *dao.SheetDao
}

func NewSheetDomainService() *SheetDomainService {
	return singleutils.CreateObj[*SheetDomainService](func() *SheetDomainService {
		return &SheetDomainService{
			repos: dao.NewSheetDao(config.DBKey),
		}
	})
}

func (f *SheetDomainService) Create(ctx context.Context, m *model.Sheet, opts ...idao.CallOptions) error {
	return f.repos.Create(ctx, m, opts...).GetError()
}

func (f *SheetDomainService) CreateMany(ctx context.Context, m []*model.Sheet, opts ...idao.CallOptions) error {
	return f.repos.CreateMany(ctx, m, opts...).GetError()
}

func (f *SheetDomainService) Update(ctx context.Context, m *model.Sheet, opts ...idao.CallOptions) error {
	return f.repos.Update(ctx, m, opts...).GetError()
}

func (f *SheetDomainService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return f.repos.DeleteById(ctx, id, opts...).GetError()
}

func (f *SheetDomainService) FindById(ctx context.Context, tenantId string, id string, opts ...idao.CallOptions) (*model.Sheet, error) {
	return f.repos.FindById(ctx, id, opts...)
}

func (f *SheetDomainService) FindPaging(ctx context.Context, qry idao.FindPagingQuery, opts ...idao.CallOptions) idao.FindPagingResult[*model.Sheet] {
	return f.repos.FindPaging(ctx, qry, opts...)
}

func (f *SheetDomainService) FindByName(ctx context.Context, qry *query.FindSheetByNameQuery) ([]*model.Sheet, error) {
	res := f.repos.FindByName(ctx, qry)
	return res.GetData(), res.GetError()
}
