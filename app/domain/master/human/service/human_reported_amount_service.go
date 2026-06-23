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

// HumanReportedAmountService 报案人金额业务服务
type HumanReportedAmountService struct {
	dao *dao.HumanReportedAmountDao
}

var _humanReportedAmountService *HumanReportedAmountService
var _humanReportedAmountServiceOnce sync.Once

func NewHumanReportedAmountService() *HumanReportedAmountService {
	_humanReportedAmountServiceOnce.Do(func() {
		_humanReportedAmountService = &HumanReportedAmountService{dao: dao.NewHumanReportedAmountDao(config.DBKey)}
	})
	return _humanReportedAmountService
}

func (s *HumanReportedAmountService) Create(ctx context.Context, data *model.HumanReportedAmount, opts ...idao.CallOptions) error {
	return s.dao.Create(ctx, data, opts...).GetError()
}

func (s *HumanReportedAmountService) Update(ctx context.Context, data *model.HumanReportedAmount, opts ...idao.CallOptions) error {
	return s.dao.Update(ctx, data, opts...).GetError()
}

func (s *HumanReportedAmountService) Submit(ctx context.Context, data *model.HumanReportedAmount, opts ...idao.CallOptions) error {
	old, err := s.dao.FindById(ctx, data.Id, opts...)
	if err != nil {
		return err
	}
	if old != nil {
		return s.Update(ctx, data, opts...)
	}
	return s.Create(ctx, data, opts...)
}

func (s *HumanReportedAmountService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *HumanReportedAmountService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.HumanReportedAmount, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *HumanReportedAmountService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanReportedAmount] {
	return s.dao.FindPaging(ctx, qry, opts...)
}