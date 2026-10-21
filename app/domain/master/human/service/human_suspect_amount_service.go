package service

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

// HumanSuspectAmountService 嫌疑人金额业务服务
type HumanSuspectAmountService struct {
	dao *dao.HumanSuspectAmountDao
}

func NewHumanSuspectAmountService(dbKey string) *HumanSuspectAmountService {
	return &HumanSuspectAmountService{dao: dao.NewHumanSuspectAmountDao(dbKey)}
}

func (s *HumanSuspectAmountService) Create(ctx context.Context, data *model.HumanSuspectAmount, opts ...idao.CallOptions) error {
	return s.dao.Create(ctx, data, opts...).GetError()
}

func (s *HumanSuspectAmountService) Update(ctx context.Context, data *model.HumanSuspectAmount, opts ...idao.CallOptions) error {
	return s.dao.Update(ctx, data, opts...).GetError()
}

func (s *HumanSuspectAmountService) Submit(ctx context.Context, data *model.HumanSuspectAmount, opts ...idao.CallOptions) error {
	old, err := s.dao.FindById(ctx, data.Id, opts...)
	if err != nil {
		return err
	}
	if old != nil {
		return s.Update(ctx, data, opts...)
	}
	return s.Create(ctx, data, opts...)
}

func (s *HumanSuspectAmountService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *HumanSuspectAmountService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.HumanSuspectAmount, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *HumanSuspectAmountService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanSuspectAmount] {
	return s.dao.FindPaging(ctx, qry, opts...)
}
