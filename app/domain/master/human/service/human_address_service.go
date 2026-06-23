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

// HumanAddressService 人员-地址关联业务服务
type HumanAddressService struct {
	dao *dao.HumanAddressDao
}

var _humanaddressService *HumanAddressService
var _humanaddressServiceOnce sync.Once

func NewHumanAddressService() *HumanAddressService {
	_humanaddressServiceOnce.Do(func() {
		_humanaddressService = &HumanAddressService{dao: dao.NewHumanAddressDao(config.DBKey)}
	})
	return _humanaddressService
}

func (s *HumanAddressService) CreateMany(ctx context.Context, items []*model.HumanAddress, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *HumanAddressService) UpdateMany(ctx context.Context, items []*model.HumanAddress, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *HumanAddressService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *HumanAddressService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.HumanAddress, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *HumanAddressService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.HumanAddress, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *HumanAddressService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanAddress] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *HumanAddressService) FindByHumanId(ctx context.Context, qry store.FindPagingQuery, humanId string, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanAddress] {
	return s.dao.FindPagingByHumanId(ctx, qry, humanId, opts...)
}

func (s *HumanAddressService) SubmitMany(ctx context.Context, insertData, updateData []*model.HumanAddress, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.HumanAddress `json:"insertData"`
		UpdateData []*model.HumanAddress `json:"updateData"`
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