package impl

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
)

func (d *DaoBase) Delete(ctx context.Context, entity map[string]any, opts ...*db.CallOptions) *db.Result {
	id := d.dao.GetId(entity)
	return d.DeleteById(ctx, id, opts...)
}

func (d *DaoBase) DeleteById(ctx context.Context, id string, opts ...*db.CallOptions) *db.Result {
	tenantId := d.GetTenantId(ctx)

	res := d.dao.DeleteById(ctx, tenantId, id, db.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}

	if d.GetIsPubEvent() {
		entity := map[string]any{
			TenantId: tenantId,
			Id:       id,
		}
		d.PublishEvent(ctx, db.AccessTypeDelete, entity, opts...)
	}
	return db.NewResult(res)
}

func (d *DaoBase) deleteById(ctx context.Context, id string, opts ...*db.CallOptions) int64 {
	tenantId := d.GetTenantId(ctx)

	res := d.dao.DeleteById(ctx, tenantId, id, db.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}

	if d.GetIsPubEvent() {
		entity := map[string]any{
			TenantId: tenantId,
			Id:       id,
		}
		d.PublishEvent(ctx, db.AccessTypeDelete, entity, opts...)
	}
	return res.RowsAffected
}

func (d *DaoBase) DeleteByIds(ctx context.Context, ids []string, opts ...*db.CallOptions) *db.Result {

	if d.GetIsPubEvent() {
		count := int64(0)
		for _, id := range ids {
			count += d.deleteById(ctx, id, opts...)
		}
		return db.NewResult(nil).SetRowsAffected(count)
	}

	tenantId := d.GetTenantId(ctx)
	res := d.dao.DeleteByIds(ctx, tenantId, ids, db.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}

	if d.GetIsPubEvent() {
		var list []map[string]any
		for _, id := range ids {
			e := map[string]any{}
			e[Id] = id
			e[TenantId] = tenantId
			list = append(list, e)
		}

		d.PublishBatchEvent(ctx, db.AccessTypeBatchDelete, list, opts...)
	}
	return db.NewResult(res)
}

func (d *DaoBase) DeleteAll(ctx context.Context, opts ...*db.CallOptions) *db.Result {
	tenantId := d.GetTenantId(ctx)
	res := d.dao.DeleteAll(ctx, tenantId, db.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}
	return db.NewResult(res)
}

func (d *DaoBase) DeleteByRSQL(ctx context.Context, filterRSQL string, opts ...*db.CallOptions) *db.Result {
	tenantId := d.GetTenantId(ctx)
	res := d.dao.DeleteByRSQL(ctx, tenantId, filterRSQL, db.NewRepositoryOptions(opts)...)
	return db.NewResult(res)
}

/*
func (d *DaoBase) DeleteByMap(ctx context.Context, filterMap map[string]interface{}, opts ...*db.CallOptions) *db.Result {
	tenantId := d.GetTenantId(ctx)
	res := d.dao.DeleteByMap(ctx, tenantId, filterMap, db.NewRepositoryOptions(opts)...)
	return db.NewResult(res)
}
*/
