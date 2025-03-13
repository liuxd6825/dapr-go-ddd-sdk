package impl

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
)

/*
func (d *DaoBase) CountByMap(ctx context.Context, filterData any, opts ...*db.CallOptions) int64 {
	opt := db.NewCallOptions(opts...)

	tenantId := d.GetTenantId(ctx)
	count, err := d.dao.CountByMap(ctx, tenantId, filterData, opt)
	if err != nil {
		panic(err)
	}
	return count
}*/

func (d *DaoBase) CountByRSQL(ctx context.Context, rSQL string, opts ...*db.CallOptions) int64 {
	opt := db.NewCallOptions(opts...)
	tenantId := d.GetTenantId(ctx)
	count, err := d.dao.CountByRSQL(ctx, tenantId, rSQL, opt)
	if err != nil {
		panic(err)
	}
	return count
}

func (d *DaoBase) SumEntity(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, opts ...*db.CallOptions) []map[string]any {
	data, _, err := d.dao.SumEntity(ctx, qry, db.NewRepositoryOptions(opts)...)
	if err != nil {
		panic(err)
	}
	return data
}

/*
func (d *DaoBase) SumMap(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, opts ...*db.CallOptions) []map[string]any {
	data, _, err := d.dao.SumMap(ctx, qry, db.NewRepositoryOptions(opts)...)
	if err != nil {
		panic(err)
	}
	return data
}
*/

func (d *DaoBase) SumByQuery(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, data any, opts ...*db.CallOptions) any {
	data, _, err := d.dao.SumByQuery(ctx, qry, data, db.NewRepositoryOptions(opts)...)
	if err != nil {
		panic(err)
	}
	return data
}

func (d *DaoBase) SumByRSQL(ctx context.Context, rSQL string, valueCols []*ddd_repository.ValueCol, opts ...*db.CallOptions) map[string]any {
	opt := db.NewCallOptions(opts...)
	tenantId := d.GetTenantId(ctx)
	return d.dao.SumByRSQL(ctx, tenantId, rSQL, valueCols, opt)
}
