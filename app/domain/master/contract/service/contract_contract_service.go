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

// ContractContractService 合同-合同关联业务服务
type ContractContractService struct {
	dao *dao.ContractContractDao
}

var _contractContractService *ContractContractService
var _contractContractServiceOnce sync.Once

func NewContractContractService() *ContractContractService {
	_contractContractServiceOnce.Do(func() {
		_contractContractService = &ContractContractService{
			dao: dao.NewContractContractDao(config.DBKey),
		}
	})
	return _contractContractService
}

func (s *ContractContractService) CreateMany(ctx context.Context, items []*model.ContractContract, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *ContractContractService) UpdateMany(ctx context.Context, items []*model.ContractContract, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *ContractContractService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *ContractContractService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.ContractContract, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *ContractContractService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.ContractContract, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *ContractContractService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.ContractContract] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *ContractContractService) FindByContractId(ctx context.Context, qry store.FindPagingQuery, contractId string, opts ...idao.CallOptions) store.FindPagingResult[*model.ContractContract] {
	return s.dao.FindPagingByContractId(ctx, qry, contractId, opts...)
}

func (s *ContractContractService) SubmitMany(ctx context.Context, insertData, updateData []*model.ContractContract, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.ContractContract `json:"insertData"`
		UpdateData []*model.ContractContract `json:"updateData"`
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