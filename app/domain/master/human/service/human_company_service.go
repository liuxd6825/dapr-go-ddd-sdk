package service

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

// HumanCompanyService 人员-公司关联业务服务
type HumanCompanyService struct {
	dao *dao.HumanCompanyDao
}

func NewHumanCompanyService(dbKey string) *HumanCompanyService {
	return &HumanCompanyService{dao: dao.NewHumanCompanyDao(dbKey)}
}

func (s *HumanCompanyService) CreateMany(ctx context.Context, items []*model.HumanCompany, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *HumanCompanyService) UpdateMany(ctx context.Context, items []*model.HumanCompany, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *HumanCompanyService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *HumanCompanyService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.HumanCompany, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *HumanCompanyService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.HumanCompany, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *HumanCompanyService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanCompany] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *HumanCompanyService) FindByHumanId(ctx context.Context, qry store.FindPagingQuery, humanId string, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanCompany] {
	return s.dao.FindPagingByHumanId(ctx, qry, humanId, opts...)
}

func (s *HumanCompanyService) SubmitMany(ctx context.Context, insertData, updateData []*model.HumanCompany, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.HumanCompany `json:"insertData"`
		UpdateData []*model.HumanCompany `json:"updateData"`
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
