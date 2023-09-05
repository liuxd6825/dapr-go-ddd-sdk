package ddd_mongodb

import (
	"context"
	"fmt"
	"github.com/liuxd6825/components-contrib/liuxd/common/utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

const (
	ConstIdField       = "_id"
	ConstTenantIdField = "tenant_id"
)

type Dao[T ddd.Entity] struct {
	entityBuilder *ddd_repository.EntityBuilder[T]
	collection    *mongo.Collection
	mongodb       *MongoDB
	null          T
	newFun        func() T                                                // 新建实体结构方法
	initfu        func() (mongodb *MongoDB, collection *mongo.Collection) // 初始化
}

func NewDao[T ddd.Entity](initfu func() (mongodb *MongoDB, collection *mongo.Collection)) *Dao[T] {
	r := &Dao[T]{}
	r.initfu = initfu
	return r
}

func (r *Dao[T]) Init(mongodb *MongoDB, collection *mongo.Collection) {
	r.collection = collection
	r.mongodb = mongodb
}

func (r *Dao[T]) NewEntity() (T, error) {
	return reflectutils.NewStruct[T]()
}

func (r *Dao[T]) NewEntityList() ([]T, error) {
	return reflectutils.NewSlice[[]T]()
}

func (r *Dao[T]) getCollection() *mongo.Collection {
	if r.collection != nil {
		return r.collection
	}
	r.Init(r.initfu())
	return r.collection
}

func (r *Dao[T]) Save(ctx context.Context, data *ddd.SetData[T], opts ...ddd_repository.Options) (setResult *ddd_repository.SetResult[T]) {
	defer func() {
		if e := recover(); e != nil {
			if err := errors.GetRecoverError(e); err != nil {
				setResult = ddd_repository.NewSetResultError[T](err)
			}
		}
	}()
	for _, item := range data.Items() {
		statue := item.Statue()
		entity := item.Data().(T)
		var err error
		switch statue {
		case ddd.DataStatueCreate:
			err = r.Insert(ctx, entity, opts...).GetError()
		case ddd.DataStatueUpdate:
			err = r.Update(ctx, entity, opts...).GetError()
		case ddd.DataStatueDelete:
			err = r.DeleteById(ctx, entity.GetTenantId(), entity.GetId(), opts...).GetError()
		case ddd.DataStatueCreateOrUpdate:
			err = r.InsertOrUpdate(ctx, entity, opts...).GetError()
		}
		if err != nil {
			return ddd_repository.NewSetResultError[T](err)
		}
	}
	return ddd_repository.NewSetResultError[T](nil)
}

func (r *Dao[T]) InsertOrUpdate(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	if err := assert.NotEmpty(entity.GetTenantId(), assert.NewOptions("tenantId is empty")); err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	return r.DoSet(func() (T, error) {
		filter := r.NewFilter(entity.GetTenantId(), map[string]interface{}{"id": entity.GetId()})
		findOneOptions := getFindOneOptions(opts...)
		isFound := true

		if err := r.getCollection().FindOne(ctx, filter, findOneOptions).Err(); err != nil {
			if err == mongo.ErrNoDocuments {
				// 没有找到，设置 isFound 状态
				isFound = false
			} else {
				var null T
				return null, err
			}
		}

		// 是否找到数据
		if isFound {
			return r.updateById(ctx, entity, opts...)
		} else {
			_, err := r.getCollection().InsertOne(ctx, entity, getInsertOneOptions(opts...))
			return entity, err
		}
	})
	return ddd_repository.NewSetResultError[T](nil)
}

func (r *Dao[T]) Insert(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	if err := assert.NotEmpty(entity.GetTenantId(), assert.NewOptions("tenantId is empty")); err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	return r.DoSet(func() (T, error) {
		_, err := r.getCollection().InsertOne(ctx, entity, getInsertOneOptions(opts...))
		return entity, err
	})
}

func (r *Dao[T]) InsertMap(ctx context.Context, tenantId string, data map[string]interface{}, opts ...ddd_repository.Options) error {
	if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is empty")); err != nil {
		return err
	}
	data["tenant_id"] = tenantId
	_, err := r.getCollection().InsertOne(ctx, data, getInsertOneOptions(opts...))
	return err
}

func (r *Dao[T]) InsertMany(ctx context.Context, entitits []T, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T] {
	if entitits == nil || len(entitits) == 0 {
		return ddd_repository.NewSetManyResultError[T](errors.New("entitits is nil"))
	}

	for _, e := range entitits {
		if err := assert.NotEmpty(e.GetTenantId(), assert.NewOptions("tenantId is empty")); err != nil {
			return ddd_repository.NewSetManyResultError[T](err)
		}
	}

	var docs []interface{}
	for _, e := range entitits {
		docs = append(docs, e)
	}

	return r.DoSetMany(func() ([]T, error) {
		_, err := r.getCollection().InsertMany(ctx, docs, getInsertManyOptions(opts...))
		return entitits, err
	})
}

func (r *Dao[T]) Update(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	if err := assert.NotEmpty(entity.GetTenantId(), assert.NewOptions("tenantId is empty")); err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	return r.DoSet(func() (T, error) {
		return r.updateById(ctx, entity, opts...)
	})
}

func (r *Dao[T]) updateById(ctx context.Context, entity T, opts ...ddd_repository.Options) (T, error) {
	opt := ddd_repository.NewOptions(opts...)
	data := r.getUpdateData(entity, opt)
	uopt := getUpdateOptions(opts...)
	setData := bson.M{"$set": data}
	_, err := r.getCollection().UpdateByID(ctx, entity.GetId(), setData, uopt)
	if err != nil {
		return entity, err
	}
	return entity, err
}

func (r *Dao[T]) UpdateManyByFilter(ctx context.Context, tenantId, filter string, data any, opts ...ddd_repository.Options) *ddd_repository.SetManyCountResult {
	filterMap, err := r.getFilterMap(tenantId, filter)
	if err != nil {
		return ddd_repository.NewSetManyCountResultError(err)
	}
	opt := ddd_repository.NewOptions(opts...)
	return r.DoSetManyCount(func() (*mongo.UpdateResult, error) {
		mdata := r.getUpdateData(data, opt)
		setData := bson.M{"$set": mdata}
		updateOptions := getUpdateOptions(opts...)
		res, err := r.getCollection().UpdateMany(ctx, filterMap, setData, updateOptions)
		return res, err
	})
}

func (r *Dao[T]) UpdateManyById(ctx context.Context, entities []T, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T] {
	if entities == nil || len(entities) == 0 {
		return ddd_repository.NewSetManyResult[T](entities, nil)
	}

	var list []mongo.WriteModel
	for _, entity := range entities {
		data := bson.M{"$set": entity}
		model := mongo.NewUpdateOneModel().SetFilter(bson.D{{"_id", entity.GetId()}}).SetUpdate(data).SetUpsert(true)
		list = append(list, model)
	}
	_, err := r.BulkWrite(ctx, list)
	if err != nil {
		return ddd_repository.NewSetManyResultError[T](err)
	}
	return ddd_repository.NewSetManyResult[T](entities, err)
}

func (r *Dao[T]) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...ddd_repository.Options) (*ddd_repository.BulkWriteResult, error) {
	if len(models) == 0 {
		return &ddd_repository.BulkWriteResult{}, nil
	}
	opt := &options.BulkWriteOptions{}
	bulkRes, err := r.getCollection().BulkWrite(ctx, models, opt)
	if err != nil {
		return nil, err
	}

	res := r.newBulkWriteResult(bulkRes)
	return res, err
}

func (r *Dao[T]) newBulkWriteResult(bulkRes *mongo.BulkWriteResult) *ddd_repository.BulkWriteResult {
	res := ddd_repository.NewBulkWriteResult()
	if bulkRes != nil {
		res.InsertedCount = bulkRes.InsertedCount
		res.MatchedCount = bulkRes.MatchedCount
		res.ModifiedCount = bulkRes.ModifiedCount
		res.DeletedCount = bulkRes.DeletedCount
		res.UpsertedCount = bulkRes.UpsertedCount
		res.UpsertedIDs = bulkRes.UpsertedIDs
		res.SetEmpty(false)
	}
	if res.UpsertedIDs == nil {
		res.UpsertedIDs = map[int64]interface{}{}
	}
	return res
}

func (r *Dao[T]) UpdateManyMaskById(ctx context.Context, entities []T, mask []string, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T] {
	if entities == nil || len(entities) == 0 {
		return ddd_repository.NewSetManyResultError[T](errors.New("entities is nil"))
	}

	for _, e := range entities {
		if err := assert.NotEmpty(e.GetTenantId(), assert.NewOptions("tenantId is empty")); err != nil {
			return ddd_repository.NewSetManyResultError[T](err)
		}
	}

	var docs []interface{}
	for _, e := range entities {
		if len(mask) == 0 {
			docs = append(docs, e)
		} else {
			m := make(map[string]interface{})
			if err := types.MaskMapper(e, &m, mask); err != nil {
				return ddd_repository.NewSetManyResultError[T](err)
			}
			m[ConstIdField] = e.GetId()
			doc := asDocument(m)
			docs = append(docs, doc)
		}
	}

	return r.DoSetMany(func() ([]T, error) {
		for _, doc := range docs {
			id, err := getDocumentId(doc)
			if err != nil {
				return nil, err
			}
			updateOptions := getUpdateOptions(opts...)
			filter := bson.D{{ConstIdField, id}}

			setData := bson.M{"$set": doc}
			_, err = r.getCollection().UpdateOne(ctx, filter, setData, updateOptions)
			if err != nil {
				return nil, err
			}
		}
		return entities, nil
	})
}

func (r *Dao[T]) getUpdateData(data any, opts ...ddd_repository.Options) any {
	if opts == nil {
		return data
	}
	opt := ddd_repository.NewOptions(opts...)
	updateCancel := opt.GetUpdateCancel()
	updateFields := opt.GetUpdateFields()

	if (updateCancel == nil || len(updateCancel) == 0) && (updateFields == nil || len(updateFields) == 0) {
		if m, ok := data.(map[string]any); ok {
			return r.getMap(m)
		} else {
			return data
		}

	}
	m := make(map[string]any)

	maskType := types.MaskTypeContain
	mask := updateFields
	mask = append(mask, "UpdatedTime", "UpdaterId", "UpdaterName")
	if updateCancel != nil {
		maskType = types.MaskTypeExclude
		mask = updateCancel
	}
	if err := types.MaskMapperType(data, &m, mask, maskType); err != nil {
		return ddd_repository.NewSetManyResultError[T](err)
	}
	m = r.getMap(m)
	return m
}

func (r *Dao[T]) getMap(m map[string]any) map[string]any {
	m = maputils.MapToSnakeKey(m)
	if idVal, ok := m["id"]; ok {
		m["_id"] = idVal
		delete(m, "id")
	}
	return m
}

func (r *Dao[T]) UpdateMapById(ctx context.Context, tenantId string, id string, data map[string]any, opts ...ddd_repository.Options) error {
	filter := bson.M{"tenant_id": tenantId, "_id": id}
	m := r.getMap(data)
	_, err := r.UpdateMapAndGetCount(ctx, tenantId, filter, m, opts...)
	return err
}

func (r *Dao[T]) FindOneAndUpdateById(ctx context.Context, tenantId string, id string, data map[string]any, opts ...ddd_repository.Options) (T, error) {
	var null T
	if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is empty")); err != nil {
		return null, err
	}

	if err := assert.NotEmpty(id, assert.NewOptions("id is empty")); err != nil {
		return null, err
	}
	filter := bson.M{"tenant_id": tenantId, "_id": id}
	udpate := r.getMap(data)
	_, err := r.getCollection().UpdateOne(ctx, filter, udpate)
	if err != nil {
		return null, err
	}
	find := r.FindById(ctx, tenantId, id, opts...)
	return find.GetData(), find.GetError()
}

func (r *Dao[T]) UpdateMap(ctx context.Context, tenantId string, filter any, data any, opts ...ddd_repository.Options) error {
	_, err := r.UpdateMapAndGetCount(ctx, tenantId, filter, data, opts...)
	return err
}

func (r *Dao[T]) UpdateMapAndGetCount(ctx context.Context, tenantId string, filter any, data any, opts ...ddd_repository.Options) (int64, error) {
	if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is empty")); err != nil {
		return 0, err
	}

	if err := assert.NotNil(filter, assert.NewOptions("filterMap is nil")); err != nil {
		return 0, err
	}

	updateOptions := getUpdateOptions(opts...)
	var f any
	if v, ok := filter.(map[string]any); ok {
		f = r.NewFilter(tenantId, v)
	} else {
		f = filter
	}
	res, err := r.getCollection().UpdateOne(ctx, f, data, updateOptions)
	if err != nil {
		return 0, err
	}
	return res.UpsertedCount, nil
}

func (r *Dao[T]) Delete(ctx context.Context, entity ddd.Entity, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	return r.DeleteById(ctx, entity.GetTenantId(), entity.GetId(), opts...)
}

func (r *Dao[T]) DeleteByFilter(ctx context.Context, tenantId, filter string, opts ...ddd_repository.Options) error {
	if filterMap, err := r.getFilterMap(tenantId, filter); err != nil {
		return err
	} else if err := r.DeleteByMap(ctx, tenantId, filterMap).GetError(); err != nil {
		return err
	}
	return nil
}

func (r *Dao[T]) DeleteById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	data := map[string]interface{}{
		ConstIdField:  id,
		TenantIdField: tenantId,
	}
	return r.DeleteByMap(ctx, tenantId, data)
}

func (r *Dao[T]) DeleteByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) error {
	if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is empty")); err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}

	setResult := r.DoSet(func() (T, error) {
		var null T
		filter := bson.D{}
		filter = append(filter, bson.E{Key: ConstIdField, Value: bson.M{"$in": ids}})
		filter = append(filter, bson.E{Key: ConstTenantIdField, Value: tenantId})
		deleteOptions := getDeleteOptions(opts...)
		_, err := r.getCollection().DeleteMany(ctx, filter, deleteOptions)
		return null, err
	})
	return setResult.GetError()
}

func (r *Dao[T]) DeleteAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	data := map[string]interface{}{}
	return r.DeleteByMap(ctx, tenantId, data)
}

func (r *Dao[T]) DeleteByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	if err := assert.NotNil(filterMap, assert.NewOptions("filterMap is nil")); err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is empty")); err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	return r.DoSet(func() (T, error) {
		filter := r.NewFilter(tenantId, filterMap)
		deleteOptions := getDeleteOptions(opts...)
		_, err := r.getCollection().DeleteMany(ctx, filter, deleteOptions)
		var result T
		return result, err
	})
}

func (r *Dao[T]) NewFilter(tenantId string, filterMap map[string]interface{}) bson.M {
	/*	conditions := bson.M{"name": bson.M{"$regex": "me"},
		"$or": []bson.M{
			bson.M{"repair": bson.M{"$eq": "ac"}},
		},
		"$and": []bson.M{
			bson.M{"repair": bson.M{"$eq": "tv"}},
			bson.M{"phone": bson.M{"$gte": 1091, "$lte": 1100}},
		}}
	*/

	if filterMap == nil || len(filterMap) == 0 {
		return bson.M{ConstTenantIdField: bson.M{"$eq": tenantId}}
	}

	ands := make([]bson.M, 1)
	ands[0] = bson.M{ConstTenantIdField: bson.M{"$eq": tenantId}}

	//filter := bson.D{{"name", "Bagels N Buns"}}
	for fieldName, fieldValue := range filterMap {
		if fieldName != ConstIdField {
			fieldName = AsFieldName(fieldName)
		}
		item := bson.M{fieldName: bson.M{"$eq": fieldValue}}
		ands = append(ands, item)
	}
	filter := bson.M{
		"$and": ands,
	}
	return filter
}
func (r *Dao[T]) FindById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) *ddd_repository.FindOneResult[T] {
	idMap := map[string]interface{}{
		ConstIdField: id,
	}
	return r.FindOneByMap(ctx, tenantId, idMap, opts...)
}

func (r *Dao[T]) FindByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	filter := bson.D{}
	filter = append(filter, bson.E{Key: ConstIdField, Value: bson.M{"$in": ids}})
	return r.FindListByMap(ctx, tenantId, filter.Map(), opts...)
}

func (r *Dao[T]) FindOneByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...ddd_repository.Options) *ddd_repository.FindOneResult[T] {
	var null T
	return r.DoFindOne(func() (T, bool, error) {
		filter := r.NewFilter(tenantId, filterMap)
		findOneOptions := getFindOneOptions(opts...)
		data, err := r.NewEntity()
		if err != nil {
			return null, false, err
		}
		result := r.getCollection().FindOne(ctx, filter, findOneOptions)
		if result.Err() != nil {
			return null, false, result.Err()
		}
		if err := result.Decode(data); err != nil {
			return null, false, err
		}
		return data, true, nil
	})
}

func (r *Dao[T]) FindListByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	return r.DoFindList(func() ([]T, bool, error) {
		var list []T
		filter := r.NewFilter(tenantId, filterMap)
		findOptions := getFindOptions(opts...)
		cursor, err := r.getCollection().Find(ctx, filter, findOptions)
		if err != nil {
			return nil, false, err
		}
		err = cursor.All(ctx, &list)
		return list, len(list) > 0, err
	})
}

func (r *Dao[T]) FindListByBsonM(ctx context.Context, tenantId string, filter bson.M, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	return r.DoFindList(func() ([]T, bool, error) {
		var list []T
		findOptions := getFindOptions(opts...)
		cursor, err := r.getCollection().Find(ctx, filter, findOptions)
		if err != nil {
			return nil, false, err
		}
		err = cursor.All(ctx, &list)
		return list, len(list) > 0, err
	})
}

func (r *Dao[T]) FindByRSQL(ctx context.Context, tenantId string, rsql string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	return r.DoList(tenantId, rsql, func(filterMap map[string]interface{}) ([]T, bool, error) {
		var list []T
		findOptions := getFindOptions(opts...)
		cursor, err := r.getCollection().Find(ctx, filterMap, findOptions)
		if err != nil {
			return nil, false, err
		}
		err = cursor.All(ctx, &list)
		return list, len(list) > 0, err
	})
}

func (r *Dao[T]) FindAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	return r.FindListByMap(ctx, tenantId, nil, opts...)
}

func (r *Dao[T]) findPaging(ctx context.Context, query ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	return r.DoFilter(query.GetTenantId(), query.GetFilter(), func(filter map[string]interface{}) (*ddd_repository.FindPagingResult[T], bool, error) {
		if err := assert.NotEmpty(query.GetTenantId(), assert.NewOptions("tenantId is empty")); err != nil {
			return nil, false, err
		}

		data, err := r.NewEntityList()
		if err != nil {
			return nil, false, err
		}

		findOptions := getFindOptions(opts...)
		if query.GetPageSize() > 0 {
			findOptions.SetLimit(query.GetPageSize())
			findOptions.SetSkip(query.GetPageSize() * query.GetPageNum())
		}
		if len(query.GetSort()) > 0 {
			sort, err := r.getSort(query.GetSort())
			if err != nil {
				return nil, false, err
			}
			findOptions.SetSort(sort)
		}

		cursor, err := r.getCollection().Find(ctx, filter, findOptions)
		if err != nil {
			return nil, false, err
		}

		err = cursor.All(ctx, &data)
		var totalRows *int64
		if query.GetIsTotalRows() {
			total, err := r.getCollection().CountDocuments(ctx, filter)
			if err != nil {
				return nil, false, err
			}
			totalRows = &total
		}

		findData := ddd_repository.NewFindPagingResult[T](data, totalRows, query, err)
		return findData, findData.IsFound, err
	})
}

/*
func (r Dao[T]) FindPaging2(ctx context.Context, query ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	var err error
	findOptions := getFindOptions(opts...)
	queryGroup := NewQueryGroup(query)

	g, err := queryGroup.GetGroup()
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}

	gpFilter, err := queryGroup.GetGroupPagingBsonFilter()
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}

	gFilter, err := queryGroup.GetGroupExpandGroupNoPagingBsonFilter()
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}

	gFilter1, err := queryGroup.GetGroupNoPagingFilter()
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}

	filter, err := queryGroup.GetFilter()
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}

	gSort, err := queryGroup.GetBsonFilterSort()
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}

	sort, err := queryGroup.GetFilterSort()
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}

	data, err := r.NewEntityList()
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}

	coll := r.getCollection()
	var findData *ddd_repository.FindPagingResult[T]
	var cur *mongo.Cursor
	///var curt *mongo.Cursor
	var errt error
	var totalRows int64

	isGroup := queryGroup.IsGroup()
	isPaging := queryGroup.IsPaging()
	isLeaf := queryGroup.IsLeaf()

	if isGroup {
		if isPaging {
			pipeline := mongo.Pipeline{}
			if gpFilter != nil && len(gpFilter) > 0 {
				pipeline = append(pipeline, gpFilter)
			}
			if g != nil && len(g) > 0 {
				pipeline = append(pipeline, g)
			}
			if sort != nil && len(sort) > 0 {
				pipeline = append(pipeline, sort)
			}

			skip := query.GetPageSize() * query.GetPageNum()
			pipeline = append(pipeline, bson.D{{"$skip", skip}})

			limit := query.GetPageSize()
			pipeline = append(pipeline, bson.D{{"$limit", limit}})

			cur, err = coll.Aggregate(ctx, pipeline)
			if err == nil {
				totalRows = int64(cur.RemainingBatchLength())
			}
		} else if !isLeaf {
			pipeline := mongo.Pipeline{}
			if gFilter != nil && len(gFilter) > 0 {
				pipeline = append(pipeline, gFilter)
			}
			if g != nil && len(g) > 0 {
				pipeline = append(pipeline, g)
			}
			if gSort != nil && len(gSort) > 0 {
				pipeline = append(pipeline, gSort)
			}
			if cur, err = coll.Aggregate(ctx, pipeline); err == nil {
				totalRows = int64(cur.RemainingBatchLength())
			}
		} else if isLeaf {
			findOptions.SetSort(sort)
			cur, err = coll.Find(ctx, gFilter1, findOptions)
			if err == nil {
				totalRows, errt = coll.CountDocuments(ctx, gFilter1)
			}
		}
	} else if !isGroup {
		findOptions.SetSort(sort)
		if query.GetPageSize() > 0 {
			findOptions.SetLimit(query.GetPageSize())
			findOptions.SetSkip(query.GetPageSize() * query.GetPageNum())
		}
		cur, err = coll.Find(ctx, filter, findOptions)
		if query.GetIsTotalRows() {
			totalRows, errt = coll.CountDocuments(ctx, filter)
		}
	}
	if err != nil || errt != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err, errt)
	}

	err = cur.All(ctx, &data)
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}

	findData = ddd_repository.NewFindPagingResult[T](data, &totalRows, query, err)
	// 进行汇总计算
	if len(query.GetValueCols()) > 0 {
		sumData, _, err := r.Sum(ctx, query, opts...)
		findData.SetSum(true, sumData, err)
	}
	return findData
}
*/
func (r Dao[T]) FindPaging(ctx context.Context, query ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	var err error
	findOptions := getFindOptions(opts...)
	queryGroup := NewQueryGroup(query)

	g, err := queryGroup.GetGroup()
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}

	totalGroup, err := queryGroup.GetTotalGroup()
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}

	//gpFilter, err := queryGroup.GetGroupPagingBsonFilter()
	//if err != nil {
	//	return ddd_repository.NewFindPagingResultWithError[T](err)
	//}
	//
	//gFilter, err := queryGroup.GetGroupExpandGroupNoPagingBsonFilter()
	//if err != nil {
	//	return ddd_repository.NewFindPagingResultWithError[T](err)
	//}
	//
	//gFilter1, err := queryGroup.GetGroupNoPagingFilter()
	//if err != nil {
	//	return ddd_repository.NewFindPagingResultWithError[T](err)
	//}

	filter, err := queryGroup.GetFilter()
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}

	filter1, err := queryGroup.GetGroupExpandFilter()
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}

	gSort, err := queryGroup.GetBsonFilterSort()
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}

	sort, err := queryGroup.GetFilterSort()
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}

	data, err := r.NewEntityList()
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}

	coll := r.getCollection()
	var findData *ddd_repository.FindPagingResult[T]
	var cur *mongo.Cursor
	///var curt *mongo.Cursor
	var errt error
	var totalRows int64

	isGroup := queryGroup.IsGroup()
	//isPaging := queryGroup.IsPaging()
	isLeaf := queryGroup.IsLeaf()

	if isGroup {
		if !isLeaf {
			pipeline := mongo.Pipeline{}
			if filter1 != nil && len(filter1) > 0 {
				pipeline = append(pipeline, bson.D{{"$match", filter1}})
			}
			if g != nil && len(g) > 0 {
				pipeline = append(pipeline, g)
			}
			if gSort != nil && len(gSort) > 0 {
				pipeline = append(pipeline, gSort)
			}
			pipeline = append(pipeline, totalGroup)

			pipeline = append(pipeline, bson.D{{
				"$project", map[string]interface{}{
					"_id":        "$_id",
					"data":       map[string]interface{}{"$slice": []interface{}{"$data", query.GetPageSize() * query.GetPageNum(), query.GetPageSize()}},
					"total_rows": "$total_rows",
				},
			}})

			//skip := query.GetPageSize() * query.GetPageNum()
			//pipeline = append(pipeline, bson.D{{"$skip", skip}})
			//
			//limit := query.GetPageSize()
			//pipeline = append(pipeline, bson.D{{"$limit", limit}})

			cur, err = coll.Aggregate(ctx, pipeline)
			//if err == nil {
			//	totalRows = int64(cur.RemainingBatchLength())
			//}
		}
	}
	if !isGroup || isLeaf {
		f := make(map[string]interface{})
		if isLeaf {
			f = filter1
		} else {
			f = filter
		}
		findOptions.SetSort(sort)
		if query.GetPageSize() > 0 {
			findOptions.SetLimit(query.GetPageSize())
			findOptions.SetSkip(query.GetPageSize() * query.GetPageNum())
		}
		cur, err = coll.Find(ctx, f, findOptions)
		if query.GetIsTotalRows() {
			totalRows, errt = coll.CountDocuments(ctx, f)
		}
	}
	if err != nil || errt != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err, errt)
	}

	if isGroup && !isLeaf {
		d := make([]struct {
			Data      []T   `json:"data" bson:"data"`
			TotalRows int64 `json:"totalRows" bson:"total_rows"`
		}, 0)
		err = cur.All(ctx, &d)
		if err != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err)
		}
		if d != nil && len(d) == 1 {
			data = d[0].Data
			totalRows = d[0].TotalRows
		}
	} else {
		err = cur.All(ctx, &data)
		if err != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err)
		}
	}
	if data == nil {
		data = []T{}
	}
	findData = ddd_repository.NewFindPagingResult[T](data, &totalRows, query, err)
	// 进行汇总计算
	if len(query.GetValueCols()) > 0 {
		sumData, _, err := r.Sum(ctx, query, opts...)
		findData.SetSum(true, sumData, err)
	} else {
		sumData := []T{}
		findData.SetSum(false, sumData, err)
	}
	return findData
}

func (r *Dao[T]) AggregateByPipeline(ctx context.Context, pipeline mongo.Pipeline, data interface{}, opts ...ddd_repository.Options) error {
	options := getAggregateOptions(opts...)
	cur, err := r.getCollection().Aggregate(ctx, pipeline, options)
	if err != nil {
		return err
	}
	err = cur.All(ctx, data)
	return err
}

func (r *Dao[T]) CopyTo(ctx context.Context, tenantId string, rsql string, toCollectionName string, opts ...ddd_repository.Options) error {
	options := getAggregateOptions(opts...)
	//db.record.aggregate([{$match:{opp_bank_name:"工商银行"}},{$out:"record1"}])
	filterMap, err := r.getFilterMap(tenantId, rsql)
	if err != nil {
		return err
	}
	pipeline := mongo.Pipeline{}
	if filterMap != nil && len(filterMap) > 0 {
		pipeline = append(pipeline, bson.D{{"$match", filterMap}})
	}
	pipeline = append(pipeline, bson.D{{"$out", toCollectionName}})
	_, err = r.getCollection().Aggregate(ctx, pipeline, options)
	return err
}

func (r Dao[T]) Sum(ctx context.Context, query ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) ([]T, bool, error) {
	if len(query.GetValueCols()) == 0 {
		return nil, false, nil
	}

	p := NewMongoProcess()
	filter := query.GetFilter()
	mustFilter := query.GetMustFilter()
	if len(filter) > 0 && len(mustFilter) > 0 {
		filter = fmt.Sprintf("(%s) and (%s)", filter, mustFilter)
	} else if len(mustFilter) > 0 {
		filter = mustFilter
	}
	if err := rsql.ParseProcess(filter, p); err != nil {
		return nil, false, err
	}
	filterMap, err := p.GetFilter(query.GetTenantId())
	if err != nil {
		return nil, false, err
	}

	return r.sum(ctx, filterMap, query.GetValueCols(), opts...)
}

func (r Dao[T]) sum(ctx context.Context, filterMap map[string]any, valueCols []*ddd_repository.ValueCol, opts ...ddd_repository.Options) ([]T, bool, error) {
	coll := r.getCollection()
	data, err := r.NewEntityList()
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err).DataResult()
	}
	var cur *mongo.Cursor
	summaryMap := make(map[string]interface{})
	summaryMap["_id"] = "total"
	for _, col := range valueCols {
		field := utils.SnakeString(col.Field)
		aggFunc := utils.SnakeString(col.AggFunc.Name())
		summaryMap[field] = map[string]interface{}{"$" + aggFunc: "$" + field}
	}

	pipeline := mongo.Pipeline{}
	if filterMap != nil {
		pipeline = append(pipeline, bson.D{{"$match", filterMap}})
	}
	if summaryMap != nil {
		pipeline = append(pipeline, bson.D{{"$group", summaryMap}})
	}
	cur, err = coll.Aggregate(ctx, pipeline)
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err).DataResult()
	}
	err = cur.All(ctx, &data)
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err).DataResult()
	}
	return data, true, nil

}

func (r *Dao[T]) CountRows(ctx context.Context, tenantId string, filterData any, opts ...ddd_repository.Options) (int64, error) {
	total, err := r.getCollection().CountDocuments(ctx, filterData)
	if err != nil {
		return 0, err
	}
	return total, err
}

func (r *Dao[T]) Count(ctx context.Context, tenantId string, rsql string, opts ...ddd_repository.Options) (int64, error) {
	f, err := r.getFilterMap(tenantId, rsql)
	total, err := r.getCollection().CountDocuments(ctx, f)
	if err != nil {
		return 0, err
	}
	return total, err
}

func (r *Dao[T]) DoList(tenantId, rsql string, fun func(filter map[string]interface{}) ([]T, bool, error)) *ddd_repository.FindListResult[T] {
	if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is empty")); err != nil {
		return ddd_repository.NewFindListResultError[T](err)
	}
	filterData, err := r.getFilterMap(tenantId, rsql)
	if err != nil {
		return ddd_repository.NewFindListResultError[T](err)
	}
	data, ok, err := fun(filterData)
	if err != nil {
		if errors.IsErrorMongoNoDocuments(err) {
			err = nil
		}
	}
	return ddd_repository.NewFindListResult(data, ok, err)
}

func (r *Dao[T]) DoFilter(tenantId, rsql string, fun func(filter map[string]interface{}) (*ddd_repository.FindPagingResult[T], bool, error)) *ddd_repository.FindPagingResult[T] {
	if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is empty")); err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}
	filterData, err := r.getFilterMap(tenantId, rsql)
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}
	data, _, err := fun(filterData)
	if err != nil {
		if errors.IsErrorMongoNoDocuments(err) {
			err = nil
		}
	}
	return data
}

func (r *Dao[T]) GetFilterMap(tenantId, rsqlstr string) (map[string]interface{}, error) {
	return r.getFilterMap(tenantId, rsqlstr)
}

func (r *Dao[T]) getFilterMap(tenantId, rsqlstr string) (map[string]interface{}, error) {
	p := NewMongoProcess()
	if err := rsql.ParseProcess(rsqlstr, p); err != nil {
		return nil, err
	}
	filterMap, err := p.GetFilter(tenantId)
	if err != nil {
		return nil, err
	}
	return filterMap, nil
}

func (r *Dao[T]) DoFindList(fun func() ([]T, bool, error)) *ddd_repository.FindListResult[T] {
	data, isFound, err := fun()
	if err != nil {
		if errors.IsErrorMongoNoDocuments(err) {
			isFound = false
			err = nil
		}
	}
	return ddd_repository.NewFindListResult[T](data, isFound, err)
}

func (r *Dao[T]) DoFindOne(fun func() (T, bool, error)) *ddd_repository.FindOneResult[T] {
	data, isFound, err := fun()
	if err != nil {
		if errors.IsErrorMongoNoDocuments(err) {
			isFound = false
			err = nil
		}
	}
	return ddd_repository.NewFindOneResult[T](data, isFound, err)
}

func (r *Dao[T]) DoSet(fun func() (T, error)) *ddd_repository.SetResult[T] {
	data, err := fun()
	return ddd_repository.NewSetResult[T](data, err)
}

func (r *Dao[T]) DoSetMany(fun func() ([]T, error)) *ddd_repository.SetManyResult[T] {
	data, err := fun()
	return ddd_repository.NewSetManyResult[T](data, err)
}

func (r *Dao[T]) DoSetManyCount(fun func() (*mongo.UpdateResult, error)) *ddd_repository.SetManyCountResult {
	res, err := fun()
	return ddd_repository.NewSetManyCountResult(res, err)
}

/*
func (r *Dao[T]) DoSetMap(fun func() (map[string]interface{}, error)) *ddd_repository.SetResult[map[string]interface{}] {
	data, err := fun()
	return ddd_repository.NewSetResult[T](data, err)
}
*/

//
//  getSort
//  @Description: 返回排序bson.D
//  @receiver r
//  @param sort  排序语句 "name:desc,id:asc"
//  @return bson.D
//  @return error
//
func (r *Dao[T]) getSort(sort string) (bson.D, error) {
	if len(sort) == 0 {
		return bson.D{}, nil
	}
	// 输入
	// name:desc,id:asc
	// 输出
	/*	sort := bson.D{
		bson.E{"update_time", -1},
		bson.E{"goods_id", -1},
	}*/
	res := bson.D{}
	list := strings.Split(sort, ",")
	for _, s := range list {
		sortItem := strings.Split(s, ":")
		name := sortItem[0]
		name = strings.Trim(name, " ")
		if name == "id" {
			name = ConstIdField
		}
		order := "asc"
		if len(sortItem) > 1 {
			order = sortItem[1]
			order = strings.ToLower(order)
			order = strings.Trim(order, " ")
		}

		// 其中 1 为升序排列，而-1是用于降序排列.
		orderVal := 1
		var oerr error
		switch order {
		case "asc":
			orderVal = 1
		case "desc":
			orderVal = -1
		default:
			oerr = errors.New("order " + order + " is error")
		}
		if oerr != nil {
			return nil, oerr
		}
		item := bson.E{name, orderVal}
		res = append(res, item)
	}
	return res, nil
}
