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

// HumanContractService 人员-合同关联业务服务
type HumanContractService struct {
	dao *dao.HumanContractDao
}

var _humancontractService *HumanContractService
var _humancontractServiceOnce sync.Once

func NewHumanContractService() *HumanContractService {
	_humancontractServiceOnce.Do(func() {
		_humancontractService = &HumanContractService{dao: dao.NewHumanContractDao(config.DBKey)}
	})
	return _humancontractService
}

func (s *HumanContractService) CreateMany(ctx context.Context, items []*model.HumanContract, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *HumanContractService) UpdateMany(ctx context.Context, items []*model.HumanContract, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *HumanContractService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *HumanContractService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.HumanContract, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *HumanContractService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.HumanContract, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *HumanContractService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanContract] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *HumanContractService) FindByHumanId(ctx context.Context, qry store.FindPagingQuery, humanId string, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanContract] {
	return s.dao.FindPagingByHumanId(ctx, qry, humanId, opts...)
}

func (s *HumanContractService) SubmitMany(ctx context.Context, insertData, updateData []*model.HumanContract, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.HumanContract `json:"insertData"`
		UpdateData []*model.HumanContract `json:"updateData"`
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