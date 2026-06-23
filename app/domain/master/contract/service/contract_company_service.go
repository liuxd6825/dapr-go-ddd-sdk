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

// ContractCompanyService 合同-公司关联业务服务
type ContractCompanyService struct {
	dao *dao.ContractCompanyDao
}

var _contractCompanyService *ContractCompanyService
var _contractCompanyServiceOnce sync.Once

func NewContractCompanyService() *ContractCompanyService {
	_contractCompanyServiceOnce.Do(func() {
		_contractCompanyService = &ContractCompanyService{
			dao: dao.NewContractCompanyDao(config.DBKey),
		}
	})
	return _contractCompanyService
}

func (s *ContractCompanyService) CreateMany(ctx context.Context, items []*model.ContractCompany, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *ContractCompanyService) UpdateMany(ctx context.Context, items []*model.ContractCompany, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *ContractCompanyService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *ContractCompanyService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.ContractCompany, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *ContractCompanyService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.ContractCompany, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *ContractCompanyService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.ContractCompany] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *ContractCompanyService) FindByContractId(ctx context.Context, qry store.FindPagingQuery, contractId string, opts ...idao.CallOptions) store.FindPagingResult[*model.ContractCompany] {
	return s.dao.FindPagingByContractId(ctx, qry, contractId, opts...)
}

func (s *ContractCompanyService) SubmitMany(ctx context.Context, insertData, updateData []*model.ContractCompany, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.ContractCompany `json:"insertData"`
		UpdateData []*model.ContractCompany `json:"updateData"`
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