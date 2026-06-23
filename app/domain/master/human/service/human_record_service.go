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

// HumanRecordService 人员-记录业务服务
type HumanRecordService struct {
	dao *dao.HumanRecordDao
}

var _humanrecordService *HumanRecordService
var _humanrecordServiceOnce sync.Once

func NewHumanRecordService() *HumanRecordService {
	_humanrecordServiceOnce.Do(func() {
		_humanrecordService = &HumanRecordService{dao: dao.NewHumanRecordDao(config.DBKey)}
	})
	return _humanrecordService
}

func (s *HumanRecordService) CreateMany(ctx context.Context, items []*model.HumanRecord, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *HumanRecordService) UpdateMany(ctx context.Context, items []*model.HumanRecord, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *HumanRecordService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *HumanRecordService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.HumanRecord, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *HumanRecordService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.HumanRecord, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *HumanRecordService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanRecord] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *HumanRecordService) FindByHumanId(ctx context.Context, qry store.FindPagingQuery, humanId string, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanRecord] {
	return s.dao.FindPagingByHumanId(ctx, qry, humanId, opts...)
}

func (s *HumanRecordService) SubmitMany(ctx context.Context, insertData, updateData []*model.HumanRecord, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.HumanRecord `json:"insertData"`
		UpdateData []*model.HumanRecord `json:"updateData"`
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