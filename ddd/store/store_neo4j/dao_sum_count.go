package store_neo4j

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
)

func (d *Dao[T]) SumEntity(ctx context.Context, qry store.FindPagingQuery, opts ...store.Options) ([]T, bool, error) {
	data := d.NewEntityList()
	_, found, err := d.SumByQuery(ctx, qry, &data, opts...)
	return data, found, err
}

func (d *Dao[T]) SumByQuery(ctx context.Context, qry store.FindPagingQuery, resData any, opts ...store.Options) (any, bool, error) {
	var err error
	if len(qry.GetValueCols()) == 0 {
		return nil, false, nil
	}
	f1 := qry.GetFilter()
	f2 := qry.GetMustFilter()
	f3 := ""
	mustWhere, ok := qry.(store.FindPagingQueryMustWhere)
	if ok {
		f3, err = mustWhere.GetMustWhere()
		if err != nil {
			return nil, false, err
		}
	}
	tenantId := appctx.GetTenantId2(ctx)
	filter := getSqlAnds(f1, f2, f3)
	res, found, err := d.sum(ctx, tenantId, filter, qry.GetValueCols(), resData, opts...)
	return res, found, err
}

func (d *Dao[T]) SumByRSQL(ctx context.Context, tenantId, rSql string, valueCols []*store.ValueCol, opts ...store.Options) map[string]any {
	data := map[string]any{}
	_, _, err := d.sum(ctx, tenantId, rSql, valueCols, data, opts...)
	if err != nil {
		panic(err)
	}
	return data
}

func (d *Dao[T]) sum(ctx context.Context, tenantId, rSql string, valueCols []*store.ValueCol, resData any, opts ...store.Options) (any, bool, error) {
	cr, err := d.Cypher.Sum(ctx, tenantId, rSql, valueCols)
	if err != nil {
		return nil, false, err
	}
	cypher := cr.Cypher()
	result, err := d.Query(ctx, cypher, cr.Params())
	if err != nil {
		return nil, false, err
	}
	err = result.GetSum(resData)
	return resData, resData != nil, err
}

func (d *Dao[T]) CountByRSQL(ctx context.Context, tenantId string, rSql string, opts ...store.Options) (int64, error) {
	cr, err := d.Cypher.Count(ctx, tenantId, rSql)
	if err != nil {
		return 0, err
	}
	cypher := cr.Cypher()
	result, err := d.Query(ctx, cypher, cr.Params())
	if err != nil {
		return 0, err
	}
	count := result.GetInt("rows")
	return count, nil
}

func getSqlAnds(s ...string) string {
	res := ""
	for _, item := range s {
		res = getSqlAnd(res, item)
	}
	return res
}

func getSqlAnd(s1 string, s2 string) string {
	b1 := len(s1) > 0
	b2 := len(s2) > 0
	if b1 && b2 {
		return fmt.Sprintf("(%s) and (%s)", s1, s2)
	} else if b1 {
		return s1
	}
	return s2
}
