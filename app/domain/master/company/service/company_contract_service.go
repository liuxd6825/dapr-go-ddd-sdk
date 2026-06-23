package service

import (
	"context"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

type CompanyContractService struct {
	dao *dao.CompanyContractDao
}

var _companyContractService *CompanyContractService
var _companyContractServiceOnce sync.Once

func NewCompanyContractService() *CompanyContractService {
	_companyContractServiceOnce.Do(func() {
		_companyContractService = &CompanyContractService{
			dao: dao.NewCompanyContractDao(config.DBKey),
		}
	})
	return _companyContractService
}

func (s *CompanyContractService) CreateMany(ctx context.Context, items []*model.CompanyContract, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *CompanyContractService) UpdateMany(ctx context.Context, items []*model.CompanyContract, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *CompanyContractService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *CompanyContractService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.CompanyContract, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *CompanyContractService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.CompanyContract, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *CompanyContractService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.CompanyContract] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *CompanyContractService) FindByCompanyId(ctx context.Context, qry store.FindPagingQuery, companyId string, opts ...idao.CallOptions) store.FindPagingResult[*model.CompanyContract] {
	return s.dao.FindPagingByCompanyId(ctx, qry, companyId, opts...)
}

func (s *CompanyContractService) SubmitMany(ctx context.Context, insertData, updateData []*model.CompanyContract, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.CompanyContract `json:"insertData"`
		UpdateData []*model.CompanyContract `json:"updateData"`
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