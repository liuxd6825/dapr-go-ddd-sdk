package service

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
)

type TranDetailService struct {
	dao *dao.TranDao
}

func NewTranDetailService(dbKey string) *TranDetailService {
	return &TranDetailService{
		dao: dao.NewTranDao(dbKey),
	}
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
