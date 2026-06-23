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

// HumanCredentialService 人员-证件业务服务
type HumanCredentialService struct {
	dao *dao.HumanCredentialDao
}

var _humancredentialService *HumanCredentialService
var _humancredentialServiceOnce sync.Once

func NewHumanCredentialService() *HumanCredentialService {
	_humancredentialServiceOnce.Do(func() {
		_humancredentialService = &HumanCredentialService{dao: dao.NewHumanCredentialDao(config.DBKey)}
	})
	return _humancredentialService
}

func (s *HumanCredentialService) CreateMany(ctx context.Context, items []*model.HumanCredential, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, items, opts...).GetError()
}

func (s *HumanCredentialService) UpdateMany(ctx context.Context, items []*model.HumanCredential, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, items, opts...).GetError()
}

func (s *HumanCredentialService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *HumanCredentialService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.HumanCredential, error) {
	return s.dao.FindById(ctx, id, opts...)
}

func (s *HumanCredentialService) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]*model.HumanCredential, error) {
	return s.dao.FindByIds(ctx, ids, opts...)
}

func (s *HumanCredentialService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanCredential] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *HumanCredentialService) FindByHumanId(ctx context.Context, qry store.FindPagingQuery, humanId string, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanCredential] {
	return s.dao.FindPagingByHumanId(ctx, qry, humanId, opts...)
}

func (s *HumanCredentialService) SubmitMany(ctx context.Context, insertData, updateData []*model.HumanCredential, opts ...idao.CallOptions) (any, error) {
	res := struct {
		InsertData []*model.HumanCredential `json:"insertData"`
		UpdateData []*model.HumanCredential `json:"updateData"`
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