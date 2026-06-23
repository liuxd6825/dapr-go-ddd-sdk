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

// ContractRecordService 合同-记录关联业务服务
type ContractRecordService struct {
	dao *dao.ContractRecordDao
}

var _contractRecordService *ContractRecordService
var _contractRecordServiceOnce sync.Once

func NewContractRecordService() *ContractRecordService {
	_contractRecordServiceOnce.Do(func() {
		_contractRecordService = &ContractRecordService{
			dao: dao.NewContractRecordDao(config.DBKey),
		}
	})
	return _contractRecordService
}

func (s *ContractRecordService) CreateMany(ctx context.Context, items []*model.ContractRecord, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *ContractRecordService) UpdateMany(ctx context.Context, items []*model.ContractRecord, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *ContractRecordService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *ContractRecordService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.ContractRecord, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *ContractRecordService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.ContractRecord, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *ContractRecordService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.ContractRecord] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *ContractRecordService) FindByContractId(ctx context.Context, qry store.FindPagingQuery, contractId string, opts ...idao.CallOptions) store.FindPagingResult[*model.ContractRecord] {
	return s.dao.FindPagingByContractId(ctx, qry, contractId, opts...)
}

func (s *ContractRecordService) SubmitMany(ctx context.Context, insertData, updateData []*model.ContractRecord, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.ContractRecord `json:"insertData"`
		UpdateData []*model.ContractRecord `json:"updateData"`
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