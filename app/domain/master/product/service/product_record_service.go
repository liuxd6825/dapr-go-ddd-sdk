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

// ProductRecordService 产品-record关联业务服务
type ProductRecordService struct {
	dao *dao.ProductRecordDao
}

var _productRecordService *ProductRecordService
var _productRecordServiceOnce sync.Once

func NewProductRecordService() *ProductRecordService {
	_productRecordServiceOnce.Do(func() {
		_productRecordService = &ProductRecordService{dao: dao.NewProductRecordDao(config.DBKey)}
	})
	return _productRecordService
}

func (s *ProductRecordService) CreateMany(ctx context.Context, items []*model.ProductRecord, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *ProductRecordService) UpdateMany(ctx context.Context, items []*model.ProductRecord, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *ProductRecordService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *ProductRecordService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.ProductRecord, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *ProductRecordService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.ProductRecord, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *ProductRecordService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.ProductRecord] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *ProductRecordService) FindByProductId(ctx context.Context, qry store.FindPagingQuery, productId string, opts ...idao.CallOptions) store.FindPagingResult[*model.ProductRecord] {
	return s.dao.FindPagingByProductId(ctx, qry, productId, opts...)
}

func (s *ProductRecordService) SubmitMany(ctx context.Context, insertData, updateData []*model.ProductRecord, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.ProductRecord `json:"insertData"`
		UpdateData []*model.ProductRecord `json:"updateData"`
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