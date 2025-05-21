package store_neo4j

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
)

func (d *Dao[T]) DeleteLabelById(ctx context.Context, tenantId string, id string, label string) error {
	cr, err := d.Cypher.DeleteLabelById(ctx, tenantId, id, label)
	if err != nil || cr == nil {
		return err
	}
	_, err = d.DoSet(ctx, tenantId, cr.Cypher(), cr.Params())
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao[T]) DeleteLabelByFilter(ctx context.Context, tenantId string, filter string, labels ...string) error {
	cr, err := d.Cypher.DeleteLabelByFilter(ctx, tenantId, filter, labels...)
	if err != nil || cr == nil {
		return err
	}
	_, err = d.DoSet(ctx, tenantId, cr.Cypher(), cr.Params())
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao[T]) DeleteById(ctx context.Context, tenantId string, id string, opts ...store.Options) *store.SetResult[T] {
	res := store.NewSetResultEmpty[T]()
	gp.Try(func() error {
		cr, err := d.Cypher.DeleteById(ctx, tenantId, id)
		if err != nil {
			return err
		}
		nRes, err := d.DoSet(ctx, tenantId, cr.Cypher(), cr.Params(), opts...)
		if nRes != nil {
			res.SetRowsAffected(nRes.GetRowsAffected())
		}
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) DeleteByIds(ctx context.Context, tenantId string, ids []string, opts ...store.Options) *store.SetResult[T] {
	res := store.NewSetResultEmpty[T]()
	gp.Try(func() error {
		cr, err := d.Cypher.DeleteByIds(ctx, tenantId, ids)
		if err != nil {
			return err
		}
		_, err = d.DoSet(ctx, tenantId, cr.Cypher(), cr.Params(), opts...)
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) DeleteAll(ctx context.Context, tenantId string, opts ...store.Options) *store.SetResult[T] {
	res := store.NewSetResultEmpty[T]()
	gp.Try(func() error {
		cr, err := d.Cypher.DeleteAll(ctx, tenantId)
		if err != nil {
			return err
		}
		_, err = d.DoSet(ctx, tenantId, cr.Cypher(), cr.Params(), opts...)
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) DeleteByFilter(ctx context.Context, tenantId string, filter string, opts ...store.Options) error {
	cr, err := d.Cypher.DeleteByRSQL(ctx, tenantId, filter)
	if err != nil {
		return err
	}
	_, err = d.DoSet(ctx, tenantId, cr.Cypher(), cr.Params(), opts...)
	return err
}

func (d *Dao[T]) DeleteByGraphId(ctx context.Context, tenantId string, graphId string, opts ...store.Options) error {
	return d.DeleteByFilter(ctx, tenantId, fmt.Sprintf("graphId=='%v'", graphId))
}

func (d *Dao[T]) DeleteByCaseId(ctx context.Context, tenantId string, caseId string, opts ...store.Options) error {
	return d.DeleteByFilter(ctx, tenantId, fmt.Sprintf("caseId=='%v'", caseId))
}

func (d *Dao[T]) DeleteByTenantId(ctx context.Context, tenantId string, opts ...store.Options) error {
	cr, err := d.Cypher.DeleteByTenantId(ctx, tenantId)
	if err != nil {
		return err
	}
	_, err = d.DoSet(ctx, tenantId, cr.Cypher(), cr.Params(), opts...)
	return err
}

func (d *Dao[T]) Delete(ctx context.Context, entity T, opts ...store.Options) *store.SetResult[T] {
	id := d.GetId(entity)
	tenantId := d.GetTenantId(entity)
	return d.DeleteById(ctx, tenantId, id)
}

func (d *Dao[T]) DeleteByRSQL(ctx context.Context, tenantId, rSQL string, opts ...store.Options) *store.SetResult[T] {
	res := store.NewSetResult[T]()
	gp.Try(func() error {
		cr, err := d.Cypher.DeleteByRSQL(ctx, tenantId, rSQL)
		if err != nil {
			return err
		}
		cypher := cr.Cypher()
		nRes, err := d.DoSet(ctx, tenantId, cypher, cr.Params(), opts...)
		if err == nil {
			res.SetRowsAffected(nRes.GetRowsAffected())
		}
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res

}
