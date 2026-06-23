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

// HumanHumanService 人员-人员关联业务服务
type HumanHumanService struct {
	dao *dao.HumanHumanDao
}

var _humanhumanService *HumanHumanService
var _humanhumanServiceOnce sync.Once

func NewHumanHumanService() *HumanHumanService {
	_humanhumanServiceOnce.Do(func() {
		_humanhumanService = &HumanHumanService{dao: dao.NewHumanHumanDao(config.DBKey)}
	})
	return _humanhumanService
}

func (s *HumanHumanService) CreateMany(ctx context.Context, items []*model.HumanHuman, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *HumanHumanService) UpdateMany(ctx context.Context, items []*model.HumanHuman, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *HumanHumanService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *HumanHumanService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.HumanHuman, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *HumanHumanService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.HumanHuman, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *HumanHumanService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanHuman] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *HumanHumanService) FindByHumanId(ctx context.Context, qry store.FindPagingQuery, humanId string, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanHuman] {
	return s.dao.FindPagingByHumanId(ctx, qry, humanId, opts...)
}

func (s *HumanHumanService) SubmitMany(ctx context.Context, insertData, updateData []*model.HumanHuman, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.HumanHuman `json:"insertData"`
		UpdateData []*model.HumanHuman `json:"updateData"`
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