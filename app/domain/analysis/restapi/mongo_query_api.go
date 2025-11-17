package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoQueryAPI struct {
	service  *service.MongoService
	env      *env.Env
	rootPath string
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoQueryAPI(env *env.Env, rootPath string) *MongoQueryAPI {
	dbKey := env.GetDBKeyValue(config.DBKey)
	dbItem := env.GetDB(dbKey)
	if dbItem == nil {
		panic("mongoDbItem is nil")
	}
	client := dbItem.GetMongo().Client()
	database := dbItem.GetMongo().GetDatabase()
	mongoService := service.NewMongoService(env, rootPath)
	return &MongoQueryAPI{
		service:  mongoService,
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

	return s.service.Aggregate(ctx, qry)
}
