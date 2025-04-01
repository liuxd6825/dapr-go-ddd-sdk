package impl

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

//type OnCreate func(ctx context.Context, entity any, opts ...*idao.CallOptions)
//var OnCreateBefore []OnCreate

func (d *DaoBase[T]) Create(ctx context.Context, entity T, opts ...*idao.CallOptions) *idao.Result {
	if d.IsNil(entity) {
		panic(fmt.Errorf("Dao.Create() entity is nil"))
	}
	tenantId := d.GetTenantId(ctx)
	d.store.SetTenantId(entity, tenantId)

	res := d.store.Insert(ctx, entity, idao.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(fmt.Sprintf("tableName: %s; %s ", d.tableName, res.Error))
	}

	d.PublishEvent(ctx, idao.AccessTypeCreate, entity, opts...)
	return idao.NewResult(res)
}

func (d *DaoBase[T]) CreateMany(ctx context.Context, list []T, opts ...*idao.CallOptions) *idao.Result {
	tenantId := d.GetTenantId(ctx)

	for _, entity := range list {
		d.store.SetTenantId(entity, tenantId)
	}

	res := d.store.InsertMany(ctx, tenantId, list, idao.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}
	return idao.NewResult(res)
}
