package restapp

import (
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strings"
)

type MongoManager struct {
}

func NewMongoManager() DbManager {
	return &MongoManager{}
}

func NewMongoScriptManager() DbScriptManager {
	return &MongoManager{}
}

func (m *MongoManager) GetInitScript(ctx context.Context, dbKey string, tables []*Table, env *EnvConfig, options *CreateOptions) string {
	cfg := env.Mongo[dbKey]
	user := cfg.UserName
	pwd := cfg.Password
	dbName := cfg.Database
	dbscript := fmt.Sprintf(`
db.createUser({
	user:"%s",
	pwd:"%s",
	roles:[
		{role:"userAdmin",db:"%s"},
		{role:"dbAdmin",db:"%s"},
		{role:"readWrite",db:"%s"}
	]
}); \r\n`, user, pwd, dbName, dbName, dbName)

	sb := strings.Builder{}
	sb.WriteString(dbscript)
	for _, table := range tables {
		sb.WriteString(m.getTableScript(ctx, env, dbKey, table, options))
	}
	return sb.String()
}

func (m *MongoManager) getTableScript(ctx context.Context, env *EnvConfig, dbKey string, table *Table, options *CreateOptions) string {
	opt := options
	if opt == nil {
		opt = NewCreateOptions()
	}
	collName := options.Prefix + table.TableName
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("use %s \r\n", dbKey))
	sb.WriteString(fmt.Sprintf("db.createCollection(\"%s\") \r\n", collName))

	e := table.Object
	t := reflect.TypeOf(e)
	elem := t.Elem()

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

		desc := f.Tag.Get("desc")
		if len(desc) == 0 {
			desc = f.Name
		}

		if isCreate {
			name := f.Tag.Get("bson")
			if len(name) == 0 {
				name = stringutils.SnakeString(f.Name)
			}
			if name != "_id" { //主键不需要创建
				name = name + "_"
				sb.WriteString(fmt.Sprintf("db.%v.createIndex({\"%s\":%v,\"description\":%v})", collName, name, order, desc))
				if isUnique {
					sb.WriteString(fmt.Sprintf(",{unique:true}"))
				}
				sb.WriteString("} \r\n")
			}
		}
	}
	return sb.String()
}

func (m *MongoManager) Create(ctx context.Context, table *Table, env *EnvConfig, options *CreateOptions) {
	opt := options
	if opt == nil {
		opt = NewCreateOptions()
	}
	opt.Printf("dbKey:%s; tableName:%s; ", table.DbKey, table.TableName)

	mongodb, ok := GetMongoByKey(table.DbKey)
	if !ok {
		opt.Log(errors.New(fmt.Sprintf("Dbkey:%s 不存在", table.DbKey)))
		return
	}

	opt.Printf("dbName:%s; ", mongodb.Name())
	collName := options.Prefix + table.TableName

	has, err := mongodb.ExistCollection(ctx, collName)
	if err != nil {
		opt.Log(err)
		return
	}
	if !has {
		err = mongodb.CreateCollection(collName)
		if err != nil {
			opt.Log(err)
			return
		}
	}

	coll := mongodb.GetCollection(collName)
	names, err := m.createIndexes(ctx, table, coll)
	if err != nil {
		opt.Log(err)
		return
	}

	opt.Printf("indexCount:%v; indexName:[", len(names))
	for _, n := range names {
		opt.Printf("%s,", n)
	}
	opt.Printf("]; ")

	id := idutils.NewId()
	if _, err = coll.InsertOne(ctx, bson.M{"_id": id}); err != nil {
		opt.Log(err)
		return
	}
	if _, err = coll.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		opt.Log(err)
		return
	}
	opt.Printf("status: success。\r\n")
	return
}

func (m *MongoManager) Update(ctx context.Context, table *Table, env *EnvConfig, options *UpdateOptions) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoManager) getDb(ctx context.Context, dbKey string) (*ddd_mongodb.MongoDB, error) {
	mongodb, ok := GetMongoByKey(dbKey)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Dbkey %s 不存在", dbKey))
	}
	return mongodb, nil
}

func (m *MongoManager) getCollection(ctx context.Context, dbKey string, prefix string, tableName string) (*mongo.Collection, error) {
	collName := prefix + tableName
	mongodb, ok := GetMongoByKey(collName)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Dbkey %s 不存在", dbKey))
	}
	coll := mongodb.GetCollection(collName)
	return coll, nil
}

func (m *MongoManager) createIndexes(ctx context.Context, table *Table, coll *mongo.Collection) ([]string, error) {
	e := table.Object
	t := reflect.TypeOf(e)
	elem := t.Elem()
	var models []mongo.IndexModel
	indexes, err := m.getIndexes(ctx, coll)
	if err != nil {
		return nil, err
	}
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
			if name != "_id" { //主键不需要创建
				name = name + "_"
				if _, ok := m.findIndex(indexes, name); !ok {
					model := mongo.IndexModel{
						Keys:    bson.D{{Key: name, Value: order}},
						Options: options.Index().SetUnique(isUnique).SetName(name),
					}
					models = append(models, model)
				}
			}
		}
	}

	if len(models) == 0 {
		return []string{}, nil
	}

	names, err := coll.Indexes().CreateMany(ctx, models)
	return names, err
}

func (m *MongoManager) findIndex(list []bson.M, name string) (bson.M, bool) {
	for _, item := range list {
		if item["name"] == name {
			return item, true
		}
	}
	return nil, false
}

func (m *MongoManager) getIndexes(ctx context.Context, coll *mongo.Collection) ([]bson.M, error) {
	cursor, err := coll.Indexes().List(ctx)
	if err != nil {
		return nil, err
	}

	var indexs []bson.M
	if err = cursor.All(context.TODO(), &indexs); err != nil {
		return nil, err
	}
	return indexs, err

}
