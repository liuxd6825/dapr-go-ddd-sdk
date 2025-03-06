package ddd_neo4j

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
)

func (d *Dao[T]) DeleteLabelById(ctx context.Context, tenantId string, id string, label string) error {
	cr, err := d.cypher.DeleteLabelById(ctx, tenantId, id, label)
	if err != nil || cr == nil {
		return err
	}
	_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params())
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao[T]) DeleteLabelByFilter(ctx context.Context, tenantId string, filter string, labels ...string) error {
	cr, err := d.cypher.DeleteLabelByFilter(ctx, tenantId, filter, labels...)
	if err != nil || cr == nil {
		return err
	}
	_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params())
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao[T]) DeleteById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResultEmpty[T]()
	gp.Try(func() error {
		cr, err := d.cypher.DeleteById(ctx, tenantId, id)
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

func (d *Dao[T]) DeleteByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResultEmpty[T]()
	gp.Try(func() error {
		cr, err := d.cypher.DeleteByIds(ctx, tenantId, ids)
		if err != nil {
			return err
		}
		_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params(), opts...)
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) DeleteAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResultEmpty[T]()
	gp.Try(func() error {
		cr, err := d.cypher.DeleteAll(ctx, tenantId)
		if err != nil {
			return err
		}
		_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params(), opts...)
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) DeleteByFilter(ctx context.Context, tenantId string, filter string, opts ...ddd_repository.Options) error {
	cr, err := d.cypher.DeleteByFilter(ctx, tenantId, filter)
	if err != nil {
		return err
	}
	_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params(), opts...)
	return err
}

func (d *Dao[T]) DeleteByGraphId(ctx context.Context, tenantId string, graphId string, opts ...ddd_repository.Options) error {
	return d.DeleteByFilter(ctx, tenantId, fmt.Sprintf("graphId=='%v'", graphId))
}

func (d *Dao[T]) DeleteByCaseId(ctx context.Context, tenantId string, caseId string, opts ...ddd_repository.Options) error {
	return d.DeleteByFilter(ctx, tenantId, fmt.Sprintf("caseId=='%v'", caseId))
}

func (d *Dao[T]) DeleteByTenantId(ctx context.Context, tenantId string, opts ...ddd_repository.Options) error {
	cr, err := d.cypher.DeleteByTenantId(ctx, tenantId)
	if err != nil {
		return err
	}
	_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params(), opts...)
	return err
}

func (d *Dao[T]) Delete(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) DeleteByRSQL(ctx context.Context, tenantId, filter string, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) DeleteByMap(ctx context.Context, tenantId string, filterMap map[string]any, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	//TODO implement me
	panic("implement me")
}
