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

// HumanProductService 人员-产品关联业务服务
type HumanProductService struct {
	dao *dao.HumanProductDao
}

var _humanproductService *HumanProductService
var _humanproductServiceOnce sync.Once

func NewHumanProductService() *HumanProductService {
	_humanproductServiceOnce.Do(func() {
		_humanproductService = &HumanProductService{dao: dao.NewHumanProductDao(config.DBKey)}
	})
	return _humanproductService
}

func (s *HumanProductService) CreateMany(ctx context.Context, items []*model.HumanProduct, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *HumanProductService) UpdateMany(ctx context.Context, items []*model.HumanProduct, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *HumanProductService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *HumanProductService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.HumanProduct, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *HumanProductService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.HumanProduct, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *HumanProductService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanProduct] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *HumanProductService) FindByHumanId(ctx context.Context, qry store.FindPagingQuery, humanId string, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanProduct] {
	return s.dao.FindPagingByHumanId(ctx, qry, humanId, opts...)
}

func (s *HumanProductService) SubmitMany(ctx context.Context, insertData, updateData []*model.HumanProduct, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.HumanProduct `json:"insertData"`
		UpdateData []*model.HumanProduct `json:"updateData"`
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