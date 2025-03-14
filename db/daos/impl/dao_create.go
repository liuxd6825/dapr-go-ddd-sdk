package impl

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos/idao"
)

func (d *DaoBase[T]) Create(ctx context.Context, entity T, opts ...*idao.CallOptions) *idao.Result {
	if d.IsNil(entity) {
		panic(fmt.Errorf("Dao.Create() entity is nil"))
	}
	tenantId := d.GetTenantId(ctx)
	d.dao.SetTenantId(entity, tenantId)
	res := d.dao.Insert(ctx, entity, idao.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}
	d.PublishEvent(ctx, idao.AccessTypeCreate, entity, opts...)
	return idao.NewResult(res)
}

func (d *DaoBase[T]) CreateMany(ctx context.Context, list []T, opts ...*idao.CallOptions) *idao.Result {
	tenantId := d.GetTenantId(ctx)
	res := d.dao.InsertMany(ctx, tenantId, list, idao.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}
	return idao.NewResult(res)
}
