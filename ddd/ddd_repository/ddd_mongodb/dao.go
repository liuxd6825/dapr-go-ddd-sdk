package ddd_mongodb

import (
	"context"
	"fmt"
	"github.com/dapr/components-contrib/liuxd/common/utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/rsql"
	rsql_mongo2 "github.com/liuxd6825/dapr-go-ddd-sdk/db/rsql/rsql_mongo"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	assert2 "github.com/liuxd6825/dapr-go-ddd-sdk/errors/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongo_options "go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strings"
)

const (
	ConstIdField       = "_id"
	ConstTenantIdField = "tenant_id"
	TenantIdField      = "tenant_id"
)

type Dao[T any] struct {
	eb         ddd.EntityBuilder[T] // 实体构造器
	decoder    Decoder[T]           // mongo数据解码器
	collection *mongo.Collection
	mongodb    *MongoDB
	null       T
	newFun     func() T                                                                   // 新建实体结构方法
	initFun    func(ctx context.Context) (mongodb *MongoDB, collection *mongo.Collection) // 初始化
	options    *Options[T]
	metadata   map[string]any
	schema     *dbschema.Schema
}

func NewDao[T any](dbSch *dbschema.Schema, initFun func(ctx context.Context) (mongodb *MongoDB, collection *mongo.Collection), opts ...*Options[T]) ddd_repository.Dao[T] {
	return NewMongoDao(dbSch, initFun, opts...)
}

func NewMongoDao[T any](dbSch *dbschema.Schema, initFun func(ctx context.Context) (mongodb *MongoDB, collection *mongo.Collection), opts ...*Options[T]) *Dao[T] {
	r := &Dao[T]{
		metadata: make(map[string]any),
	}
	r.initFun = initFun
	r.options = NewOptions[T](opts...)
	if r.options.entityBuilder == nil {
		r.options.entityBuilder = ddd.NewAnyEntityBuilder[T](dbSch)
		r.decoder = NewStructDecoder[T]()
	}
	r.eb = r.options.entityBuilder
	if r.eb.GetConfig().IsMap {
		r.decoder = NewStructDecoder[T]()
	} else {
		r.decoder = NewMapDecoder[T]()
	}
	ctx := context.Background()
	mongodb, coll := initFun(ctx)
	if err := r.Init(ctx, mongodb, coll); err != nil {
		panic(err)
	}
	return r
}

func (r *Dao[T]) SetMetadata(metadata map[string]any) {
	r.metadata = metadata
}

func (r *Dao[T]) GetMetadata() map[string]any {
	return r.metadata
}

func (r *Dao[T]) AddMetadata(key string, val any) {
	r.metadata[key] = val
}

func (r *Dao[T]) GetSchema() *dbschema.Schema {
	return r.schema
}

func (r *Dao[T]) NewEntity() T {
	return r.eb.NewEntity()
}

func (r *Dao[T]) NewEntityList() []T {
	return r.eb.NewEntityList()
}

func (r *Dao[T]) GetTenantId(entity T) string {
	return r.eb.GetTenantId(entity)
}

func (r *Dao[T]) SetTenantId(entity T, tenantId string) {
	r.eb.SetTenantId(entity, tenantId)
}

func (r *Dao[T]) GetId(entity T) string {
	return r.eb.GetId(entity)
}

func (r *Dao[T]) SetId(entity T, id string) {
	r.eb.SetId(entity, id)
}

func (d *Dao[T]) GetAggId(entity T) string {
	return d.eb.GetAggId(entity)
}

func (r *Dao[T]) Init(ctx context.Context, mongodb *MongoDB, collection *mongo.Collection) error {
	r.mongodb = mongodb
	r.collection = collection
	/*
		if r.options.GetAutoCreateCollection() {
			find, err := mongodb.ExistCollection(ctx, collection.TableName())
			if err != nil {
				return err
			}

			if !find {
				if err := mongodb.CreateCollection(collection.TableName()); err != nil {
					return err
				}
				if r.options.GetAutoCreateIndex() {
					if err := r.CreateIndexes(ctx); err != nil {
						return err
					}
				}
			}
		}
	*/
	return nil
}

// CreateIndexes
//
//	 @Description: 根据index标签创建数据库索引。index:asc,desc,unique
//	 示例
//	 type Index struct {
//		    Id        string `bson:"_id" `
//		    TenantId  string
//		    TableName      string `bson:"name" index:"" `
//		    Asc       int64  `bson:"asc" index:" asc"`
//		    Desc      int64  `bson:"desc" index:" desc "`
//		    Unique    string `index:"unique"`
//		    AscUnique string `bson:"asc_unique" index:"asc, unique "`
//	  }
//	 @receiver r
//	 @param ctx  上下文
//	 @return error 错误
func (r *Dao[T]) CreateIndexes(ctx context.Context) error {
	e := r.NewEntity()
	t := reflect.TypeOf(e)
	elem := t.Elem()
	var models []mongo.IndexModel
	for i := 0; i < elem.NumField(); i++ {
		f := elem.Field(i)

		isCreate := false   //是否创建索引
		isUnique := false   //是否唯一
		var order int32 = 1 // 排序规则 -1:降序； 1:升序

		if strings.Contains(string(f.Tag), " index:") {
			isCreate = true
		}

		indexTag := f.Tag.Get("index")
		indexTag = strings.Trim(indexTag, " ")
		if len(indexTag) > 0 {
			isCreate = true
			if strings.Contains(indexTag, ",") {
				values := strings.Split(indexTag, ",")
				for _, key := range values {
					switch strings.Trim(key, " ") {
					case "unique":
						isUnique = true
						break
					case "asc":
						order = 1
						break
					case "desc":
						order = 0
						break
					}
				}
			} else {
				indexTag = strings.ToLower(strings.Trim(indexTag, " "))
				switch indexTag {
				case "unique":
					isUnique = true
					break
				case "asc":
					order = 1
					break
				case "desc":
					order = -1
					break
				}
			}
		} else if !isCreate {
			gormTag := f.Tag.Get("gorm")
			if strings.Contains(gormTag, "index:") {
				isCreate = true
			}
		}

		if isCreate {
			name := f.Tag.Get("bson")
			if len(name) == 0 {
				name = stringutils.SnakeString(f.Name)
			}
			model := mongo.IndexModel{
				Keys:    bson.D{{Key: name, Value: order}},
				Options: mongo_options.Index().SetUnique(isUnique).SetName(name + "_idx"),
			}
			models = append(models, model)
		}
	}

	if len(models) == 0 {
		return nil
	}

	col := r.getCollection(ctx)
	_, err := col.Indexes().DropAll(ctx)
	if err != nil {
		return err
	}

	_, err = col.Indexes().CreateMany(ctx, models)

	return err
}

func (r *Dao[T]) getCollection(ctx context.Context) *mongo.Collection {
	mongodb, coll := r.initFun(ctx)
	if err := r.Init(ctx, mongodb, coll); err != nil {
		panic(err)
	}
	return r.collection
}

func (r *Dao[T]) Save(ctx context.Context, data *ddd.SetData[T], opts ...ddd_repository.Options) (setResult *ddd_repository.SetResult[T]) {
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
			err = r.Insert(ctx, entity, opts...).GetError()
		case ddd.DataStatueUpdate:
			err = r.Update(ctx, entity, opts...).GetError()
		case ddd.DataStatueDelete:
			err = r.DeleteById(ctx, r.GetTenantId(entity), r.GetId(entity), opts...).GetError()
		case ddd.DataStatueCreateOrUpdate:
			err = r.InsertOrUpdate(ctx, entity, opts...).GetError()
		}
		if err != nil {
			return ddd_repository.NewSetResultError[T](err)
		}
	}
	return ddd_repository.NewSetResultError[T](nil)
}

func (r *Dao[T]) getSessionCtx(ctx context.Context) context.Context {
	sCtx := getSessionContext(ctx, r.mongodb.database.Name())
	if sCtx == nil {
		return ctx
	}
	return sCtx
}

func (r *Dao[T]) InsertOrUpdate(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResultEmpty[T]()
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

func (r *Dao[T]) Insert(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResultEmpty[T]()
	gp.Try(func() error {
		ctx := r.getSessionCtx(ctx)
		if err := assert2.NotEmpty(r.GetTenantId(entity), assert2.NewOptions("tenantId is empty")); err != nil {
			return err
		}
		r.eb.SetCreatedInfo(ctx, entity)
		mRes, err := r.getCollection(ctx).InsertOne(ctx, entity, getInsertOneOptions(opts...))
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
func (r *Dao[T]) InsertMap(ctx context.Context, tenantId string, data map[string]interface{}, opts ...ddd_repository.Options) (res *ddd_repository.SetResult[T]) {
	res = ddd_repository.NewSetResultEmpty[T]()
	gp.Try(func() error {
		ctx := r.getSessionCtx(ctx)
		if err := assert2.NotEmpty(tenantId, assert2.NewOptions("tenantId is empty")); err != nil {
			return err
		}
		data["tenant_id"] = tenantId
		inRes, err := r.getCollection(ctx).InsertOne(ctx, data, getInsertOneOptions(opts...))
		if inRes != nil && inRes.InsertedID != nil {
			res.SetRowsAffected(1)
		}
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (r *Dao[T]) InsertMany(ctx context.Context, tenantId string, entities []T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResultEmpty[T]()
	gp.Try(func() error {
		ctx := r.getSessionCtx(ctx)
		if entities == nil || len(entities) == 0 {
			return errors.New("entities is nil")
		}
		for _, e := range entities {
			if err := assert2.NotEmpty(r.GetTenantId(e), assert2.NewOptions("tenantId is empty")); err != nil {
				return err
			}
		}
		var docs []interface{}
		for _, e := range entities {
			docs = append(docs, e)
			r.eb.SetCreatedInfo(ctx, e)
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

func (r *Dao[T]) Update(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	return r.updateById(ctx, entity, opts...)
}

func (r *Dao[T]) updateById(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResultEmpty[T]()
	gp.Try(func() error {
		opt := ddd_repository.NewOptions(opts...)
		data := r.getUpdateData(entity, opt)
		uopt := getUpdateOptions(opts...)
		setData := bson.M{"$set": data}
		sCtx := r.getSessionCtx(ctx)
		mRes, err := r.getCollection(ctx).UpdateByID(sCtx, r.GetId(entity), setData, uopt)
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

func (r *Dao[T]) UpdateByRSQL(ctx context.Context, tenantId, filterRSQL string, data T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResultEmpty[T]()
	gp.Try(func() error {
		filter, err := r.getFilter(tenantId, filterRSQL)
		if err != nil {
			return err
		}
		if filter.IsAggregate() {
			return errors.New("aggregate update is not supported")
		}
		opt := ddd_repository.NewOptions(opts...)
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

func (r *Dao[T]) UpdateMany(ctx context.Context, tenantId string, entities []T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResultEmpty[T]()
	gp.Try(func() error {
		if entities == nil || len(entities) == 0 {
			return errors.New("entities is nil")
		}
		var list []mongo.WriteModel
		for _, entity := range entities {
			data := bson.M{"$set": entity}
			model := mongo.NewUpdateOneModel().SetFilter(bson.D{{"_id", r.GetId(entity)}}).SetUpdate(data).SetUpsert(true)
			list = append(list, model)
		}

		mRes, err := r.BulkWrite(ctx, list)
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

func (r *Dao[T]) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...ddd_repository.Options) (*ddd_repository.BulkWriteResult, error) {
	if len(models) == 0 {
		return &ddd_repository.BulkWriteResult{}, nil
	}
	opt := &mongo_options.BulkWriteOptions{}
	sCtx := r.getSessionCtx(ctx)
	bulkRes, err := r.getCollection(ctx).BulkWrite(sCtx, models, opt)
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

func (r *Dao[T]) UpdateManyMaskById(ctx context.Context, entities []T, mask []string, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResultEmpty[T]()
	gp.Try(func() error {
		if entities == nil || len(entities) == 0 {
			return errors.New("entities is nil")
		}

		for _, e := range entities {
			if err := assert2.NotEmpty(r.GetTenantId(e), assert2.NewOptions("tenantId is empty")); err != nil {
				return err
			}
		}

		var docs []interface{}
		for _, e := range entities {
			if len(mask) == 0 {
				docs = append(docs, e)
			} else {
				m := make(map[string]interface{})
				if err := types.MaskMapper(e, &m, mask); err != nil {
					return err
				}
				m[ConstIdField] = r.GetId(e)
				doc := asDocument(m)
				docs = append(docs, doc)
			}
		}
		sCtx := r.getSessionCtx(ctx)
		for _, doc := range docs {
			id, err := getDocumentId(doc)
			if err != nil {
				return err
			}
			updateOptions := getUpdateOptions(opts...)
			filter := bson.D{{ConstIdField, id}}

			setData := bson.M{"$set": doc}
			_, err = r.getCollection(ctx).UpdateOne(sCtx, filter, setData, updateOptions)
			if err != nil {
				return err
			}
		}
		return nil
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
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
		return ddd_repository.NewSetResultEmpty[T]().SetError(err)
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

func (r *Dao[T]) UpdateMapById(ctx context.Context, tenantId string, id string, data map[string]any, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	filter := bson.M{"tenant_id": tenantId, "_id": id}
	m := r.getMap(data)
	r.eb.SetUpdatedInfo(ctx, data)
	res := r.UpdateMapAndGetCount(ctx, tenantId, filter, m, opts...)
	return res
}

func (r *Dao[T]) FindOneAndUpdateById(ctx context.Context, tenantId string, id string, data map[string]any, opts ...ddd_repository.Options) (T, error) {
	var null T
	if err := assert2.NotEmpty(tenantId, assert2.NewOptions("tenantId is empty")); err != nil {
		return null, err
	}

	if err := assert2.NotEmpty(id, assert2.NewOptions("id is empty")); err != nil {
		return null, err
	}
	filter := bson.M{"tenant_id": tenantId, "_id": id}
	udpate := r.getMap(data)
	sCtx := r.getSessionCtx(ctx)
	_, err := r.getCollection(ctx).UpdateOne(sCtx, filter, udpate)
	if err != nil {
		return null, err
	}
	find := r.FindById(ctx, tenantId, id, opts...)
	return find.GetData(), find.GetError()
}

func (r *Dao[T]) UpdateMap(ctx context.Context, tenantId string, id string, data map[string]any, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResultEmpty[T]()
	gp.Try(func() error {
		filter, err := r.getFilter(tenantId, fmt.Sprintf("id=='%s'", id))
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

func (r *Dao[T]) UpdateMapAndGetCount(ctx context.Context, tenantId string, filter any, data any, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResultEmpty[T]()
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
		switch fieldName {
		case "$and", "$or":
			ands = append(ands, bson.M{fieldName: fieldValue})
		default:
			if fieldName != ConstIdField {
				fieldName = AsFieldName(fieldName)
			}
			item := bson.M{fieldName: bson.M{"$eq": fieldValue}}
			ands = append(ands, item)
		}
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
		data := r.NewEntity()
		sCtx := r.getSessionCtx(ctx)
		result := r.getCollection(ctx).FindOne(sCtx, filter, findOneOptions)
		if err := r.DecodeSingle(ctx, result, data); err != nil {
			return null, false, err
		}
		return data, true, nil
	})
}

func (r *Dao[T]) DecodeSingle(ctx context.Context, result *mongo.SingleResult, data T) error {
	return r.decoder.Single(ctx, result, data)
}

func (r *Dao[T]) DecodeList(ctx context.Context, cursor *mongo.Cursor, data *[]T) error {
	return r.decoder.List(ctx, cursor, data)
}

func (r *Dao[T]) FindListByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	return r.DoFindList(func() ([]T, bool, error) {
		var list []T
		filter := r.NewFilter(tenantId, filterMap)
		findOptions := getFindOptions(opts...)
		sCtx := r.getSessionCtx(ctx)
		cursor, err := r.getCollection(ctx).Find(sCtx, filter, findOptions)
		if err != nil {
			return nil, false, err
		}
		err = cursor.All(ctx, &list)
		count := len(list)
		if count == 0 {
			list = []T{}
		}
		return list, count > 0, err
	})
}

func (r *Dao[T]) FindListByBsonM(ctx context.Context, tenantId string, filter bson.M, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	return r.DoFindList(func() ([]T, bool, error) {
		var list []T
		sCtx := r.getSessionCtx(ctx)
		findOptions := getFindOptions(opts...)
		cursor, err := r.getCollection(ctx).Find(sCtx, filter, findOptions)
		if err != nil {
			return nil, false, err
		}
		err = cursor.All(ctx, &list)
		count := len(list)
		if count == 0 {
			list = []T{}
		}
		return list, len(list) > 0, err
	})
}

func (r *Dao[T]) FindByRSQL(ctx context.Context, tenantId string, rsql string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	return r.doList(tenantId, rsql, func(filter *rsql_mongo2.Filter) ([]T, bool, error) {
		list := r.NewEntityList()
		ctx = r.getSessionCtx(ctx)
		findOpts := &findByFilterOptions{
			resultsData: &list,
		}
		err := r.findByFilter(ctx, filter, findOpts)
		if err != nil {
			return nil, false, err
		}
		return list, len(list) > 0, err
	})
}

func (r *Dao[T]) FindAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	return r.FindListByMap(ctx, tenantId, nil, opts...)
}

func (r *Dao[T]) findPaging(ctx context.Context, query ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	return r.doFilter(query.GetTenantId(), query.GetFilter(), func(filter *rsql_mongo2.Filter) (*ddd_repository.FindPagingResult[T], bool, error) {
		if err := assert2.NotEmpty(query.GetTenantId(), assert2.NewOptions("tenantId is empty")); err != nil {
			return nil, false, err
		}
		ctx = r.getSessionCtx(ctx)
		data := r.NewEntityList()

		inOut := newFindOption(filter, query, &data)
		err := r.find(ctx, inOut)
		if err != nil {
			return nil, false, err
		}

		findData := ddd_repository.NewFindPagingResult[T](data, inOut.totalRows, query, err)
		return findData, findData.IsFound, err
	})
}

type findOption struct {
	filter    *rsql_mongo2.Filter            // rsql的查询条件
	query     ddd_repository.FindPagingQuery // 分页查询条件
	results   any                            // 返回数据
	totalRows int64                          // 返回记录数
}

func newFindOption(filter *rsql_mongo2.Filter, query ddd_repository.FindPagingQuery, results any) *findOption {
	return &findOption{
		filter:    filter,
		query:     query,
		results:   results,
		totalRows: 0,
	}
}

func (r *Dao[T]) newFindOptions(query ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) (*mongo_options.FindOptions, error) {
	findOptions := getFindOptions(opts...)
	if query == nil {
		return findOptions, nil
	}

	if query.GetPageSize() > 0 {
		findOptions.SetLimit(query.GetPageSize())
		findOptions.SetSkip(query.GetPageSize() * query.GetPageNum())
	}
	if len(query.GetSort()) > 0 {
		sort, err := r.getSort(query.GetSort())
		if err != nil {
			return nil, err
		}
		findOptions.SetSort(sort)
	}

	if projection := r.getFindOptionsProjection(query); projection != nil {
		findOptions.SetProjection(projection)
	}
	return findOptions, nil
}

func (r *Dao[T]) getFindOptionsProjection(query ddd_repository.FindPagingQuery) bson.D {
	var projection bson.D
	if len(query.GetFields()) > 0 {
		fields := strings.Split(query.GetFields(), ",")
		projection = bson.D{}
		for _, f := range fields {
			key := stringutils.SnakeString(strings.Trim(f, " "))
			projection = append(projection, bson.E{Key: key, Value: 1})
		}
	}
	return projection
}

/*
	func (r *Dao[T]) FindPaging2(ctx context.Context, query ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
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

		coll := r.getCollection(ctx)
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

func (r *Dao[T]) FindAutoComplete(ctx context.Context, qry ddd_repository.FindAutoCompleteQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	f := ddd_repository.NewFindPagingQuery()
	groupCols := []*ddd_repository.GroupCol{
		{Field: qry.GetField(), DataType: types.DataTypeString},
	}

	f.SetGroupCols(groupCols)
	f.SetTenantId(qry.GetTenantId())
	f.SetFields(qry.GetFields())
	f.SetFilter(qry.GetFilter())
	f.SetMustFilter(qry.GetMustWhere())

	f.SetPageNum(qry.GetPageNum())
	f.SetPageSize(qry.GetPageSize())
	f.SetSort(qry.GetSort())
	f.SetIsTotalRows(false)

	return r.FindPaging(ctx, f, opts...)
}

func (r *Dao[T]) FindDistinct(ctx context.Context, qry ddd_repository.FindDistinctQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	f := ddd_repository.NewFindPagingQuery()

	f.SetGroupCols(qry.GetGroupCols())
	f.SetTenantId(qry.GetTenantId())
	f.SetFields(qry.GetFields())
	f.SetFilter(qry.GetFilter())
	f.SetMustFilter(qry.GetMustWhere())

	f.SetPageNum(qry.GetPageNum())
	f.SetPageSize(qry.GetPageSize())
	f.SetSort(qry.GetSort())
	f.SetIsTotalRows(false)

	return r.FindPaging(ctx, f, opts...)
}

func (r *Dao[T]) AggregateByPipeline(ctx context.Context, pipeline mongo.Pipeline, data interface{}, opts ...ddd_repository.Options) error {
	sCtx := r.getSessionCtx(ctx)
	options := getAggregateOptions(opts...)
	cur, err := r.getCollection(ctx).Aggregate(sCtx, pipeline, options)
	if err != nil {
		return err
	}
	err = cur.All(ctx, data)
	return err
}

func (r *Dao[T]) CopyTo(ctx context.Context, tenantId string, rsql string, toCollectionName string, opts ...ddd_repository.Options) error {
	options := getAggregateOptions(opts...)
	//db.record.aggregate([{$match:{opp_bank_name:"工商银行"}},{$out:"record1"}])
	filter, err := r.getFilter(tenantId, rsql)
	if err != nil {
		return err
	}
	pipeline := mongo.Pipeline{}
	if filter != nil && len(filter.Match) > 0 {
		pipeline = append(pipeline, bson.D{{"$match", filter.Match}})
	}
	pipeline = append(pipeline, bson.D{{"$out", toCollectionName}})
	sCtx := r.getSessionCtx(ctx)
	_, err = r.getCollection(ctx).Aggregate(sCtx, pipeline, options)
	return err
}

func (r *Dao[T]) SumEntity(ctx context.Context, qry ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) ([]T, bool, error) {
	data := r.NewEntityList()
	sCtx := r.getSessionCtx(ctx)
	_, found, err := r.SumByQuery(sCtx, qry, &data, opts...)
	return data, found, err
}

func (r *Dao[T]) SumMap(ctx context.Context, qry ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) ([]map[string]any, bool, error) {
	data := make([]map[string]any, 0)
	_, found, err := r.SumByQuery(r.getSessionCtx(ctx), qry, &data, opts...)
	return data, found, err
}

func (r *Dao[T]) SumByQuery(ctx context.Context, qry ddd_repository.FindPagingQuery, data any, opts ...ddd_repository.Options) (any, bool, error) {
	if len(qry.GetValueCols()) == 0 {
		return nil, false, nil
	}

	var err error
	process := rsql_mongo2.NewProcess(qry.GetTenantId())

	f1 := qry.GetFilter()
	f2 := qry.GetMustFilter()
	f3 := ""
	mustWhere, ok := qry.(ddd_repository.FindPagingQueryMustWhere)
	if ok {
		f3, err = mustWhere.GetMustWhere()
		if err != nil {
			return nil, false, err
		}
	}
	filterRSQL := getRsqlAnds(f1, f2, f3)

	if err := rsql.ParseProcess(filterRSQL, process); err != nil {
		return nil, false, err
	}
	filter, ok := process.GetFilter().(*rsql_mongo2.Filter)
	if !ok {
		return nil, false, errors.New("filter does not implement Filter")
	}

	_, found, err := r.sum(ctx, filter, qry.GetValueCols(), data, opts...)
	return data, found, err
}

func (r *Dao[T]) SumByRSQL(ctx context.Context, tenantId string, rSql string, valueCols []*ddd_repository.ValueCol, opts ...ddd_repository.Options) map[string]any {
	filter, err := r.getFilter(tenantId, rSql)
	if err != nil {
		panic(err)
	}
	var list []map[string]any
	_, _, err = r.sum(ctx, filter, valueCols, &list, opts...)
	if err != nil {
		panic(err)
	}
	return list[0]
}

func (r *Dao[T]) sum(ctx context.Context, filter *rsql_mongo2.Filter, valueCols []*ddd_repository.ValueCol, list any, opts ...ddd_repository.Options) (any, bool, error) {
	coll := r.getCollection(ctx)
	/*
		var match map[string]any
		if filter, ok := filterMap.(*rsql_mongo.Filter); ok {
			match = filter.Match
		} else if m, ok := filterMap.(map[string]any); ok {
			match = m
		} else {
			return nil, false, errors.New("sum() filterMap is not map[string]any")
		}
	*/
	var cur *mongo.Cursor
	summaryMap := make(map[string]interface{})
	summaryMap["_id"] = "total"
	for _, col := range valueCols {
		field := utils.SnakeString(col.Field)
		aggFunc := utils.SnakeString(col.AggFunc.Name())

		var fieldValue any = "$" + field
		if aggFunc == ddd_repository.AggFuncCount.Name() {
			field = "count_" + field
			aggFunc = "$sum"
			fieldValue = 1
		} else {
			aggFunc = "$" + aggFunc
		}

		summaryMap[field] = map[string]interface{}{aggFunc: fieldValue}
	}

	pipeline := mongo.Pipeline{}
	filter.AddPipelines(pipeline)

	if summaryMap != nil {
		pipeline = append(pipeline, bson.D{{"$group", summaryMap}})
	}

	cur, err := coll.Aggregate(r.getSessionCtx(ctx), pipeline)
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err).DataResult()
	}
	err = cur.All(ctx, list)
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err).DataResult()
	}
	return list, true, nil

}

func (r *Dao[T]) CountByMap(ctx context.Context, tenantId string, filterData any, opts ...ddd_repository.Options) (int64, error) {
	total, err := r.getCollection(ctx).CountDocuments(r.getSessionCtx(ctx), filterData)
	if err != nil {
		return 0, err
	}
	return total, err
}

func (r *Dao[T]) CountByRSQL(ctx context.Context, tenantId string, rsql string, opts ...ddd_repository.Options) (int64, error) {
	f, err := r.getFilter(tenantId, rsql)
	if err != nil {
		return 0, err
	}
	ctx = r.getSessionCtx(ctx)
	total, err := r.getCollection(ctx).CountDocuments(ctx, f.Match)
	if err != nil {
		return 0, err
	}
	return total, err
}

func (r *Dao[T]) doList(tenantId, rsql string, fun func(filter *rsql_mongo2.Filter) ([]T, bool, error)) *ddd_repository.FindListResult[T] {
	if err := assert2.NotEmpty(tenantId, assert2.NewOptions("tenantId is empty")); err != nil {
		return ddd_repository.NewFindListResultError[T](err)
	}
	filterData, err := r.getFilter(tenantId, rsql)
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

func (r *Dao[T]) doFilter(tenantId, rsql string, fun func(filter *rsql_mongo2.Filter) (*ddd_repository.FindPagingResult[T], bool, error)) *ddd_repository.FindPagingResult[T] {
	if err := assert2.NotEmpty(tenantId, assert2.NewOptions("tenantId is empty")); err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}
	filterData, err := r.getFilter(tenantId, rsql)
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

func (r *Dao[T]) GetFilterMap(tenantId, rsql string) *rsql_mongo2.Filter {
	data, err := r.getFilter(tenantId, rsql)
	if err != nil {
		panic(err)
	}
	return data
}

func (r *Dao[T]) getFilter(tenantId, rSql string) (*rsql_mongo2.Filter, error) {
	process := rsql_mongo2.NewProcess(tenantId)
	if err := rsql.ParseProcess(rSql, process); err != nil {
		return nil, err
	}
	filter := process.GetFilter().(*rsql_mongo2.Filter)
	return filter, nil
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
	return ddd_repository.NewSetResult[T]().SetError(err).SetData(data)
}

func (r *Dao[T]) StartTx(ctx context.Context, txFun ddd_repository.TxFunc, options ...*ddd_repository.SessionOptions) (err error) {
	return StartTx(ctx, r.mongodb, r.mongodb.Name(), txFun, options...)
}

/*
func (r *Dao[T]) DoSetMap(fun func() (map[string]interface{}, error)) *ddd_repository.SetResult[map[string]interface{}] {
	data, err := fun()
	return ddd_repository.NewSetResult[T](data, err)
}
*/

// getSort
// @Description: 返回排序bson.D
// @receiver r
// @param sort  排序语句 "name:desc,id:asc"
// @return bson.D
// @return error
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

/*
func (r *Dao[T]) getDocuments(entities []T) []any {
	var list []any
	for _, item := range entities {
		list = append(list, r.getDocument(item))
	}
	return list
}

func (r *Dao[T]) getDocument(entity any) any {
	if dataMap, ok := entity.(DataMap); ok {
		return dataMap.GetDataMap()
	}
	return entity
}

*/
