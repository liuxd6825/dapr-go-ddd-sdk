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

// CompanyAccountService 公司账号关联业务服务
type CompanyAccountService struct {
	dao *dao.CompanyAccountDao
}

var _companyAccountService *CompanyAccountService
var _companyAccountServiceOnce sync.Once

func NewCompanyAccountService() *CompanyAccountService {
	_companyAccountServiceOnce.Do(func() {
		_companyAccountService = &CompanyAccountService{
			dao: dao.NewCompanyAccountDao(config.DBKey),
		}
	})
	return _companyAccountService
}

func (s *CompanyAccountService) CreateMany(ctx context.Context, items []*model.CompanyAccount, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *CompanyAccountService) UpdateMany(ctx context.Context, items []*model.CompanyAccount, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *CompanyAccountService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *CompanyAccountService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.CompanyAccount, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *CompanyAccountService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.CompanyAccount, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *CompanyAccountService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.CompanyAccount] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *CompanyAccountService) FindByCompanyId(ctx context.Context, qry store.FindPagingQuery, companyId string, opts ...idao.CallOptions) store.FindPagingResult[*model.CompanyAccount] {
	return s.dao.FindPagingByCompanyId(ctx, qry, companyId, opts...)
}

func (s *CompanyAccountService) SubmitMany(ctx context.Context, insertData, updateData []*model.CompanyAccount, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.CompanyAccount `json:"insertData"`
		UpdateData []*model.CompanyAccount `json:"updateData"`
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