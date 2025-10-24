package impl

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

//type OnCreate func(ctx context.Context, entity any, opts ...*idao.CallOptions)
//var OnCreateBefore []OnCreate

func (d *DaoBase[T]) Create(ctx context.Context, entity T, opts ...idao.CallOptions) *idao.Result {
	if d.IsNil(entity) {
		panic(fmt.Errorf("Dao.Create() entity is nil"))
	}
	tenantId := d.GetTenantId(ctx, opts...)
	d.store.SetTenantId(entity, tenantId)

	res := d.store.Insert(ctx, entity, idao.NewCallOptions(opts...))
	if res.Error != nil {
		panic(fmt.Sprintf("tableName: %s; %s ", d.tableName, res.Error))
	}

	//d.PublishEvent(ctx, idao.AccessTypeCreate, entity, opts...)
	return idao.NewResult(res)
}

type nilRowsAffected struct {
}

func (r *nilRowsAffected) GetRowsAffected() int64 {
	return 0
}

func (d *DaoBase[T]) CreateMany(ctx context.Context, list []T, opts ...idao.CallOptions) *idao.Result {
	if len(list) == 0 {
		return idao.NewResult(&nilRowsAffected{})
	}
	tenantId := d.GetTenantId(ctx, opts...)
	for _, entity := range list {
		d.store.SetTenantId(entity, tenantId)
	}
	res := d.store.InsertMany(ctx, tenantId, list, idao.NewCallOptions(opts...))
	if res.Error != nil {
		panic(res.Error)
	}
	return idao.NewResult(res)
}

func (d *DaoBase[T]) CreateUpdate(ctx context.Context, entity T, opts ...idao.CallOptions) *idao.Result {
	tenantId := d.GetTenantId(ctx, opts...)
	d.store.SetTenantId(entity, tenantId)
	res := d.store.InsertOrUpdate(ctx, entity, idao.NewCallOptions(opts...))
	return idao.NewResult(res)
}

func (d *DaoBase[T]) Merge(ctx context.Context, entity T, fields map[string]string, opts ...idao.CallOptions) *idao.Result {
	tenantId := d.GetTenantId(ctx, opts...)
	d.store.SetTenantId(entity, tenantId)
	res := d.store.Merge(ctx, entity, fields, idao.NewCallOptions(opts...))
	return idao.NewResult(res)
}
