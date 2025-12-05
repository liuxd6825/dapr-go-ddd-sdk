package impl

import (
	"context"
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

func (d *DaoBase[T]) Update(ctx context.Context, entity T, opts ...idao.CallOptions) *idao.Result {
	if d.IsNil(entity) {
		panic(fmt.Errorf("Dao.Update() entity is nil"))
	}
	tenantId := d.GetTenantId(ctx, opts...)
	d.store.SetTenantId(entity, tenantId)
	res := d.store.Update(ctx, entity, idao.NewCallOptions(opts...))
	if res.Error != nil {
		panic(res.Error)
	}

	//d.PublishEvent(ctx, idao.AccessTypeUpdate, entity, opts...)
	return idao.NewResult(res)
}

func (d *DaoBase[T]) UpdateNotNull(ctx context.Context, entity T, opts ...idao.CallOptions) *idao.Result {
	options := idao.NewCallOptions(opts...)
	options.SetNotUpdateNull(true)
	return d.Update(ctx, entity, options)
}

func (d *DaoBase[T]) UpdateMap(ctx context.Context, id string, data map[string]any, opts ...idao.CallOptions) *idao.Result {
	tenantId := d.GetTenantId(ctx, opts...)
	res := d.store.UpdateMap(ctx, tenantId, id, data, idao.NewCallOptions(opts...))
	if res.Error != nil {
		panic(res.Error)
	}
	//d.PublishEvent(ctx, idao.AccessTypeUpdate, data, opts...)
	return idao.NewResult(res)
}

func (d *DaoBase[T]) UpdateMany(ctx context.Context, list []T, opts ...idao.CallOptions) *idao.Result {
	tenantId := d.GetTenantId(ctx, opts...)
	res := d.store.UpdateMany(ctx, tenantId, list, idao.NewCallOptions(opts...))
	if res.Error != nil {
		panic(res.Error)
	}
	return idao.NewResult(res)
}

func (d *DaoBase[T]) UpdateByRSQL(ctx context.Context, filterRSQL string, data T, opts ...idao.CallOptions) *idao.Result {
	tenantId := d.GetTenantId(ctx, opts...)
	res := d.store.UpdateByRSQL(ctx, tenantId, filterRSQL, data, idao.NewCallOptions(opts...))
	if res.Error != nil {
		panic(res.Error)
	}
	return idao.NewResult(res)
}

func (d *DaoBase[T]) UpdateMapByRSQL(ctx context.Context, rsql string, data map[string]any, opts ...idao.CallOptions) *idao.Result {
	tenantId := d.GetTenantId(ctx, opts...)
	res := d.store.UpdateMapByRSQL(ctx, tenantId, rsql, data, idao.NewCallOptions(opts...))
	if res.Error != nil {
		panic(res.Error)
	}
	return idao.NewResult(res)
}
func (d *DaoBase[T]) IsNil(entity T) bool {
	return any(entity) == nil
}
