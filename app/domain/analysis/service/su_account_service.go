package service

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ddd/ddd_query"
)

type SuAccountService struct {
	dao *dao.SuAccountDao
}

func NewSuTaskAccountService() *SuAccountService {
	return &SuAccountService{
		dao: dao.NewSuAccountDao(config.DBKey),
	}
}

func (s *SuAccountService) Create(ctx context.Context, v *model.SuTaskAccount, opts ...idao.CallOptions) error {
	return s.dao.Create(ctx, v, opts...).GetError()
}

func (s *SuAccountService) CreateMany(ctx context.Context, v []*model.SuTaskAccount, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, v, opts...).GetError()
}

func (s *SuAccountService) Update(ctx context.Context, v *model.SuTaskAccount, opts ...idao.CallOptions) error {
	return s.dao.Update(ctx, v, opts...).GetError()
}

func (s *SuAccountService) UpdateMany(ctx context.Context, v []*model.SuTaskAccount, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, v, opts...).GetError()
}

func (s *SuAccountService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *SuAccountService) DeleteByIds(ctx context.Context, listId []string, opts ...idao.CallOptions) error {
	return s.dao.DeleteByIds(ctx, listId, opts...).GetError()
}

func (s *SuAccountService) FindById(ctx context.Context, qry *query.SuTaskAccountFindByIdQuery, opts ...idao.CallOptions) (*model.SuTaskAccount, error) {
	return s.dao.FindById(ctx, qry.Id, opts...)
}

func (s *SuAccountService) FindPaging(ctx context.Context, qry *ddd_query.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.SuTaskAccount] {
	return s.dao.FindPaging(ctx, qry, opts...)
}
