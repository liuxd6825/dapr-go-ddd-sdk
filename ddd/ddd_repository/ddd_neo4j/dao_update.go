package ddd_neo4j

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
)

func (d *Dao[T]) UpdateByRSQL(ctx context.Context, tenantId, filterRSQL string, data map[string]any, opts ...ddd_repository.Options) *ddd_repository.SetManyCountResult {
	//TODO implement me
	panic("implement me")
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
	cr, err := d.cypher.Update(ctx, entity)
	res, err := d.doSet(ctx, entity.GetTenantId(), cr.Cypher(), cr.Params(), opts...)
	if err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	if _, err := res.GetOne("", entity); err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	return ddd_repository.NewSetResult(entity, err)
}

func (d *Dao[T]) UpdateMany(ctx context.Context, list []T, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T] {
	for _, entity := range list {
		if cr, err := d.cypher.Update(ctx, entity); err != nil {
			return ddd_repository.NewSetManyResultError[T](err)
		} else {
			if res, err := d.doSet(ctx, entity.GetTenantId(), cr.Cypher(), cr.Params(), opts...); err != nil {
				return ddd_repository.NewSetManyResultError[T](err)
			} else if _, err := res.GetOne(cr.ResultKeys()[0], entity); err != nil {
				return ddd_repository.NewSetManyResultError[T](err)
			}
		}
	}
	return ddd_repository.NewSetManyResult(list, nil)
}

func (d *Dao[T]) UpdateLabelById(ctx context.Context, tenantId string, id string, label string) error {
	cr, err := d.cypher.UpdateLabelById(ctx, tenantId, id, label)
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
