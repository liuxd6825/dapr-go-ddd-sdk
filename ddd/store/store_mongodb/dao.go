package store_mongodb

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	assert2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/mapperutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	mongo_options "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	ConstIdField       = "id"
	ConstTenantIdField = "tenant_id"
	TenantIdField      = "tenant_id"
)

type IMongoDB interface {
	GetClient() *mongo.Client
	GetOperationTimeout() time.Duration
	GetDatabase() *mongo.Database
	GetServerCount() int
	ExistCollection(ctx context.Context, name any) (bool, error)
	GetCollection(collectionName string) *mongo.Collection
	CreateCollection(collectionName string, ops ...*mongo_options.CreateCollectionOptions) error
}

type ObjectId string

type InitOptionsFunc = func(options *mongo_options.ClientOptions) error

func (i ObjectId) String() string {
	return string(i)
}

type Dao[T any] struct {
	eb         store.EntityBuilder[T] // 实体构造器
	decoder    Decoder[T]             // mongo数据解码器
	collection *mongo.Collection
	mongodb    IMongoDB
	null       T
	newFun     func() T                                                                   // 新建实体结构方法
	initFun    func(ctx context.Context) (mongodb IMongoDB, collection *mongo.Collection) // 初始化
	options    *Options[T]
	metadata   map[string]any
	schema     *store.DBSchema
}

func NewDao[T any](dbSch *store.DBSchema, initFun func(ctx context.Context) (mongodb IMongoDB, collection *mongo.Collection), opts ...*Options[T]) store.IStore[T] {
	return NewMongoDao(dbSch, initFun, opts...)
}

func NewMongoDao[T any](dbSch *store.DBSchema, initFun func(ctx context.Context) (mongodb IMongoDB, collection *mongo.Collection), opts ...*Options[T]) *Dao[T] {
	if dbSch == nil {
		panic("NewMongoDao() dbSch is nil")
	}
	if dbSch.GormSchema == nil {
		dbSch.GormSchema = dbschema.NewGormSchema(dbSch)
	}
	r := &Dao[T]{
		metadata: make(map[string]any),
		schema:   dbSch,
	}

	r.initFun = initFun
	r.options = NewOptions[T](opts...)
	if r.options.entityBuilder == nil {
		r.options.entityBuilder = store.NewAnyEntityBuilder[T](dbSch)
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

func (r *Dao[T]) GetSchema() *store.DBSchema {
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

func (r *Dao[T]) GetAggId(entity T) string {
	return r.eb.GetAggId(entity)
}

func (r *Dao[T]) Init(ctx context.Context, mongodb IMongoDB, collection *mongo.Collection) error {
	r.mongodb = mongodb
	r.collection = collection

	return nil
}

func (r *Dao[T]) CreateCollection(ctx context.Context) error {
	if r.options.GetAutoCreateCollection() {
		find, err := r.mongodb.ExistCollection(ctx, r.schema.TableName)
		if err != nil {
			return err
		}

		if !find {
			validator := r.CreateValidator(ctx)
			opts := &mongo_options.CreateCollectionOptions{
				Validator: validator,
			}
			if err := r.mongodb.CreateCollection(r.schema.TableName, opts); err != nil {
				return err
			}

			if r.options.GetAutoCreateIndex() {
				if err := r.CreateIndexes(ctx); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (r *Dao[T]) CreateValidator(ctx context.Context) bson.M {
	// 定义 JSON Schema 验证规则
	var required []string
	properties := bson.M{}
	for _, field := range r.schema.Fields {
		if field.NotNull {
			required = append(required, field.DBName)
		}
		prop := bson.M{
			"bsonType":    field.DataType,
			"description": field.Comment,
		}
		properties[field.DBName] = prop

	}
	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType":   "object",
			"required":   required, // 必填字段
			"properties": properties,
		},
	}
	return validator
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

	var models []mongo.IndexModel
	for _, field := range r.schema.Fields {
		if field.PrimaryKey {
			// 为业务ID创建唯一索引
			indexModel := mongo.IndexModel{
				Keys:    bson.D{{Key: field.DBName, Value: 1}},
				Options: mongo_options.Index().SetUnique(true),
			}
			models = append(models, indexModel)

			// 可以同时为tenant_id + id创建复合唯一索引
			tenantIdIndex := mongo.IndexModel{
				Keys:    bson.D{{Key: "tenant_id", Value: 1}, {Key: field.DBName, Value: 1}},
				Options: mongo_options.Index().SetUnique(true),
			}
			models = append(models, tenantIdIndex)
		}

		if field.Unique {
			model := mongo.IndexModel{
				Keys:    bson.D{{Key: field.DBName, Value: 1}},
				Options: mongo_options.Index().SetUnique(field.Unique).SetName(field.DBName + "_unique_idx"),
			}
			models = append(models, model)
		}

		var order int32 = 0 // 排序规则 -1:降序； 1:升序
		switch field.OrderType {
		case store.OrderType_Asc:
			order = 1
		case store.OrderType_Desc:
			order = -1
		}
		if order != 0 {
			model := mongo.IndexModel{
				Keys:    bson.D{{Key: field.DBName, Value: order}},
				Options: mongo_options.Index().SetName(field.DBName + "_order_idx"),
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
	return coll
}

func (r *Dao[T]) Save(ctx context.Context, data *ddd.SetData[T], opts ...store.Options) (setResult *store.SetResult[T]) {
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
			err = r.Insert(ctx, entity, opts...).GetError()
		case ddd.DataStatueUpdate:
			err = r.Update(ctx, entity, opts...).GetError()
		case ddd.DataStatueDelete:
			err = r.DeleteById(ctx, r.GetTenantId(entity), r.GetId(entity), opts...).GetError()
		case ddd.DataStatueCreateOrUpdate:
			err = r.InsertOrUpdate(ctx, entity, opts...).GetError()
		}
		if err != nil {
			return store.NewSetResultError[T](err)
		}
	}
	return store.NewSetResultError[T](nil)
}

func (r *Dao[T]) getSessionCtx(ctx context.Context) context.Context {
	sCtx := getSessionContext(ctx, r.mongodb.GetDatabase().Name())
	if sCtx == nil {
		return ctx
	}
	return sCtx
}

func (r *Dao[T]) Merge(ctx context.Context, entity T, fields map[string]string, opts ...store.Options) (setResult *store.SetResult[T]) {
	panic("implement me")
}

func (r *Dao[T]) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...store.Options) (*store.BulkWriteResult, error) {
	if len(models) == 0 {
		return &store.BulkWriteResult{}, nil
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

func (r *Dao[T]) newBulkWriteResult(bulkRes *mongo.BulkWriteResult) *store.BulkWriteResult {
	res := store.NewBulkWriteResult()
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

func (r *Dao[T]) UpdateManyMaskById(ctx context.Context, entities []T, mask []string, opts ...store.Options) *store.SetResult[T] {
	res := store.NewSetResultEmpty[T]()
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
				if err := mapperutils.MaskMapper(e, &m, mask); err != nil {
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
func (r *Dao[T]) getDbMap(m map[string]any) map[string]any {
	data := map[string]any{}
	for _, field := range r.schema.Fields {
		data[field.DBName] = m[field.Name]
	}
	return data
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
			item := bson.M{fieldName: bson.M{"$eq": fieldValue}}
			ands = append(ands, item)
		}
	}
	filter := bson.M{
		"$and": ands,
	}
	return filter
}

func (r *Dao[T]) DecodeSingle(ctx context.Context, result *mongo.SingleResult, data T) error {
	return r.decoder.Single(ctx, result, data)
}

func (r *Dao[T]) DecodeList(ctx context.Context, cursor *mongo.Cursor, data *[]T) error {
	return r.decoder.List(ctx, cursor, data)
}

// AsFieldName
// @Description: 转换为mongodb规范的字段名称
// @param name
// @return string
func AsFieldName(name string) string {
	return stringutils.SnakeString(name)
}

// entity2db
func (r *Dao[T]) entity2db(entity any) map[string]any {
	data := map[string]any{}
	entityIsMap := r.eb.GetConfig().IsMap
	eMap, isMap := entity.(map[string]any)
	if entityIsMap {
		if !isMap {
			panic("entity is not map")
		}
	}

	if isMap {
		if entityIsMap {
			for _, field := range r.schema.Fields {
				fieldValue, ok := eMap[field.Name]
				if !ok && field.Name != field.DBName {
					fieldValue, ok = eMap[field.DBName]
				}
				if !ok {
					fieldValue = field.DefaultValueInterface
				}
				data[field.DBName] = fieldValue
			}
		} else {
			newMap := map[string]any{}
			for key, val := range eMap {
				field := r.schema.LookedField(key)
				if field != nil {
					newMap[field.DBName] = val
				}
			}
			return newMap
		}
	} else {
		for _, field := range r.schema.Fields {
			val := reflectutils.GetField(entity, field.Name)
			data[field.DBName] = val
		}
	}
	return data
}

func (r *Dao[T]) db2entity(data map[string]any) T {
	entity := r.NewEntity()
	isMap := r.eb.GetConfig().IsMap
	if isMap {
		eAny := any(entity)
		eMap, isMap := eAny.(map[string]any)
		if !isMap {
			panic("store_mongodb.dao entity is not a map")
		}
		for _, field := range r.schema.Fields {
			val := data[field.DBName]
			if field.DataType == store.DataType_Date || field.DataType == store.DataType_Time {
				if pDate, ok := val.(primitive.DateTime); ok {
					timeVal := pDate.Time().In(times.GetLocalTimeZone())
					val = &timeVal
				}
			} else if field.DataType == store.DataType_Array {
				if arr, ok := val.(primitive.A); ok {
					timeVal := arr
					val = &timeVal
				}
			}
			eMap[field.Name] = val
		}
	} else {
		var errFieldName string
		gp.Try(func() error {
			for _, field := range r.schema.Fields {
				errFieldName = field.Name
				val := data[field.DBName]
				if field.DataType == store.DataType_Date || field.DataType == store.DataType_Time {
					if pDate, ok := val.(primitive.DateTime); ok {
						timeVal := pDate.Time()
						val = &timeVal
					}
				} else if field.DataType == store.DataType_Array {
					if arr, ok := val.(primitive.A); ok {
						timeVal := arr
						val = &timeVal
					}
				}
				if val != nil {
					err := reflectutils.SetField(entity, field.Name, val)
					if err != nil {
						return errors.New("set field %s to entity %T error: %s", field.Name, entity, err.Error())
					}
				}
			}
			return nil
		}).Catch(func(err error) {
			panic(fmt.Sprintf("转换db值到属性%s时出错：%s", errFieldName, err.Error()))
		})
	}
	return entity
}
