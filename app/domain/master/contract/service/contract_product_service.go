package service

import (
	"context"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

// ContractProductService 合同-产品关联业务服务
type ContractProductService struct {
	dao *dao.ContractProductDao
}

var _contractProductService *ContractProductService
var _contractProductServiceOnce sync.Once

func NewContractProductService() *ContractProductService {
	_contractProductServiceOnce.Do(func() {
		_contractProductService = &ContractProductService{
			dao: dao.NewContractProductDao(config.DBKey),
		}
	})
	return _contractProductService
}

func (s *ContractProductService) CreateMany(ctx context.Context, items []*model.ContractProduct, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *ContractProductService) UpdateMany(ctx context.Context, items []*model.ContractProduct, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *ContractProductService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *ContractProductService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.ContractProduct, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *ContractProductService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.ContractProduct, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *ContractProductService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.ContractProduct] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *ContractProductService) FindByContractId(ctx context.Context, qry store.FindPagingQuery, contractId string, opts ...idao.CallOptions) store.FindPagingResult[*model.ContractProduct] {
	return s.dao.FindPagingByContractId(ctx, qry, contractId, opts...)
}

func (s *ContractProductService) SubmitMany(ctx context.Context, insertData, updateData []*model.ContractProduct, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.ContractProduct `json:"insertData"`
		UpdateData []*model.ContractProduct `json:"updateData"`
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