package ddd_neo4j

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
)

func (d *Dao[T]) InsertMany(ctx context.Context, entities []T, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T] {
	res := ddd_repository.NewSetManyResultEmpty[T]()
	gp.Try(func() error {
		for _, e := range entities {
			if err := d.Insert(ctx, e, opts...).GetError(); err != nil {
				return err
			}
		}
		return nil
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) InsertOrUpdate(ctx context.Context, entity T, opts ...ddd_repository.Options) (setResult *ddd_repository.SetResult[T]) {
	var err error
	defer func() {
		if err = errors.GetRecoverError(err, recover()); err != nil {
			setResult = ddd_repository.NewSetResultError[T](err)
		}
	}()

	cr, err1 := d.cypher.InsertOrUpdate(ctx, entity)
	if err1 != nil {
		err = err1
		return ddd_repository.NewSetResultError[T](err)
	}

	tenantId := d.eb.GetTenantId(entity)
	_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params(), opts...)
	if err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	return ddd_repository.NewSetResult(entity, err)
}

func (d *Dao[T]) InsertOrUpdateMany(ctx context.Context, entities []T, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T] {
	for _, e := range entities {
		if err := d.InsertOrUpdate(ctx, e, opts...).GetError(); err != nil {
			return ddd_repository.NewSetManyResultError[T](err)
		}
	}
	return ddd_repository.NewSetManyResult[T](entities, nil)
}

func (d *Dao[T]) Save(ctx context.Context, data *ddd.SetData[T], opts ...ddd_repository.Options) (setResult *ddd_repository.SetResult[T]) {
	var err error
	defer func() {
		if err = errors.GetRecoverError(err, recover()); err != nil {
			setResult = ddd_repository.NewSetResultError[T](err)
		}
	}()

	for _, item := range data.Items() {
		statue := item.Statue()
		entity := item.Data()

		switch statue {
		case ddd.DataStatueCreate:
			err = d.Insert(ctx, entity, opts...).GetError()
		case ddd.DataStatueUpdate:
			err = d.Update(ctx, entity, opts...).GetError()
		case ddd.DataStatueDelete:
			err = d.DeleteById(ctx, d.GetTenantId(entity), d.GetId(entity), opts...)
		case ddd.DataStatueCreateOrUpdate:
			err = d.InsertOrUpdate(ctx, entity, opts...).GetError()
		}
		if err != nil {
			return ddd_repository.NewSetResultError[T](err)
		}
	}
	return ddd_repository.NewSetResultError[T](nil)
}

func (d *Dao[T]) Insert(ctx context.Context, entity T, opts ...ddd_repository.Options) (res *ddd_repository.SetResult[T]) {
	res = ddd_repository.NewSetResultEmpty[T]()
	gp.Try(func() error {
		cr, err := d.cypher.Insert(ctx, entity)
		if err != nil {
			return err
		}
		d.eb.SetUpdatedInfo(ctx, entity)
		_, err = d.doSet(ctx, d.GetTenantId(entity), cr.Cypher(), cr.Params(), opts...)
		if err != nil {
			return err
		}
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) InsertMap(ctx context.Context, tenantId string, data map[string]interface{}, opts ...ddd_repository.Options) (res *ddd_repository.SetResult[T]) {
	//TODO implement me
	panic("implement me")
}
