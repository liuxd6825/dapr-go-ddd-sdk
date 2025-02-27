package impl

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
)

func (d *DaoBase) SoftDeleteById(ctx context.Context, id string, opts ...*db.CallOptions) {
	entity := map[string]any{}
	entity[Id] = id
	entity[TenantId] = d.GetTenantId(ctx)
	err := d.dao.Update(ctx, entity, db.NewRepositoryOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}

	d.PublishEvent(ctx, db.AccessTypeUpdate, entity, opts...)
}

func (d *DaoBase) DeleteById(ctx context.Context, id string, opts ...*db.CallOptions) {
	tenantId := d.GetTenantId(ctx)

	err := d.dao.DeleteById(ctx, tenantId, id, db.NewRepositoryOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}

	if d.GetIsPubEvent() {
		entity := map[string]any{
			TenantId: tenantId,
			Id:       id,
		}
		d.PublishEvent(ctx, db.AccessTypeDelete, entity, opts...)
	}
}

func (d *DaoBase) DeleteByIds(ctx context.Context, ids []string, opts ...*db.CallOptions) {
	if d.GetIsPubEvent() {
		for _, id := range ids {
			d.DeleteById(ctx, id, opts...)
		}
		return
	}

	tenantId := d.GetTenantId(ctx)
	err := d.dao.DeleteByIds(ctx, tenantId, ids, db.NewRepositoryOptions(opts)...)
	if err != nil {
		panic(err)
	}

	if d.GetIsPubEvent() {
		list := []map[string]any{}
		for _, id := range ids {
			e := map[string]any{}
			e[Id] = id
			e[TenantId] = tenantId
			list = append(list, e)
		}

		d.PublishBatchEvent(ctx, db.AccessTypeBatchDelete, list, opts...)
	}
}

func (d *DaoBase) DeleteAll(ctx context.Context, opts ...*db.CallOptions) {
	tenantId := d.GetTenantId(ctx)
	err := d.dao.DeleteAll(ctx, tenantId, db.NewRepositoryOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}

func (d *DaoBase) DeleteByFilter(ctx context.Context, filterRSQL string, opts ...*db.CallOptions) {
	tenantId := d.GetTenantId(ctx)
	err := d.dao.DeleteByFilter(ctx, tenantId, filterRSQL, db.NewRepositoryOptions(opts)...)
	if err != nil {
		panic(err)
	}
}

func (d *DaoBase) DeleteByMap(ctx context.Context, filterMap map[string]interface{}, opts ...*db.CallOptions) {
	tenantId := d.GetTenantId(ctx)
	err := d.dao.DeleteByMap(ctx, tenantId, filterMap, db.NewRepositoryOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}
