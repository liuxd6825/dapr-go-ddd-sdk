package store_mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	assert2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *Dao[T]) InsertOrUpdate(ctx context.Context, entity T, opts ...store.Options) *store.SetResult[T] {
	res := store.NewSetResultEmpty[T]()
	gp.Try(func() error {
		if err := assert2.NotEmpty(r.GetTenantId(entity), assert2.NewOptions("tenantId is empty")); err != nil {
			return err
		}
		filter := r.NewFilter(r.GetTenantId(entity), map[string]interface{}{"id": r.GetId(entity)})
		findOneOptions := getFindOneOptions(opts...)
		isFound := true
		sCtx := r.getSessionCtx(ctx)
		if err := r.getCollection(ctx).FindOne(sCtx, filter, findOneOptions).Err(); err != nil {
			if err == mongo.ErrNoDocuments {
				isFound = false
			} else {
				return err
			}
		}
		// 是否找到数据
		if isFound {
			res = r.updateById(sCtx, entity, opts...)
			return res.Error
		} else {
			res = r.Insert(ctx, entity, opts...)
			return res.Error
		}
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res

}

func (r *Dao[T]) Insert(ctx context.Context, entity T, opts ...store.Options) *store.SetResult[T] {
	res := store.NewSetResultEmpty[T]()
	gp.Try(func() error {
		tenantId := r.GetTenantId(entity)
		if err := assert2.NotEmpty(tenantId, assert2.NewOptions("tenantId is empty")); err != nil {
			return err
		}
		doc := r.getInsertData(ctx, tenantId, entity)
		ctx := r.getSessionCtx(ctx)
		mRes, err := r.getCollection(ctx).InsertOne(ctx, doc, getInsertOneOptions(opts...))
		if err != nil {
			return err
		}
		if mRes != nil && mRes.InsertedID != nil {
			res.SetRowsAffected(1)
		}
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

// InsertMap
// @Description: 插入数据
// @receiver r
// @param ctx
// @param tenantId
// @param data
// @param opts
// @return error
func (r *Dao[T]) InsertMap(ctx context.Context, tenantId string, data map[string]interface{}, opts ...store.Options) (res *store.SetResult[T]) {
	res = store.NewSetResultEmpty[T]()
	gp.Try(func() error {
		if err := assert2.NotEmpty(tenantId, assert2.NewOptions("tenantId is empty")); err != nil {
			return err
		}
		doc := r.getInsertData(ctx, tenantId, data)
		ctx := r.getSessionCtx(ctx)
		inRes, err := r.getCollection(ctx).InsertOne(ctx, doc, getInsertOneOptions(opts...))
		if inRes != nil && inRes.InsertedID != nil {
			res.SetRowsAffected(1)
		}
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (r *Dao[T]) InsertMany(ctx context.Context, tenantId string, entities []T, opts ...store.Options) *store.SetResult[T] {
	res := store.NewSetResultEmpty[T]()
	gp.Try(func() error {
		ctx := r.getSessionCtx(ctx)
		if entities == nil || len(entities) == 0 {
			return errors.New("entities is nil")
		}
		var docs []any
		for _, e := range entities {
			r.eb.SetTenantId(e, tenantId)
			doc := r.getInsertData(ctx, tenantId, e)
			docs = append(docs, doc)
		}

		mRes, err := r.getCollection(ctx).InsertMany(ctx, docs, getInsertManyOptions(opts...))
		if err == nil && mRes != nil {
			count := int64(len(mRes.InsertedIDs))
			res.SetRowsAffected(count)
		}
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}
