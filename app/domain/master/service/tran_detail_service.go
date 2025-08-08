package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"sync"
)

type TranDetailService struct {
	dao *dao.TranDetailDao
}

var _tranDetailService *TranDetailService
var _tranDetailServiceOnce sync.Once

func NewTranDetailService() *TranDetailService {
	_tranDetailServiceOnce.Do(func() {
		_tranDetailService = &TranDetailService{
			dao: dao.NewTranDetailDao(config.DBKey),
		}
	})
	return _tranDetailService
}

func (r *TranDetailService) CreateMany(ctx context.Context, v []*model.TranDetail, opts ...idao.CallOptions) error {
	return r.dao.CreateMany(ctx, v, opts...).GetError()
}
