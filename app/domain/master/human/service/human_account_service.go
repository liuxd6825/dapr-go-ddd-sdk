package service

import (
	"context"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

// HumanAccountService 人员-账户关联业务服务
type HumanAccountService struct {
	dao *dao.HumanAccountDao
}

var _humanAccountService *HumanAccountService
var _humanAccountServiceOnce sync.Once

func NewHumanAccountService() *HumanAccountService {
	_humanAccountServiceOnce.Do(func() {
		_humanAccountService = &HumanAccountService{dao: dao.NewHumanAccountDao(config.DBKey)}
	})
	return _humanAccountService
}

func (s *HumanAccountService) CreateMany(ctx context.Context, items []*model.HumanAccount, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *HumanAccountService) UpdateMany(ctx context.Context, items []*model.HumanAccount, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *HumanAccountService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *HumanAccountService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.HumanAccount, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *HumanAccountService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.HumanAccount, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *HumanAccountService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanAccount] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *HumanAccountService) FindByHumanId(ctx context.Context, qry store.FindPagingQuery, humanId string, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanAccount] {
	return s.dao.FindPagingByHumanId(ctx, qry, humanId, opts...)
}

func (s *HumanAccountService) SubmitMany(ctx context.Context, insertData, updateData []*model.HumanAccount, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.HumanAccount `json:"insertData"`
		UpdateData []*model.HumanAccount `json:"updateData"`
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