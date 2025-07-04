package store_mongodb

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	assert2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/mapperutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *Dao[T]) Update(ctx context.Context, entity T, opts ...store.Options) *store.SetResult[T] {
	return r.updateById(ctx, entity, opts...)
}

func (r *Dao[T]) updateById(ctx context.Context, entity T, opts ...store.Options) *store.SetResult[T] {
	res := store.NewSetResultEmpty[T]()
	gp.Try(func() error {
		r.eb.SetUpdatedInfo(ctx, entity)

		opt := store.NewOptions(opts...)
		data := r.getUpdateData(entity, opt)

		uOpt := getUpdateOptions(opts...)
		setData := bson.M{"$set": data}

		sCtx := r.getSessionCtx(ctx)
		id := r.GetId(entity)
		tenantId := r.GetTenantId(entity)

		filter := bson.M{
			ConstIdField:       id,
			ConstTenantIdField: tenantId,
		}
		mRes, err := r.getCollection(ctx).UpdateOne(sCtx, filter, setData, uOpt)
		if err != nil {
			return err
		}
		if mRes != nil {
			res.SetRowsAffected(mRes.ModifiedCount)
		}
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res

}

func (r *Dao[T]) UpdateByRSQL(ctx context.Context, tenantId, filterRSQL string, data T, opts ...store.Options) *store.SetResult[T] {
	res := store.NewSetResultEmpty[T]()
	gp.Try(func() error {
		filter, err := r.getFilter(tenantId, filterRSQL)
		if err != nil {
			return err
		}
		if filter.IsAggregate() {
			return errors.New("aggregate update is not supported")
		}
		opt := store.NewOptions(opts...)
		r.eb.SetUpdatedInfo(ctx, data)
		mData := r.getUpdateData(data, opt)
		setData := bson.M{"$set": mData}
		updateOptions := getUpdateOptions(opts...)
		sCtx := r.getSessionCtx(ctx)
		mRes, err := r.getCollection(ctx).UpdateMany(sCtx, filter.Match, setData, updateOptions)
		if mRes != nil {
			res.SetRowsAffected(mRes.ModifiedCount)
		}
		return err

	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (r *Dao[T]) UpdateMany(ctx context.Context, tenantId string, entities []T, opts ...store.Options) *store.SetResult[T] {
	res := store.NewSetResultEmpty[T]()
	gp.Try(func() error {
		if entities == nil || len(entities) == 0 {
			return errors.New("entities is nil")
		}
		var list []mongo.WriteModel
		for _, entity := range entities {
			r.eb.SetUpdatedInfo(ctx, entity)
			item := r.entity2db(entity)
			data := bson.M{"$set": item}
			id := r.GetId(entity)

			model := mongo.NewUpdateOneModel().SetFilter(
				bson.M{ConstTenantIdField: tenantId, ConstIdField: id},
			).SetUpdate(data).SetUpsert(false)
			list = append(list, model)
		}

		mRes, err := r.BulkWrite(ctx, list)
		if err != nil {
			return err
		}
		if mRes != nil {
			res.SetRowsAffected(mRes.ModifiedCount + mRes.UpsertedCount + mRes.DeletedCount + mRes.InsertedCount)
		}
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res

}

func (r *Dao[T]) UpdateMapById(ctx context.Context, tenantId string, id string, data map[string]any, opts ...store.Options) *store.SetResult[T] {
	filter := bson.M{ConstTenantIdField: tenantId, ConstIdField: id}
	m := r.getDbMap(data)
	r.eb.SetUpdatedInfo(ctx, data)
	res := r.UpdateMapAndGetCount(ctx, tenantId, filter, m, opts...)
	return res
}

func (r *Dao[T]) UpdateMap(ctx context.Context, tenantId string, id string, data map[string]any, opts ...store.Options) *store.SetResult[T] {
	res := store.NewSetResultEmpty[T]()
	gp.Try(func() error {
		filter, err := r.getFilter(tenantId, fmt.Sprintf("%s=='%s'", ConstIdField, id))
		if err != nil {
			return err
		}
		res = r.UpdateMapAndGetCount(ctx, tenantId, filter.Match, data, opts...)
		return res.Error
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (r *Dao[T]) UpdateMapAndGetCount(ctx context.Context, tenantId string, filter any, data any, opts ...store.Options) *store.SetResult[T] {
	res := store.NewSetResultEmpty[T]()
	gp.Try(func() error {
		if err := assert2.NotEmpty(tenantId, assert2.NewOptions("tenantId is empty")); err != nil {
			return err
		}

		if err := assert2.NotNil(filter, assert2.NewOptions("filterMap is nil")); err != nil {
			return err
		}
		updateOptions := getUpdateOptions(opts...)
		var f any
		if v, ok := filter.(map[string]any); ok {
			f = r.NewFilter(tenantId, v)
		} else {
			f = filter
		}
		sCtx := r.getSessionCtx(ctx)
		r.eb.SetUpdatedInfo(ctx, data)
		upeRes, err := r.getCollection(ctx).UpdateMany(sCtx, f, data, updateOptions)
		res.SetError(err)
		if upeRes != nil {
			res.SetRowsAffected(upeRes.ModifiedCount)
		}
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (r *Dao[T]) getUpdateData(data any, opts ...store.Options) any {
	if opts == nil {
		return data
	}
	opt := store.NewOptions(opts...)
	updateCancel := opt.GetUpdateCancel()
	updateFields := opt.GetUpdateFields()

	if (updateCancel == nil || len(updateCancel) == 0) && (updateFields == nil || len(updateFields) == 0) {
		if m, ok := data.(map[string]any); ok {
			return r.getDbMap(m)
		} else {
			doc := r.entity2db(data)
			return doc
		}

	}
	m := make(map[string]any)

	maskType := mapperutils.MaskTypeContain
	mask := updateFields
	mask = append(mask, "UpdatedTime", "UpdaterId", "UpdaterName")
	if updateCancel != nil {
		maskType = mapperutils.MaskTypeExclude
		mask = updateCancel
	}
	if err := mapperutils.MaskMapperType(data, &m, mask, maskType); err != nil {
		return store.NewSetResultEmpty[T]().SetError(err)
	}
	m = r.getDbMap(m)
	return m
}
