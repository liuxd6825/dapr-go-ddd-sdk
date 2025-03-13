package impl

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db"
)

func (d *DaoBase) Create(ctx context.Context, entity map[string]any, opts ...*db.CallOptions) *db.Result {
	if entity == nil {
		panic(fmt.Errorf("Dao.Create() entity is nil"))
	}
	tenantId := d.GetTenantId(ctx)
	d.dao.SetTenantId(entity, tenantId)
	res := d.dao.Insert(ctx, entity, db.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}
	d.PublishEvent(ctx, db.AccessTypeCreate, entity, opts...)
	return db.NewResult(res)
}

func (d *DaoBase) CreateMany(ctx context.Context, list []map[string]any, opts ...*db.CallOptions) *db.Result {
	tenantId := d.GetTenantId(ctx)
	res := d.dao.InsertMany(ctx, tenantId, list, db.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}
	return db.NewResult(res)
}
