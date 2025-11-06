package store_neo4j

import (
	"context"

	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
)

func (d *Dao[T]) Insert(ctx context.Context, entity T, opts ...store2.Options) (res *store2.SetResult[T]) {
	res = store2.NewSetResultEmpty[T]()
	gp.Try(func() error {
		tenantId := d.config.EntityBuilder.GetTenantId(entity)
		d.config.EntityBuilder.SetCreatedInfo(ctx, entity)
		cr, err := d.Cypher.Insert(ctx, tenantId, entity)
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

func (d *Dao[T]) InsertMany(ctx context.Context, tenantId string, list []T, opts ...store2.Options) *store2.SetResult[T] {
	res := store2.NewSetResultEmpty[T]()
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

func (d *Dao[T]) InsertOrUpdate(ctx context.Context, entity T, opts ...store2.Options) (res *store2.SetResult[T]) {
	res = store2.NewSetResultEmpty[T]()
	gp.Try(func() error {
		cr, err := d.Cypher.InsertOrUpdate(ctx, entity)
		if err != nil {
			return err
		}

		tenantId := d.config.EntityBuilder.GetTenantId(entity)
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

func (d *Dao[T]) InsertOrUpdateMany(ctx context.Context, entities []T, opts ...store2.Options) *store2.SetResult[T] {
	for _, e := range entities {
		if err := d.InsertOrUpdate(ctx, e, opts...).GetError(); err != nil {
			return store2.NewSetResultEmpty[T]().SetError(err)
		}
	}
	return store2.NewSetResultEmpty[T]()
}

func (d *Dao[T]) Merge(ctx context.Context, entity T, fields map[string]string, opts ...store2.Options) (res *store2.SetResult[T]) {
	res = store2.NewSetResultEmpty[T]()
	gp.Try(func() error {
		cr, err := d.Cypher.Merge(ctx, entity, fields)
		if err != nil {
			return err
		}

		tenantId := d.config.EntityBuilder.GetTenantId(entity)
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

func (d *Dao[T]) Save(ctx context.Context, data *ddd.SetData[T], opts ...store2.Options) (setResult *store2.SetResult[T]) {
	var err error
	defer func() {
		if err = errors.GetRecoverError(err, recover()); err != nil {
			setResult = store2.NewSetResultError[T](err)
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
			return store2.NewSetResultError[T](err)
		}
	}
	return store2.NewSetResultError[T](nil)
}

func (d *Dao[T]) InsertMap(ctx context.Context, tenantId string, data map[string]interface{}, opts ...store2.Options) (res *store2.SetResult[T]) {
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
