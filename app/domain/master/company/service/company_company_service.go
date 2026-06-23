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

type CompanyCompanyService struct {
	dao *dao.CompanyCompanyDao
}

var _companyCompanyService *CompanyCompanyService
var _companyCompanyServiceOnce sync.Once

func NewCompanyCompanyService() *CompanyCompanyService {
	_companyCompanyServiceOnce.Do(func() {
		_companyCompanyService = &CompanyCompanyService{
			dao: dao.NewCompanyCompanyDao(config.DBKey),
		}
	})
	return _companyCompanyService
}

func (s *CompanyCompanyService) CreateMany(ctx context.Context, items []*model.CompanyCompany, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *CompanyCompanyService) UpdateMany(ctx context.Context, items []*model.CompanyCompany, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *CompanyCompanyService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *CompanyCompanyService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.CompanyCompany, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *CompanyCompanyService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.CompanyCompany, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *CompanyCompanyService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.CompanyCompany] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *CompanyCompanyService) FindByCompanyId(ctx context.Context, qry store.FindPagingQuery, companyId string, opts ...idao.CallOptions) store.FindPagingResult[*model.CompanyCompany] {
	return s.dao.FindPagingByCompanyId(ctx, qry, companyId, opts...)
}

func (s *CompanyCompanyService) SubmitMany(ctx context.Context, insertData, updateData []*model.CompanyCompany, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.CompanyCompany `json:"insertData"`
		UpdateData []*model.CompanyCompany `json:"updateData"`
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