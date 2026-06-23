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

// HumanExtService 人员-扩展信息业务服务
type HumanExtService struct {
	dao *dao.HumanExtDao
}

var _humanextService *HumanExtService
var _humanextServiceOnce sync.Once

func NewHumanExtService() *HumanExtService {
	_humanextServiceOnce.Do(func() {
		_humanextService = &HumanExtService{dao: dao.NewHumanExtDao(config.DBKey)}
	})
	return _humanextService
}

func (s *HumanExtService) CreateMany(ctx context.Context, items []*model.HumanExt, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *HumanExtService) UpdateMany(ctx context.Context, items []*model.HumanExt, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *HumanExtService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *HumanExtService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.HumanExt, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *HumanExtService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.HumanExt, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *HumanExtService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanExt] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *HumanExtService) FindByHumanId(ctx context.Context, qry store.FindPagingQuery, humanId string, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanExt] {
	return s.dao.FindPagingByHumanId(ctx, qry, humanId, opts...)
}

func (s *HumanExtService) SubmitMany(ctx context.Context, insertData, updateData []*model.HumanExt, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.HumanExt `json:"insertData"`
		UpdateData []*model.HumanExt `json:"updateData"`
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