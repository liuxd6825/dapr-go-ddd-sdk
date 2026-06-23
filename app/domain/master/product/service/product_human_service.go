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

// ProductHumanService 产品-human关联业务服务
type ProductHumanService struct {
	dao *dao.ProductHumanDao
}

var _producthumanService *ProductHumanService
var _producthumanServiceOnce sync.Once

func NewProductHumanService() *ProductHumanService {
	_producthumanServiceOnce.Do(func() {
		_producthumanService = &ProductHumanService{dao: dao.NewProductHumanDao(config.DBKey)}
	})
	return _producthumanService
}

func (s *ProductHumanService) CreateMany(ctx context.Context, items []*model.ProductHuman, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *ProductHumanService) UpdateMany(ctx context.Context, items []*model.ProductHuman, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *ProductHumanService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *ProductHumanService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.ProductHuman, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *ProductHumanService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.ProductHuman, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *ProductHumanService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.ProductHuman] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *ProductHumanService) FindByProductId(ctx context.Context, qry store.FindPagingQuery, productId string, opts ...idao.CallOptions) store.FindPagingResult[*model.ProductHuman] {
	return s.dao.FindPagingByProductId(ctx, qry, productId, opts...)
}

func (s *ProductHumanService) SubmitMany(ctx context.Context, insertData, updateData []*model.ProductHuman, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.ProductHuman `json:"insertData"`
		UpdateData []*model.ProductHuman `json:"updateData"`
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