package ddd_neo4j

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
)

func (d *Dao[T]) getUpdateMap(ctx context.Context, entity T) map[string]any {
	res, err := d.schema.NewMap(context.Background(), entity)
	if err != nil {
		panic(err)
	}
	d.eb.SetUpdatedInfo(ctx, res)
	return res
}

func (d *Dao[T]) getUpdateMapList(ctx context.Context, list []T) []map[string]any {
	items := make([]map[string]any, 0)
	for _, ent := range list {
		m := d.getUpdateMap(ctx, ent)
		items = append(items, m)
	}
	return items
}

func (d *Dao[T]) UpdateByRSQL(ctx context.Context, tenantId, rSQL string, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResultEmpty[T]()
	gp.Try(func() error {
		cr, err := d.cypher.UpdateByRSQL(ctx, tenantId, rSQL, entity)
		if err != nil {
			return err
		}
		nRes, err := d.doSet(ctx, tenantId, cr.Cypher(), cr.Params(), opts...)
		if nRes != nil {
			res.SetRowsAffected(nRes.GetRowsAffected())
		}
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) UpdateMap(ctx context.Context, tenantId string, id string, data map[string]any, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) UpdateMapAndGetCount(ctx context.Context, tenantId string, filter any, data any, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) Update(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResultEmpty[T]()
	gp.Try(func() error {
		tenantId := d.GetTenantId(entity)
		cr, err := d.cypher.Update(ctx, tenantId, entity)
		if err != nil {
			return err
		}

		nRes, err := d.doSet(ctx, tenantId, cr.Cypher(), cr.Params(), opts...)
		if nRes != nil {
			res.SetRowsAffected(nRes.GetRowsAffected())
		}
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) UpdateMany(ctx context.Context, tenantId string, list []T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResultEmpty[T]()
	gp.Try(func() error {
		for _, ent := range list {
			d.eb.SetTenantId(ent, tenantId)
			d.eb.SetUpdatedInfo(ctx, ent)
		}
		cr, err := d.cypher.UpdateMany(ctx, tenantId, list)
		if err != nil {
			return err
		}
		cypher := cr.Cypher()
		params := cr.Params()
		nRes, err := d.doSet(ctx, tenantId, cypher, params, opts...)
		if nRes != nil {
			res.SetRowsAffected(nRes.GetRowsAffected())
		}
		if err != nil {
			println(cypher)
		}
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) UpdateLabelById(ctx context.Context, tenantId string, id string, label string) error {
	cr, err := d.cypher.UpdateLabelById(ctx, tenantId, id, label)
	if err != nil {
		return err
	}
	_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params())
	return err
}

func (d *Dao[T]) UpdateLabelByFilter(ctx context.Context, tenantId string, filter string, labels ...string) error {
	cr, err := d.cypher.UpdateLabelByFilter(ctx, tenantId, filter, labels...)
	if err != nil {
		return err
	}
	if cr == nil {
		return nil
	}

	_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params())
	if err != nil {
		return err
	}
	return nil
}
