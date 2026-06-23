package service

import (
	"context"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

type CompanyHumanService struct {
	dao *dao.CompanyHumanDao
}

var _companyHumanService *CompanyHumanService
var _companyHumanServiceOnce sync.Once

func NewCompanyHumanService() *CompanyHumanService {
	_companyHumanServiceOnce.Do(func() {
		_companyHumanService = &CompanyHumanService{
			dao: dao.NewCompanyHumanDao(config.DBKey),
		}
	})
	return _companyHumanService
}

func (s *CompanyHumanService) CreateMany(ctx context.Context, items []*model.CompanyHuman, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *CompanyHumanService) UpdateMany(ctx context.Context, items []*model.CompanyHuman, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *CompanyHumanService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *CompanyHumanService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.CompanyHuman, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *CompanyHumanService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.CompanyHuman, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *CompanyHumanService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.CompanyHuman] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *CompanyHumanService) FindByCompanyId(ctx context.Context, qry store.FindPagingQuery, companyId string, opts ...idao.CallOptions) store.FindPagingResult[*model.CompanyHuman] {
	return s.dao.FindPagingByCompanyId(ctx, qry, companyId, opts...)
}

func (s *CompanyHumanService) SubmitMany(ctx context.Context, insertData, updateData []*model.CompanyHuman, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.CompanyHuman `json:"insertData"`
		UpdateData []*model.CompanyHuman `json:"updateData"`
	}{}
	if len(insertData) > 0 {
		if err := s.dao.CreateMany(ctx, insertData, opts...).GetError(); err != nil {
			return nil, err
		}
		ids := make([]string, 0, len(insertData))
		for _, item := range insertData {
			ids = append(ids, item.Id)
		}
		items, err := s.dao.FindByIds(ctx, ids, opts...)
		if err != nil {
			return nil, err
		}
		res.InsertData = items
	}
	if len(updateData) > 0 {
		if err := s.dao.UpdateMany(ctx, updateData, opts...).GetError(); err != nil {
			return nil, err
		}
		ids := make([]string, 0, len(updateData))
		for _, item := range updateData {
			ids = append(ids, item.Id)
		}
		items, err := s.dao.FindByIds(ctx, ids, opts...)
		if err != nil {
			return nil, err
		}
		res.UpdateData = items
	}
	return res, nil
}