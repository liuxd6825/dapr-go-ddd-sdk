package store_mongodb

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	assert2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

func (r *Dao[T]) Update(ctx context.Context, entity T, opts ...store.Options) *store.SetResult[T] {
	return r.updateById(ctx, entity, opts...)
}

func (r *Dao[T]) updateById(ctx context.Context, entity T, opts ...store.Options) *store.SetResult[T] {
	res := store.NewSetResultEmpty[T]()
	gp.Try(func() error {
		opt := store.NewOptions(opts...)
		id := r.GetId(entity)
		tenantId := r.GetTenantId(entity)

		data := r.getUpdateData(ctx, tenantId, entity, opt)
		setData := bson.M{"$set": data}

		uOpt := getUpdateOptions(opts...)
		sCtx := r.getSessionCtx(ctx)

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
		updateOptions := getUpdateOptions(opts...)

		mData := r.getUpdateData(ctx, tenantId, data, opt)
		setData := bson.M{"$set": mData}

		sCtx := r.getSessionCtx(ctx)
		mRes, err := r.getCollection(ctx).UpdateMany(sCtx, filter.Match, setData, updateOptions)
		if mRes != nil {
			res.SetRowsAffected(mRes.ModifiedCount + mRes.UpsertedCount)
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
			setData := bson.M{"$set": item}
			id := r.GetId(entity)

			model := mongo.NewUpdateOneModel().SetFilter(
				bson.M{ConstTenantIdField: tenantId, ConstIdField: id},
			).SetUpdate(setData).SetUpsert(false)
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
	res := r.updateMap(ctx, tenantId, filter, m, opts...)
	return res
}

func (r *Dao[T]) UpdateMap(ctx context.Context, tenantId string, id string, data map[string]any, opts ...store.Options) *store.SetResult[T] {
	res := store.NewSetResultEmpty[T]()
	gp.Try(func() error {
		filter, err := r.getFilter(tenantId, fmt.Sprintf("%s=='%s'", ConstIdField, id))
		if err != nil {
			return err
		}
		res = r.updateMap(ctx, tenantId, filter.Match, data, opts...)
		return res.Error
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (r *Dao[T]) UpdateMapByRSQL(ctx context.Context, tenantId string, filterRSQL string, data map[string]any, opts ...store.Options) *store.SetResult[T] {
	res := store.NewSetResultEmpty[T]()
	gp.Try(func() error {
		filter, err := r.getFilter(tenantId, filterRSQL)
		if err != nil {
			return err
		}
		res = r.updateMap(ctx, tenantId, filter.Match, data, opts...)
		return res.Error
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (r *Dao[T]) updateMap(ctx context.Context, tenantId string, filter any, data map[string]any, opts ...store.Options) *store.SetResult[T] {
	res := store.NewSetResultEmpty[T]()
	gp.Try(func() error {
		if err := assert2.NotEmpty(tenantId, assert2.NewOptions("tenantId is empty")); err != nil {
			return err
		}

		if err := assert2.NotNil(filter, assert2.NewOptions("filterMap is nil")); err != nil {
			return err
		}
		updateOptions := getUpdateOptions(opts...)
		doc := r.getUpdateData(ctx, tenantId, data, opts...)

		setData := bson.M{"$set": doc}

		sCtx := r.getSessionCtx(ctx)
		upeRes, err := r.getCollection(ctx).UpdateMany(sCtx, filter, setData, updateOptions)
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

func (r *Dao[T]) getUpdateData(ctx context.Context, tenantId string, data any, opts ...store.Options) any {
	r.eb.SetUpdatedInfo(ctx, data)
	doc := r.entity2db(data)
	opt := store.NewOptions(opts...)

	for _, field := range r.schema.Fields {
		if !field.Updatable || field.PrimaryKey {
			delete(doc, field.DBName)
		}
	}

	updateFields := opt.GetUpdateFields()
	if len(updateFields) > 0 {
		for _, field := range r.schema.Fields {
			if IncludeField(field, updateFields) {
				continue
			}
			delete(doc, field.DBName)
		}
	}

	cancelFields := opt.GetUpdateCancel()
	if len(cancelFields) > 0 {
		for _, field := range r.schema.Fields {
			if IncludeField(field, cancelFields) {
				delete(doc, field.DBName)
			}
		}
	}
	return doc
}

func (r *Dao[T]) getInsertData(ctx context.Context, tenantId string, data any, opts ...store.Options) any {
	r.eb.SetCreatedInfo(ctx, data)
	doc := r.entity2db(data)
	for _, field := range r.schema.Fields {
		if !field.Creatable {
			delete(doc, field.DBName)
		}
	}
	return doc
}

func IncludeField(field *store.Field, fields []string) bool {
	if fields == nil || len(fields) == 0 {
		return false
	}
	for _, f := range fields {
		f = strings.ToLower(f)
		if f == strings.ToLower(field.DBName) || f == strings.ToLower(field.Name) {
			return true
		}
	}
	return false
}
