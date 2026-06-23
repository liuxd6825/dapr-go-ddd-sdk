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

// ProductContractService 产品-contract关联业务服务
type ProductContractService struct {
	dao *dao.ProductContractDao
}

var _productcontractService *ProductContractService
var _productcontractServiceOnce sync.Once

func NewProductContractService() *ProductContractService {
	_productcontractServiceOnce.Do(func() {
		_productcontractService = &ProductContractService{dao: dao.NewProductContractDao(config.DBKey)}
	})
	return _productcontractService
}

func (s *ProductContractService) CreateMany(ctx context.Context, items []*model.ProductContract, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *ProductContractService) UpdateMany(ctx context.Context, items []*model.ProductContract, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *ProductContractService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *ProductContractService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.ProductContract, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *ProductContractService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.ProductContract, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *ProductContractService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.ProductContract] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *ProductContractService) FindByProductId(ctx context.Context, qry store.FindPagingQuery, productId string, opts ...idao.CallOptions) store.FindPagingResult[*model.ProductContract] {
	return s.dao.FindPagingByProductId(ctx, qry, productId, opts...)
}

func (s *ProductContractService) SubmitMany(ctx context.Context, insertData, updateData []*model.ProductContract, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.ProductContract `json:"insertData"`
		UpdateData []*model.ProductContract `json:"updateData"`
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