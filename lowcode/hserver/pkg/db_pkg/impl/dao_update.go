package impl

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
)

func (d *DaoBase) Update(ctx context.Context, entity map[string]any, opts ...*db.CallOptions) *db.Result {
	if entity == nil {
		panic(fmt.Errorf("Dao.Update() entity is nil"))
	}
	tenantId := d.GetTenantId(ctx)
	d.dao.SetTenantId(entity, tenantId)
	res := d.dao.Update(ctx, entity, db.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}

	d.PublishEvent(ctx, db.AccessTypeUpdate, entity, opts...)
	return db.NewResult(res)
}

func (d *DaoBase) UpdateMap(ctx context.Context, id string, data map[string]any, opts ...*db.CallOptions) *db.Result {
	tenantId := d.GetTenantId(ctx)
	res := d.dao.UpdateMap(ctx, tenantId, id, data, db.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}
	d.PublishEvent(ctx, db.AccessTypeUpdate, data, opts...)
	return db.NewResult(res)
}

func (d *DaoBase) UpdateMany(ctx context.Context, entities []map[string]any, opts ...*db.CallOptions) *db.Result {
	tenantId := d.GetTenantId(ctx)
	res := d.dao.UpdateMany(ctx, tenantId, entities, db.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}
	return db.NewResult(res)
}

func (d *DaoBase) UpdateByRSQL(ctx context.Context, filterRSQL string, data map[string]any, opts ...*db.CallOptions) *db.Result {
	tenantId := d.GetTenantId(ctx)
	res := d.dao.UpdateByRSQL(ctx, tenantId, filterRSQL, data, db.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}
	return db.NewResult(res)
}
