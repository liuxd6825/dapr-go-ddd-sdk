package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/subscribe"
	masterconfigs "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/configs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/subscribe/sub_import"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func Init(app *iris.Application, baseUrl string, env *envs.Env) {
	masterconfigs.Init(env)
	dbKey := masterconfigs.GetMongoDbKey()

	RegisterSchema(app, baseUrl, env)
	RegisterHtml(app, baseUrl, env)
	RegisterCdcToNeo4j(app, baseUrl, env)
	RegisterSub(app, baseUrl, env)
	RegisterRecord(app, baseUrl, env, dbKey)
	RegisterRecordDay(app, baseUrl, env, dbKey)
}

func RegisterSchema(app *iris.Application, baseUrl string, env *envs.Env) {
	schemaAPI := NewSchemaAPI(env, baseUrl)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(schemaAPI)
	})
}

func RegisterHtml(app *iris.Application, baseUrl string, env *envs.Env) {
	htmlAPI := NewHtmlAPI(env, "web")
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(htmlAPI)
	})
}

func RegisterRecord(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewRecordAPI(baseUrl, dbKey))
}

func RegisterRecordDay(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewRecordDayViewApi(baseUrl, dbKey))
}

func RegisterCdcToNeo4j(app *iris.Application, baseUrl string, env *envs.Env) {
	//restapi.RegisterController(app, NewGraphAPI(env, baseUrl))
	restapi.RegisterController(app, subscribe.NewCdcAPI(env, baseUrl, masterconfigs.GetConfig().MongoDbKey))
}

func RegisterSub(app *iris.Application, baseUrl string, env *envs.Env) {
	restapi.RegisterController(app, sub_import.NewRecordEventHandler(env, baseUrl))
}
