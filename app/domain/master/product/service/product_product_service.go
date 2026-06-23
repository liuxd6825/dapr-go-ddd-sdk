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

// ProductProductService 产品-产品关联业务服务
type ProductProductService struct {
	dao *dao.ProductProductDao
}

var _productProductService *ProductProductService
var _productProductServiceOnce sync.Once

func NewProductProductService() *ProductProductService {
	_productProductServiceOnce.Do(func() {
		_productProductService = &ProductProductService{dao: dao.NewProductProductDao(config.DBKey)}
	})
	return _productProductService
}

func (s *ProductProductService) CreateMany(ctx context.Context, items []*model.ProductProduct, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *ProductProductService) UpdateMany(ctx context.Context, items []*model.ProductProduct, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *ProductProductService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *ProductProductService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.ProductProduct, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *ProductProductService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.ProductProduct, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *ProductProductService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.ProductProduct] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *ProductProductService) FindByProductId(ctx context.Context, qry store.FindPagingQuery, productId string, opts ...idao.CallOptions) store.FindPagingResult[*model.ProductProduct] {
	return s.dao.FindPagingByProductId(ctx, qry, productId, opts...)
}

func (s *ProductProductService) SubmitMany(ctx context.Context, insertData, updateData []*model.ProductProduct, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.ProductProduct `json:"insertData"`
		UpdateData []*model.ProductProduct `json:"updateData"`
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