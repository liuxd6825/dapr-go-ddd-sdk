package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const HumanSchema = `
	{
	  "name": "human",
	  "type": "object",
	  "required": ["id","name"],
	  "description": "人员基本信息",
	  "properties": {
		"age": {
		  "type": ["integer", "null"],
		  "title": "年龄",
		  "order": 100
		},
		"analyse": {
		  "type": ["string", "null"],
		  "title": "分析状态",
		  "order": 101
		},
		"birthday": {
		  "type": ["date", "null"],
		  "title": "出生日期",
		  "order": 102
		},
		"gender": {
		  "type": ["string", "null"],
		  "title": "性别",
		  "order": 103
		},
		"name": {
		  "type": ["string", "null"],
		  "title": "姓名",
		  "order": 104
		},
		"peopleType": {
		  "type": "array",
		  "items": { "type": "string" },
		  "title": "人员类型",
		  "order": 105
		},
		"tags": {
		  "title": "标签",
		  "type": "array",
		  "items": { "type": "string" },
		  "order": 106
		},
		"remark": {
		  "name": "remark",
		  "type": ["string", "null"],
		  "title": "备注"
		},
		"createdTime": {
		  "name": "createdTime",
		  "type": ["datetime"],
		  "title": "创建时间",
		  "order": 9002
		},
		"creatorId": {
		  "name": "creatorId",
		  "type": ["string","null"],
		  "title": "创建人Id",
		  "notNull": true,
		  "order": 9003
		},
		"creatorName": {
		  "name": "creatorName",
		  "type": ["string","null"],
		  "title": "创建人",
		  "notNull": true,
		  "order": 9004
		},
		"updatedTime": {
		  "name": "updatedTime",
		  "type": ["datetime","null"],
		  "title": "修改时间",
		  "notNull": true,
		  "order": 9005
		},
		"updaterId": {
		  "name": "updaterId",
		  "type": ["string","null"],
		  "title": "修改人Id",
		  "notNull": true,
		  "order": 9006
		},
		"updaterName": {
		  "name": "updaterName",
		  "type": ["string","null"],
		  "title": "修改人",
		  "notNull": true,
		  "order": 9007
		}
	  }
	}
	`

var client *mongo.Client

func init() {
	// 设置MongoDB连接URL
	clientOptions := options.Client().ApplyURI("mongodb://192.168.120.224:27018,192.168.120.224:27019,192.168.120.224:27020/?retryWrites=false&replicaSet=mongors&readPreference=primary&serverSelectionTimeoutMS=5000&connectTimeoutMS=10000")

	// 连接到MongoDB
	clientVal, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	client = clientVal
	// 确保连接成功
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getCollection(dbName string, collName string) *mongo.Collection {
	// 选择数据库和集合
	collection := client.Database(dbName).Collection(collName)
	return collection
}
