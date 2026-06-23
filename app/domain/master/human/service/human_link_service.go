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

// HumanLinkService 人员-联系方式业务服务
type HumanLinkService struct {
	dao *dao.HumanLinkDao
}

var _humanlinkService *HumanLinkService
var _humanlinkServiceOnce sync.Once

func NewHumanLinkService() *HumanLinkService {
	_humanlinkServiceOnce.Do(func() {
		_humanlinkService = &HumanLinkService{dao: dao.NewHumanLinkDao(config.DBKey)}
	})
	return _humanlinkService
}

func (s *HumanLinkService) CreateMany(ctx context.Context, items []*model.HumanLink, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *HumanLinkService) UpdateMany(ctx context.Context, items []*model.HumanLink, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *HumanLinkService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *HumanLinkService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.HumanLink, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *HumanLinkService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.HumanLink, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *HumanLinkService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanLink] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *HumanLinkService) FindByHumanId(ctx context.Context, qry store.FindPagingQuery, humanId string, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanLink] {
	return s.dao.FindPagingByHumanId(ctx, qry, humanId, opts...)
}

func (s *HumanLinkService) SubmitMany(ctx context.Context, insertData, updateData []*model.HumanLink, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.HumanLink `json:"insertData"`
		UpdateData []*model.HumanLink `json:"updateData"`
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