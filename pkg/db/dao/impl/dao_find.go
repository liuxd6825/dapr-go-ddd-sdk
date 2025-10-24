package impl

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

func (d *DaoBase[T]) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (T, error) {
	tenantId := d.GetTenantId(ctx, opts...)
	res := d.store.FindById(ctx, tenantId, id, idao.NewCallOptions(opts...))
	return res.Data, res.Error
}

func (d *DaoBase[T]) FindByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) ([]T, error) {
	tenantId := d.GetTenantId(ctx, opts...)
	data, _, err := d.store.FindByIds(ctx, tenantId, ids, idao.NewCallOptions(opts...)).Result()
	return data, err
}

func (d *DaoBase[T]) FindOneByRSQL(ctx context.Context, rsql string, opts ...idao.CallOptions) (T, error) {
	var null T
	tenantId := d.GetTenantId(ctx, opts...)
	list := d.store.FindByRSQL(ctx, tenantId, rsql, idao.NewCallOptions(opts...))
	if list.Error != nil {
		return null, list.Error
	}
	if len(list.Data) == 0 {
		return null, nil
	}
	return list.Data[0], nil
}

func (d *DaoBase[T]) FindByRSQL(ctx context.Context, rsql string, opts ...idao.CallOptions) ([]T, error) {
	tenantId := d.GetTenantId(ctx, opts...)
	res := d.store.FindByRSQL(ctx, tenantId, rsql, idao.NewCallOptions(opts...))
	return res.GetData(), res.GetError()
}

func (d *DaoBase[T]) FindAll(ctx context.Context, opts ...idao.CallOptions) *store.FindListResult[T] {
	tenantId := d.GetTenantId(ctx, opts...)
	res := d.store.FindAll(ctx, tenantId, idao.NewCallOptions(opts...))
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

/*
func (d *DaoBase) FindListByMap(ctx context.Context, filterMap map[string]interface{}, opts ...*db.CallOptions) *ddd_repository.FindListResult[map[string]any] {
	tenantId := d.GetTenantId(ctx, opts...)
	res := d.dao.FindListByMap(ctx, tenantId, filterMap, db.NewCallOptions(opts...))
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}
*/

func (d *DaoBase[T]) FindPaging(ctx context.Context, findPaging store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[T] {
	findQuery := d.NewFindPagingQuery(ctx, findPaging)
	res := d.store.FindPaging(ctx, findQuery, idao.NewCallOptions(opts...))
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *DaoBase[T]) FindAutoComplete(ctx context.Context, qry store.FindAutoCompleteQuery, opts ...idao.CallOptions) store.FindPagingResult[T] {
	if qry == nil {
		panic(errors.New("FindAutoComplete query is nil"))
	}
	tenantId := d.GetTenantId(ctx, opts...)
	qry.SetTenantId(tenantId)
	res := d.store.FindAutoComplete(ctx, qry, idao.NewCallOptions(opts...))
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *DaoBase[T]) FindDistinct(ctx context.Context, qry store.FindDistinctQuery, opts ...idao.CallOptions) store.FindPagingResult[T] {
	if qry == nil {
		panic(errors.New("FindDistinctQueryRequest query is nil"))
	}
	data := d.store.FindDistinct(ctx, qry, idao.NewCallOptions(opts...))
	if data.GetError() != nil {
		panic(data.GetError())
	}
	return data
}
