package impl

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
)

func (d *DaoBase) FindById(ctx context.Context, id string, opts ...*db.CallOptions) map[string]any {
	tenantId := d.GetTenantId(ctx)
	res := d.dao.FindById(ctx, tenantId, id, db.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.Data
}

func (d *DaoBase) FindByIds(ctx context.Context, ids []string, opts ...*db.CallOptions) []map[string]any {
	tenantId := d.GetTenantId(ctx)
	data, _, err := d.dao.FindByIds(ctx, tenantId, ids, db.NewRepositoryOptions(opts)...).Result()
	if err != nil {
		panic(err)
	}
	return data
}

func (d *DaoBase) FindByRSQL(ctx context.Context, rsql string, opts ...*db.CallOptions) []map[string]any {
	tenantId := d.GetTenantId(ctx)
	res := d.dao.FindByRSQL(ctx, tenantId, rsql, db.NewRepositoryOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res.Data
}

func (d *DaoBase) FindAll(ctx context.Context, opts ...*db.CallOptions) *ddd_repository.FindListResult[map[string]any] {
	tenantId := d.GetTenantId(ctx)
	res := d.dao.FindAll(ctx, tenantId, db.NewRepositoryOptions(opts)...)
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

func (d *DaoBase) FindPaging(ctx context.Context, findPaging *ddd_repository.FindPagingQueryRequest, opts ...*db.CallOptions) *ddd_repository.FindPagingResult[map[string]any] {
	findQuery := d.NewFindPagingQuery(ctx, findPaging)
	res := d.dao.FindPaging(ctx, findQuery, db.NewRepositoryOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *DaoBase) FindAutoComplete(ctx context.Context, qry *ddd_repository.FindAutoCompleteQueryRequest, opts ...*db.CallOptions) *ddd_repository.FindPagingResult[map[string]any] {
	if qry == nil {
		panic(errors.New("FindAutoComplete query is nil"))
	}
	tenantId := d.GetTenantId(ctx)
	qry.SetTenantId(tenantId)
	res := d.dao.FindAutoComplete(ctx, qry, db.NewRepositoryOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *DaoBase) FindDistinct(ctx context.Context, qry *ddd_repository.FindDistinctQueryRequest, opts ...*db.CallOptions) *ddd_repository.FindPagingResult[map[string]any] {
	if qry == nil {
		panic(errors.New("FindDistinctQueryRequest query is nil"))
	}
	data := d.dao.FindDistinct(ctx, qry, db.NewRepositoryOptions(opts)...)
	if data.GetError() != nil {
		panic(data.GetError())
	}
	return data
}
