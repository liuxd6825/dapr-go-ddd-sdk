package impl

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
)

func (d *DaoBase) Create(ctx context.Context, entity map[string]any, opts ...*db.CallOptions) int64 {
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
	return res.RowsAffected
}

func (d *DaoBase) CreateMany(ctx context.Context, entity []map[string]any, opts ...*db.CallOptions) int64 {
	tenantId := d.GetTenantId(ctx)
	for _, e := range entity {
		e[TenantId] = tenantId
	}
	res := d.dao.InsertMany(ctx, entity, db.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected
}
