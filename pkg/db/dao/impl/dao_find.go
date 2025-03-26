package impl

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

func (d *DaoBase[T]) FindById(ctx context.Context, id string, opts ...*idao.CallOptions) T {
	tenantId := d.GetTenantId(ctx)
	res := d.store.FindById(ctx, tenantId, id, idao.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.Data
}

func (d *DaoBase[T]) FindByIds(ctx context.Context, ids []string, opts ...*idao.CallOptions) []T {
	tenantId := d.GetTenantId(ctx)
	data, _, err := d.store.FindByIds(ctx, tenantId, ids, idao.NewRepositoryOptions(opts)...).Result()
	if err != nil {
		panic(err)
	}
	return data
}

func (d *DaoBase[T]) FindOneByRSQL(ctx context.Context, rsql string, opts ...*idao.CallOptions) T {
	tenantId := d.GetTenantId(ctx)
	list := d.store.FindByRSQL(ctx, tenantId, rsql, idao.NewRepositoryOptions(opts)...)
	if list.Error != nil {
		panic(list.Error)
	}
	if len(list.Data) == 0 {
		var null T
		return null
	}
	return list.Data[0]
}

func (d *DaoBase[T]) FindByRSQL(ctx context.Context, rsql string, opts ...*idao.CallOptions) []T {
	tenantId := d.GetTenantId(ctx)
	res := d.store.FindByRSQL(ctx, tenantId, rsql, idao.NewRepositoryOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res.Data
}

func (d *DaoBase[T]) FindAll(ctx context.Context, opts ...*idao.CallOptions) *store.FindListResult[T] {
	tenantId := d.GetTenantId(ctx)
	res := d.store.FindAll(ctx, tenantId, idao.NewRepositoryOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

/*
func (d *DaoBase) FindListByMap(ctx context.Context, filterMap map[string]interface{}, opts ...*db.CallOptions) *ddd_repository.FindListResult[map[string]any] {
	tenantId := d.GetTenantId(ctx)
	res := d.dao.FindListByMap(ctx, tenantId, filterMap, db.NewRepositoryOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}
*/

func (d *DaoBase[T]) FindPaging(ctx context.Context, findPaging *store.FindPagingQueryRequest, opts ...*idao.CallOptions) *store.FindPagingResult[T] {
	findQuery := d.NewFindPagingQuery(ctx, findPaging)
	res := d.store.FindPaging(ctx, findQuery, idao.NewRepositoryOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *DaoBase[T]) FindAutoComplete(ctx context.Context, qry *store.FindAutoCompleteQueryRequest, opts ...*idao.CallOptions) *store.FindPagingResult[T] {
	if qry == nil {
		panic(errors.New("FindAutoComplete query is nil"))
	}
	tenantId := d.GetTenantId(ctx)
	qry.SetTenantId(tenantId)
	res := d.store.FindAutoComplete(ctx, qry, idao.NewRepositoryOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *DaoBase[T]) FindDistinct(ctx context.Context, qry *store.FindDistinctQueryRequest, opts ...*idao.CallOptions) *store.FindPagingResult[T] {
	if qry == nil {
		panic(errors.New("FindDistinctQueryRequest query is nil"))
	}
	data := d.store.FindDistinct(ctx, qry, idao.NewRepositoryOptions(opts)...)
	if data.GetError() != nil {
		panic(data.GetError())
	}
	return data
}
