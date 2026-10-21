package restapi

import (
	"github.com/kataras/iris/v12"
	analysisconfigs "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/configs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func Register(app *iris.Application, baseUrl string, env *envs.Env) {

	analysisconfigs.Init(env)
	mongoDbKey := analysisconfigs.GetMongoDbKey()
	dorisDbKey := analysisconfigs.GetDorisDbKey()

	restapi.RegisterController(app, NewSuTaskApi(env, baseUrl, mongoDbKey))
	restapi.RegisterController(app, NewSuTaskAccountApi(env, baseUrl, mongoDbKey))
	restapi.RegisterController(app, NewSuTaskRecordApi(env, baseUrl, mongoDbKey))
	restapi.RegisterController(app, NewMongoQueryAPI(env, baseUrl, mongoDbKey))
	restapi.RegisterController(app, NewRecordChartQueryAPI(env, baseUrl, dorisDbKey))
}
