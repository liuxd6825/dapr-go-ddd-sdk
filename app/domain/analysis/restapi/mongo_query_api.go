package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoQueryAPI struct {
	accountService *service.SuAccountService
	env            *env.Env
	rootPath       string
	client         *mongo.Client
	database       *mongo.Database
}

func NewMongoQueryAPI(env *env.Env, rootPath string) *MongoQueryAPI {
	dbKey := env.GetDBKeyValue(config.DBKey)
	dbItem := env.GetDB(dbKey)
	if dbItem == nil {
		panic("mongoDbItem is nil")
	}
	client := dbItem.GetMongo().Client()
	database := dbItem.GetMongo().GetDatabase()
	return &MongoQueryAPI{
		env:      env,
		rootPath: rootPath,
		client:   client,
		database: database,
	}
}

func (s *MongoQueryAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath+"/analysis/", "analysis.MongoQueryAPI", s)
	controller.Post("aggregate", "Aggregate")
	return controller
}

func (s *MongoQueryAPI) Aggregate(ctx context.Context, qry *query.MongoQuery) (any, error) {

	// 1. 从请求体中解码聚合管道
	var pipeline = qry.Pipeline
	// 2. 检查管道是否为空
	if len(pipeline) == 0 {
		return nil, errors.New("聚合管道不能为空")
	}

	// 3. 获取集合句柄
	// 为简单起见，我们硬编码了数据库和集合名称
	// 在实际应用中，可以考虑从URL路径中动态获取
	collection := s.database.Collection(qry.Collection)

	// 4. 执行聚合查询
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.New("数据库聚合查询失败: " + err.Error())
	}
	defer cursor.Close(ctx)

	// 5. 将查询结果解码到通用结构中
	// bson.M 是 map[string]interface{} 的别名，可以表示任何BSON文档
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, errors.New("解析查询结果失败: " + err.Error())
	}

	// 如果没有结果，返回一个空数组，这比返回 null 更友好
	if results == nil {
		results = make([]bson.M, 0)
	}

	// 6. 设置响应头并返回JSON结果
	return results, nil
}
