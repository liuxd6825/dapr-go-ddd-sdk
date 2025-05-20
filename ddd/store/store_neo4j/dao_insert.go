package store_neo4j

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
)

func (d *Dao[T]) Insert(ctx context.Context, entity T, opts ...store.Options) (res *store.SetResult[T]) {
	res = store.NewSetResultEmpty[T]()
	gp.Try(func() error {
		tenantId := d.config.EntityBuilder.GetTenantId(entity)
		d.config.EntityBuilder.SetCreatedInfo(ctx, entity)
		cr, err := d.Cypher.Insert(ctx, tenantId, entity)
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

func (d *Dao[T]) InsertMany(ctx context.Context, tenantId string, list []T, opts ...store.Options) *store.SetResult[T] {
	res := store.NewSetResultEmpty[T]()
	gp.Try(func() error {
		for _, ent := range list {
			d.config.EntityBuilder.SetTenantId(ent, tenantId)
			d.config.EntityBuilder.SetCreatedInfo(ctx, ent)
		}
		cr, err := d.Cypher.InsertMany(ctx, tenantId, list)
		if err != nil {
			return err
		}
		cypher := cr.Cypher()
		params := cr.Params()
		println(cypher)
		nRes, err := d.Run(ctx, cypher, params, true, opts...)
		if nRes != nil {
			res.SetRowsAffected(nRes.GetRowsAffected())
		}
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) InsertOrUpdate(ctx context.Context, entity T, opts ...store.Options) (res *store.SetResult[T]) {
	res = store.NewSetResultEmpty[T]()
	gp.Try(func() error {
		cr, err := d.Cypher.InsertOrUpdate(ctx, entity)
		if err != nil {
			return err
		}

		tenantId := d.config.EntityBuilder.GetTenantId(entity)
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

func (d *Dao[T]) InsertOrUpdateMany(ctx context.Context, entities []T, opts ...store.Options) *store.SetResult[T] {
	for _, e := range entities {
		if err := d.InsertOrUpdate(ctx, e, opts...).GetError(); err != nil {
			return store.NewSetResultEmpty[T]().SetError(err)
		}
	}
	return store.NewSetResultEmpty[T]()
}

func (d *Dao[T]) Merge(ctx context.Context, entity T, fields map[string]string, opts ...store.Options) (res *store.SetResult[T]) {
	res = store.NewSetResultEmpty[T]()
	gp.Try(func() error {
		cr, err := d.Cypher.Merge(ctx, entity, fields)
		if err != nil {
			return err
		}

		tenantId := d.config.EntityBuilder.GetTenantId(entity)
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

func (d *Dao[T]) Save(ctx context.Context, data *ddd.SetData[T], opts ...store.Options) (setResult *store.SetResult[T]) {
	var err error
	defer func() {
		if err = errors.GetRecoverError(err, recover()); err != nil {
			setResult = store.NewSetResultError[T](err)
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
			err = d.DeleteById(ctx, d.GetTenantId(entity), d.GetId(entity), opts...).Error
		case ddd.DataStatueCreateOrUpdate:
			err = d.InsertOrUpdate(ctx, entity, opts...).GetError()
		}
		if err != nil {
			return store.NewSetResultError[T](err)
		}
	}
	return store.NewSetResultError[T](nil)
}

func (d *Dao[T]) InsertMap(ctx context.Context, tenantId string, data map[string]interface{}, opts ...store.Options) (res *store.SetResult[T]) {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) getInsertMap(ctx context.Context, entity T) map[string]any {
	res, err := d.config.DBSchema.NewMap(context.Background(), entity)
	if err != nil {
		panic(err)
	}
	d.config.EntityBuilder.SetCreatedInfo(ctx, res)
	return res
}

func (d *Dao[T]) getInsertMapList(ctx context.Context, list []T) []map[string]any {
	items := make([]map[string]any, 0)
	for _, ent := range list {
		m := d.getInsertMap(ctx, ent)
		items = append(items, m)
	}
	return items
}
