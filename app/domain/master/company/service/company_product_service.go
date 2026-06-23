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

type CompanyProductService struct {
	dao *dao.CompanyProductDao
}

var _companyProductService *CompanyProductService
var _companyProductServiceOnce sync.Once

func NewCompanyProductService() *CompanyProductService {
	_companyProductServiceOnce.Do(func() {
		_companyProductService = &CompanyProductService{
			dao: dao.NewCompanyProductDao(config.DBKey),
		}
	})
	return _companyProductService
}

func (s *CompanyProductService) CreateMany(ctx context.Context, items []*model.CompanyProduct, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *CompanyProductService) UpdateMany(ctx context.Context, items []*model.CompanyProduct, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *CompanyProductService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *CompanyProductService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.CompanyProduct, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *CompanyProductService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.CompanyProduct, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *CompanyProductService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.CompanyProduct] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *CompanyProductService) FindByCompanyId(ctx context.Context, qry store.FindPagingQuery, companyId string, opts ...idao.CallOptions) store.FindPagingResult[*model.CompanyProduct] {
	return s.dao.FindPagingByCompanyId(ctx, qry, companyId, opts...)
}

func (s *CompanyProductService) SubmitMany(ctx context.Context, insertData, updateData []*model.CompanyProduct, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.CompanyProduct `json:"insertData"`
		UpdateData []*model.CompanyProduct `json:"updateData"`
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