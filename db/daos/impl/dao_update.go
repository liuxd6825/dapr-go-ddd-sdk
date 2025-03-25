package impl

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos/idao"
)

func (d *DaoBase[T]) Update(ctx context.Context, entity T, opts ...*idao.CallOptions) *idao.Result {
	if d.IsNil(entity) {
		panic(fmt.Errorf("Dao.Update() entity is nil"))
	}
	tenantId := d.GetTenantId(ctx)
	d.store.SetTenantId(entity, tenantId)
	res := d.store.Update(ctx, entity, idao.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}

	d.PublishEvent(ctx, idao.AccessTypeUpdate, entity, opts...)
	return idao.NewResult(res)
}

func (d *DaoBase[T]) UpdateMap(ctx context.Context, id string, data map[string]any, opts ...*idao.CallOptions) *idao.Result {
	tenantId := d.GetTenantId(ctx)
	res := d.store.UpdateMap(ctx, tenantId, id, data, idao.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}
	//d.PublishEvent(ctx, idao.AccessTypeUpdate, data, opts...)
	return idao.NewResult(res)
}

func (d *DaoBase[T]) UpdateMany(ctx context.Context, list []T, opts ...*idao.CallOptions) *idao.Result {
	tenantId := d.GetTenantId(ctx)
	res := d.store.UpdateMany(ctx, tenantId, list, idao.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}
	return idao.NewResult(res)
}

func (d *DaoBase[T]) UpdateByRSQL(ctx context.Context, filterRSQL string, data T, opts ...*idao.CallOptions) *idao.Result {
	tenantId := d.GetTenantId(ctx)
	res := d.store.UpdateByRSQL(ctx, tenantId, filterRSQL, data, idao.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}
	return idao.NewResult(res)
}

func (d *DaoBase[T]) IsNil(entity T) bool {
	return any(entity) == nil
}
