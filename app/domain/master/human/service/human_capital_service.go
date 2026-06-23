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

// HumanCapitalService 人员-资产关联业务服务
type HumanCapitalService struct {
	dao *dao.HumanCapitalDao
}

var _humancapitalService *HumanCapitalService
var _humancapitalServiceOnce sync.Once

func NewHumanCapitalService() *HumanCapitalService {
	_humancapitalServiceOnce.Do(func() {
		_humancapitalService = &HumanCapitalService{dao: dao.NewHumanCapitalDao(config.DBKey)}
	})
	return _humancapitalService
}

func (s *HumanCapitalService) CreateMany(ctx context.Context, items []*model.HumanCapital, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *HumanCapitalService) UpdateMany(ctx context.Context, items []*model.HumanCapital, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *HumanCapitalService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *HumanCapitalService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.HumanCapital, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *HumanCapitalService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.HumanCapital, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *HumanCapitalService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanCapital] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *HumanCapitalService) FindByHumanId(ctx context.Context, qry store.FindPagingQuery, humanId string, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanCapital] {
	return s.dao.FindPagingByHumanId(ctx, qry, humanId, opts...)
}

func (s *HumanCapitalService) SubmitMany(ctx context.Context, insertData, updateData []*model.HumanCapital, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.HumanCapital `json:"insertData"`
		UpdateData []*model.HumanCapital `json:"updateData"`
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