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

// ContractHumanService 合同-人员关联业务服务
type ContractHumanService struct {
	dao *dao.ContractHumanDao
}

var _contractHumanService *ContractHumanService
var _contractHumanServiceOnce sync.Once

func NewContractHumanService() *ContractHumanService {
	_contractHumanServiceOnce.Do(func() {
		_contractHumanService = &ContractHumanService{
			dao: dao.NewContractHumanDao(config.DBKey),
		}
	})
	return _contractHumanService
}

func (s *ContractHumanService) CreateMany(ctx context.Context, items []*model.ContractHuman, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *ContractHumanService) UpdateMany(ctx context.Context, items []*model.ContractHuman, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *ContractHumanService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *ContractHumanService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.ContractHuman, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *ContractHumanService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.ContractHuman, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *ContractHumanService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.ContractHuman] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *ContractHumanService) FindByContractId(ctx context.Context, qry store.FindPagingQuery, contractId string, opts ...idao.CallOptions) store.FindPagingResult[*model.ContractHuman] {
	return s.dao.FindPagingByContractId(ctx, qry, contractId, opts...)
}

func (s *ContractHumanService) SubmitMany(ctx context.Context, insertData, updateData []*model.ContractHuman, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.ContractHuman `json:"insertData"`
		UpdateData []*model.ContractHuman `json:"updateData"`
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