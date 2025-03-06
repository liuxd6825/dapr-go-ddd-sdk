package ddd_mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rsql/rsql_mongo"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *Dao[T]) Delete(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	return r.DeleteById(ctx, r.GetTenantId(entity), r.GetId(entity), opts...)
}

func (r *Dao[T]) DeleteByRSQL(ctx context.Context, tenantId, rSQL string, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResult[T]()
	filter, err := r.getFilter(tenantId, rSQL)
	if err != nil {
		return res.SetError(err)
	}
	db := r.deleteByFilter(ctx, tenantId, filter)
	return res.SetError(db.GetError()).SetRowsAffected(db.RowsAffected)
}

func (r *Dao[T]) DeleteById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	data := map[string]interface{}{
		ConstIdField:  id,
		TenantIdField: tenantId,
	}
	return r.DeleteByMap(ctx, tenantId, data)
}

func (r *Dao[T]) DeleteByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResult[T]()
	if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is empty")); err != nil {
		return res.SetError(err)
	}
	if len(ids) == 0 {
		return res
	}

	_ = r.DoSet(func() (T, error) {
		var null T
		filter := bson.D{}
		filter = append(filter, bson.E{Key: ConstIdField, Value: bson.M{"$in": ids}})
		filter = append(filter, bson.E{Key: ConstTenantIdField, Value: tenantId})
		deleteOptions := getDeleteOptions(opts...)
		sCtx := r.getSessionCtx(ctx)
		del, err := r.getCollection(ctx).DeleteMany(sCtx, filter, deleteOptions)
		res.SetError(err)
		res.SetRowsAffected(del.DeletedCount)
		return null, err
	})
	return res
}

func (r *Dao[T]) DeleteAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	data := map[string]interface{}{}
	return r.DeleteByMap(ctx, tenantId, data)
}
func (r *Dao[T]) DeleteByMap(ctx context.Context, tenantId string, filterMap map[string]any, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	return r.deleteByAny(ctx, tenantId, filterMap, opts...)
}

func (r *Dao[T]) deleteByAny(ctx context.Context, tenantId string, filterAny any, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	if err := assert.NotNil(filterAny, assert.NewOptions("filterMap is nil")); err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is empty")); err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	if filterMap, ok := filterAny.(map[string]any); ok {
		return r.deleteByMap(ctx, tenantId, filterMap, opts...)
	} else if filter, ok := filterAny.(*rsql_mongo.Filter); ok {
		return r.deleteByFilter(ctx, tenantId, filter, opts...)
	}
	return ddd_repository.NewSetResultError[T](errors.New("filterMap is not map[string]any or *rsql_mongo.Filter"))
}

func (r *Dao[T]) deleteByMap(ctx context.Context, tenantId string, filterMap map[string]any, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResultEmpty[T]()
	gp.Try(func() error {
		ctx := r.getSessionCtx(ctx)
		filter := r.NewFilter(tenantId, filterMap)
		deleteOptions := getDeleteOptions(opts...)
		mRes, err := r.getCollection(ctx).DeleteMany(ctx, filter, deleteOptions)
		if mRes != nil {
			res.SetRowsAffected(mRes.DeletedCount)
		}
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (r *Dao[T]) deleteByFilter(ctx context.Context, tenantId string, filter *rsql_mongo.Filter, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResultEmpty[T]()
	gp.Try(func() error {
		if !filter.IsAggregate() {
			res = r.deleteByMap(ctx, tenantId, filter.Match, opts...)
			return res.Error
		}

		ctx = r.getSessionCtx(ctx)
		pipeline := filter.NewPipeline()
		coll := r.getCollection(ctx)
		cursor, err := coll.Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}
		defer cursor.Close(ctx)

		// 提取符合条件的 _id
		var idsToDelete []interface{}
		for cursor.Next(context.TODO()) {
			var result bson.M
			if err := cursor.Decode(&result); err != nil {
				return err
			}
			idsToDelete = append(idsToDelete, result["id"])
		}

		if err := cursor.Err(); err != nil {
			return err
		}

		// 执行删除操作
		if len(idsToDelete) > 0 {
			// 使用 $in 操作符删除符合条件的文档
			mRes, err := coll.DeleteMany(context.TODO(), bson.D{
				{"_id", bson.D{{"$in", idsToDelete}}},
			})
			if mRes != nil {
				res.SetRowsAffected(mRes.DeletedCount)
			}
			return err
		}
		return nil
	}).Catch(func(e error) {
		res.SetError(e)
	})
	return res
}
