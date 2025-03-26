package impl

import (
	"context"
	"fmt"
	idao2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

func (d *DaoBase[T]) Create(ctx context.Context, entity T, opts ...*idao2.CallOptions) *idao2.Result {
	if d.IsNil(entity) {
		panic(fmt.Errorf("Dao.Create() entity is nil"))
	}
	tenantId := d.GetTenantId(ctx)
	d.store.SetTenantId(entity, tenantId)
	res := d.store.Insert(ctx, entity, idao2.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}
	d.PublishEvent(ctx, idao2.AccessTypeCreate, entity, opts...)
	return idao2.NewResult(res)
}

func (d *DaoBase[T]) CreateMany(ctx context.Context, list []T, opts ...*idao2.CallOptions) *idao2.Result {
	tenantId := d.GetTenantId(ctx)
	res := d.store.InsertMany(ctx, tenantId, list, idao2.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}
	return idao2.NewResult(res)
}
