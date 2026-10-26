package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/configs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/service"
	envs "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/spark/livyserver"
)

func Init(app *iris.Application, baseUrl string, env *envs.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service-name ": "import"}, func() error {
		importconfigs.Init(env)

		cfg := importconfigs.GetConfig()
		livyServer := livyserver.NewLivyServer(cfg.LivyServer.IP, cfg.LivyServer.Port)
		recordService := service.NewRecordService(livyServer, cfg.DorisDbKey)

		RegisterExcelRowApi(app, baseUrl, env)
		RegisterExcelFileApi(app, baseUrl, env)
		RegisterExcelSheetApi(app, baseUrl, env)
		RegisterTaskApi(app, baseUrl, env, livyServer)
		RegisterRecordAPI(app, baseUrl, env, recordService)
		RegisterTemplateAPI(app, baseUrl, env)
		RegisterRagAPI(app, baseUrl, env)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterExcelRowApi(app *iris.Application, baseUrl string, env *envs.Env) {
	api := NewExcelRowApi(env, baseUrl, importconfigs.GetConfig().DorisDbKey)
	restapi.RegisterController(app, api)
}

func RegisterExcelFileApi(app *iris.Application, baseUrl string, env *envs.Env) {
	api := NewExcelFileApi(env, baseUrl, importconfigs.GetConfig().DorisDbKey)
	restapi.RegisterController(app, api)
}

func RegisterExcelSheetApi(app *iris.Application, baseUrl string, env *envs.Env) {
	api := NewExcelSheetApi(env, baseUrl, importconfigs.GetConfig().DorisDbKey)
	restapi.RegisterController(app, api)
}

func RegisterTaskApi(app *iris.Application, baseUrl string, env *envs.Env, livyServer *livyserver.LivyServer) {
	api := NewTaskAPI(env, baseUrl, livyServer, importconfigs.GetConfig().DorisDbKey)
	restapi.RegisterController(app, api)
}

func RegisterRecordAPI(app *iris.Application, baseUrl string, env *envs.Env, recordService *service.RecordService) {
	api := NewRecordAPI(env, baseUrl, recordService)
	restapi.RegisterController(app, api)
}

func RegisterTemplateAPI(app *iris.Application, baseUrl string, env *envs.Env) {
	restapi.RegisterController(app, NewTemplateAPI(env, baseUrl, importconfigs.GetConfig().MongoDbKey))
}

func RegisterRagAPI(app *iris.Application, baseUrl string, env *envs.Env) {
	restapi.RegisterController(app, NewRagApi(env, baseUrl))
}
