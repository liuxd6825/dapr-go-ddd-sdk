package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoQueryAPI struct {
	dbKey    string
	service  *service.MongoService
	env      *envs.Env
	rootPath string
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoQueryAPI(env *envs.Env, rootPath string, dbKey string) *MongoQueryAPI {
	dbItem := env.GetDB(dbKey)
	if dbItem == nil {
		panic("mongoDbItem is nil")
	}
	client := dbItem.GetMongo().Client()
	database := dbItem.GetMongo().GetDatabase()
	mongoService := service.NewMongoService(env, rootPath, dbKey)
	return &MongoQueryAPI{
		dbKey:    dbKey,
		service:  mongoService,
		env:      env,
		rootPath: rootPath,
		client:   client,
		database: database,
	}
}

func (s *MongoQueryAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath+"/analysis/", "analysis.MongoQueryAPI", s)
	//controller.Post("aggregate", "Aggregate")
	return controller
}

func (s *MongoQueryAPI) Aggregate(ctx context.Context, qry *query.AggregateQuery) (any, error) {
	return s.service.Aggregate(ctx, qry)
}
