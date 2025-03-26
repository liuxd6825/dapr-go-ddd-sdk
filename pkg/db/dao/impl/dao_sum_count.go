package impl

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
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
	count, err := d.store.CountByRSQL(ctx, tenantId, rSQL, opt)
	if err != nil {
		panic(err)
	}
	return count
}

func (d *DaoBase[T]) SumEntity(ctx context.Context, qry *store.FindPagingQueryRequest, opts ...*idao.CallOptions) []T {
	data, _, err := d.store.SumEntity(ctx, qry, idao.NewRepositoryOptions(opts)...)
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

func (d *DaoBase[T]) SumByQuery(ctx context.Context, qry *store.FindPagingQueryRequest, opts ...*idao.CallOptions) map[string]any {
	data := map[string]any{}
	_, _, err := d.store.SumByQuery(ctx, qry, data, idao.NewRepositoryOptions(opts)...)
	if err != nil {
		panic(err)
	}
	return data
}

func (d *DaoBase[T]) SumByRSQL(ctx context.Context, rSQL string, valueCols []*store.ValueCol, opts ...*idao.CallOptions) map[string]any {
	opt := idao.NewCallOptions(opts...)
	tenantId := d.GetTenantId(ctx)
	return d.store.SumByRSQL(ctx, tenantId, rSQL, valueCols, opt)
}
