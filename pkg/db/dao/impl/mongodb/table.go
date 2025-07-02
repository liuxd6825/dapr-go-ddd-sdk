package mongodb

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/store_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongo_options "go.mongodb.org/mongo-driver/mongo/options"
)

type Table struct {
	tableName string
	schema    *store.DBSchema
	db        store_mongodb.IMongoDB
}

func NewTable(db store_mongodb.IMongoDB, schema *store.DBSchema) idao.Table {
	return newTable(db, schema)
}

func newTable(db store_mongodb.IMongoDB, schema *store.DBSchema) *Table {
	tableName := stringutils.AsFieldName(schema.Name)
	return &Table{db: db, tableName: tableName, schema: schema}
}

func (t *Table) GetTableName() string {
	return t.tableName
}

func (t *Table) GetSchema() *store.DBSchema {
	return t.schema
}

func (t *Table) AutoMigrate(ctx context.Context) {
	var err error
	if ctx == nil {
		ctx = context.Background()
	}
	validator := t.getValidator(ctx)
	if isExist, e := t.db.ExistCollection(ctx, t.tableName); e != nil {
		err = e
	} else if !isExist {
		opts := mongo_options.CreateCollection().
			SetValidator(validator).
			SetValidationAction("error"). // 违反规则时报错
			SetValidationLevel("strict")  // 严格验证
		err = t.db.CreateCollection(t.tableName, opts)
	} else {
		cmd := bson.D{
			{"collMod", t.tableName}, // 这里修改的集合名
			{"validator", validator},
			{"validationLevel", "strict"},
			{"validationAction", "error"},
		}
		var result bson.M
		err = t.db.GetDatabase().RunCommand(ctx, cmd).Decode(&result)
		if err == nil && result != nil {
			if val, ok := result["ok"]; ok {
				if fVal, fOk := val.(float64); fOk {
					if fVal != 1 {
						err = errors.New("Table.AutoMigrate() error")
					}
				}
			}
		}
	}
	if err != nil {
		panic(err)
	}

	err = t.CreateIndexes(ctx)
	if err != nil {
		panic(err)
	}
}

func (t *Table) Exist(ctx context.Context) bool {
	if ctx == nil {
		ctx = context.Background()
	}
	isExist, err := t.db.ExistCollection(ctx, t.tableName)
	if err != nil {
		panic(err)
	}
	return isExist
}

func (t *Table) Drop(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	var err error
	if isExist, e1 := t.db.ExistCollection(ctx, t.tableName); e1 != nil {
		err = e1
	} else if isExist {
		err = t.db.GetCollection(t.tableName).Drop(ctx)
	}
	if err != nil {
		panic(err)
	}
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
func (t *Table) CreateIndexes(ctx context.Context) error {
	var models []mongo.IndexModel
	for _, field := range t.schema.Fields {
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

	col := t.db.GetCollection(t.tableName)
	_, err := col.Indexes().DropAll(ctx)
	if err != nil {
		return err
	}
	_, err = col.Indexes().CreateMany(ctx, models)
	return err
}

func (t *Table) getValidator(ctx context.Context) bson.M {
	// 定义 JSON Schema 验证规则
	jsonSchema := t.getJsonSchema(ctx)
	validator := bson.M{
		"$jsonSchema": jsonSchema,
	}
	return validator
}

func (t *Table) getJsonSchema(ctx context.Context) bson.M {
	// 定义 JSON Schema 验证规则
	required := []string{}
	properties := bson.M{}
	for _, field := range t.schema.Fields {
		notNull := field.NotNull
		if notNull || field.PrimaryKey {
			required = append(required, field.DBName)
		}
		dataType := t.getBsonType(field.DataType)
		bsonType := []string{dataType}
		if !notNull {
			bsonType = append(bsonType, "null")
		}
		prop := bson.M{
			"bsonType":    bsonType,
			"description": field.Comment,
		}
		properties[field.DBName] = prop
	}
	jsonSchema := bson.M{
		"bsonType":   "object",
		"properties": properties,
	}
	if len(required) > 0 {
		jsonSchema["required"] = required // 必填字段]
	}
	return jsonSchema
}

/*
bsonType 值	说明
"double"	双精度浮点数 (64-bit)
"string"	UTF-8 字符串
"object"	嵌套文档（子对象）
"array"		数组
"binData"	二进制数据
"undefined"	已弃用（兼容旧版）
"objectId"	MongoDB 的 ObjectId（12字节唯一ID）
"bool"		布尔值 (true/false)
"date"		日期时间（UTC 时间戳）
"null"		null 值
"regex"		正则表达式
"dbPointer"	已弃用（数据库引用）
"javascript"	JavaScript 代码
"symbol"	已弃用（类似字符串）
"javascriptWithScope"	带作用域的 JavaScript 代码
"int"		32 位整数
"timestamp"	MongoDB 内部时间戳（不同于 "date"）
"long"		64 位整数
"decimal"	高精度小数（128-bit Decimal）
"minKey"	最小键（内部排序用）
"maxKey"	最大键（内部排序用）
*/
func (t *Table) getBsonType(dataType store.DataType) string {
	switch dataType {
	case store.DataType_Bool, store.DataType_Array, store.DataType_String, store.DataType_Date:
		return string(dataType)
	case store.DataType_Bytes:
		return "binData"
	case store.DataType_Time:
		return "date"
	case store.DataType_Float:
		return "double"
	case store.DataType_Int, store.DataType_Uint:
		return "long"
	case store.DataType_Json:
		return "object"
	default:
		panic(fmt.Sprintf("Table.getBsonType(dataType) invalid data type" + string(dataType)))
	}
	return ""
}
