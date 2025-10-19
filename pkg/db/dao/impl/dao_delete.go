package impl

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

func (d *DaoBase[T]) Delete(ctx context.Context, entity T, opts ...idao.CallOptions) *idao.Result {
	id := d.store.GetId(entity)
	if id == "" {
		panic("GetId() return is nil")
	}
	return d.DeleteById(ctx, id, opts...)
}

func (d *DaoBase[T]) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) *idao.Result {
	tenantId := d.GetTenantId(ctx)

	res := d.store.DeleteById(ctx, tenantId, id, idao.NewCallOptions(opts...))
	if res.Error != nil {
		panic(res.Error)
	}

	return idao.NewResult(res)
}

func (d *DaoBase[T]) deleteById(ctx context.Context, id string, opts ...idao.CallOptions) int64 {
	tenantId := d.GetTenantId(ctx)

	res := d.store.DeleteById(ctx, tenantId, id, idao.NewCallOptions(opts...))
	if res.Error != nil {
		panic(res.Error)
	}

	return res.RowsAffected
}

func (d *DaoBase[T]) DeleteByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) *idao.Result {
	tenantId := d.GetTenantId(ctx)
	res := d.store.DeleteByIds(ctx, tenantId, ids, idao.NewCallOptions(opts...))
	if res.Error != nil {
		panic(res.Error)
	}
	return idao.NewResult(res)
}

func (d *DaoBase[T]) DeleteMany(ctx context.Context, entities []T, opts ...idao.CallOptions) *idao.Result {
	tenantId := d.GetTenantId(ctx)
	res := d.store.DeleteMany(ctx, tenantId, entities)
	if res.Error != nil {
		panic(res.Error)
	}
	return idao.NewResult(res)
}

func (d *DaoBase[T]) DeleteAll(ctx context.Context, opts ...idao.CallOptions) *idao.Result {
	tenantId := d.GetTenantId(ctx)
	res := d.store.DeleteAll(ctx, tenantId, idao.NewCallOptions(opts...))
	if res.Error != nil {
		panic(res.Error)
	}
	return idao.NewResult(res)
}

func (d *DaoBase[T]) DeleteByRSQL(ctx context.Context, filterRSQL string, opts ...idao.CallOptions) *idao.Result {
	tenantId := d.GetTenantId(ctx)
	res := d.store.DeleteByRSQL(ctx, tenantId, filterRSQL, idao.NewCallOptions(opts...))
	return idao.NewResult(res)
}

/*
func (d *DaoBase) DeleteByMap(ctx context.Context, filterMap map[string]interface{}, opts ...*db.CallOptions) *db.Result {
	tenantId := d.GetTenantId(ctx)
	res := d.dao.DeleteByMap(ctx, tenantId, filterMap, db.NewRepositoryOptions(opts)...)
	return db.NewResult(res)
}
*/
