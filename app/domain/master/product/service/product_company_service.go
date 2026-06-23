package service

import (
	"context"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

// ProductCompanyService 产品-公司关联业务服务
type ProductCompanyService struct {
	dao *dao.ProductCompanyDao
}

var _productCompanyService *ProductCompanyService
var _productCompanyServiceOnce sync.Once

func NewProductCompanyService() *ProductCompanyService {
	_productCompanyServiceOnce.Do(func() {
		_productCompanyService = &ProductCompanyService{dao: dao.NewProductCompanyDao(config.DBKey)}
	})
	return _productCompanyService
}

func (s *ProductCompanyService) CreateMany(ctx context.Context, items []*model.ProductCompany, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *ProductCompanyService) UpdateMany(ctx context.Context, items []*model.ProductCompany, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *ProductCompanyService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *ProductCompanyService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.ProductCompany, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *ProductCompanyService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.ProductCompany, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *ProductCompanyService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.ProductCompany] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *ProductCompanyService) FindByProductId(ctx context.Context, qry store.FindPagingQuery, productId string, opts ...idao.CallOptions) store.FindPagingResult[*model.ProductCompany] {
	return s.dao.FindPagingByProductId(ctx, qry, productId, opts...)
}

func (s *ProductCompanyService) SubmitMany(ctx context.Context, insertData, updateData []*model.ProductCompany, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.ProductCompany `json:"insertData"`
		UpdateData []*model.ProductCompany `json:"updateData"`
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