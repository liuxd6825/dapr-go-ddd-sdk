package impl

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
)

/*
func (d *DaoBase) CountByMap(ctx context.Context, filterData any, opts ...*idao.CallOptions) int64 {
	opt := idao.NewCallOptions(opts...)

	tenantId := d.GetTenantId(ctx)
	count, err := d.dao.CountByMap(ctx, tenantId, filterData, opt)
	if err != nil {
		panic(err)
	}
	return count
}*/

func (d *DaoBase[T]) CountByRSQL(ctx context.Context, rSQL string, opts ...*idao.CallOptions) int64 {
	opt := idao.NewCallOptions(opts...)
	tenantId := d.GetTenantId(ctx)
	count, err := d.dao.CountByRSQL(ctx, tenantId, rSQL, opt)
	if err != nil {
		panic(err)
	}
	return count
}

func (d *DaoBase[T]) SumEntity(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, opts ...*idao.CallOptions) []T {
	data, _, err := d.dao.SumEntity(ctx, qry, idao.NewRepositoryOptions(opts)...)
	if err != nil {
		panic(err)
	}
	return data
}

/*
func (d *DaoBase) SumMap(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, opts ...*idao.CallOptions) []map[string]any {
	data, _, err := d.dao.SumMap(ctx, qry, idao.NewRepositoryOptions(opts)...)
	if err != nil {
		panic(err)
	}
	return data
}
*/

func (d *DaoBase[T]) SumByQuery(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, opts ...*idao.CallOptions) map[string]any {
	data := map[string]any{}
	_, _, err := d.dao.SumByQuery(ctx, qry, data, idao.NewRepositoryOptions(opts)...)
	if err != nil {
		panic(err)
	}
	return data
}

func (d *DaoBase[T]) SumByRSQL(ctx context.Context, rSQL string, valueCols []*ddd_repository.ValueCol, opts ...*idao.CallOptions) map[string]any {
	opt := idao.NewCallOptions(opts...)
	tenantId := d.GetTenantId(ctx)
	return d.dao.SumByRSQL(ctx, tenantId, rSQL, valueCols, opt)
}
