package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
	"sync"
)

type TranDetailService struct {
	dao *dao.TranDao
}

var _tranDetailService *TranDetailService
var _tranDetailServiceOnce sync.Once

func NewTranDetailService() *TranDetailService {
	_tranDetailServiceOnce.Do(func() {
		_tranDetailService = &TranDetailService{
			dao: dao.NewTranDao(config.DBKey),
		}
	})
	return _tranDetailService
}

func (r *TranDetailService) CreateMany(ctx context.Context, v []*model.Tran, opts ...idao.CallOptions) error {
	return r.dao.CreateMany(ctx, v, opts...).GetError()
}

func (r *TranDetailService) FindThresholdQuery(ctx context.Context, qry *query.TranDetailFindThresholdQuery, opts ...idao.CallOptions) ([]*model.Tran, error) {
	build := rsql.NewBuilder().And(
		rsql.Or(
			rsql.Eq("name", qry.Name),
			rsql.Eq("oppName", qry.Name),
		),
		rsql.Gte("date", qry.StartDate),
		rsql.Lte("date", qry.EndDate),
		rsql.Gte("amount", qry.Amount),
	)
	qry.SetMustFilter(build.Build())

	result := r.dao.FindPaging(ctx, qry, opts...)
	return result.GetData(), result.GetError()
}
