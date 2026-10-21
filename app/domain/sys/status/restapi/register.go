package restapi

import (
	"github.com/kataras/iris/v12"
	sysconfigs "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/configs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *envs.Env, cfg sysconfigs.Config) {
	sysconfigs.Init(env)
	RegisterStatusApi(app, baseUrl, env, cfg)
}

func RegisterStatusApi(app *iris.Application, baseUrl string, env *envs.Env, cfg sysconfigs.Config) {
	restapi.RegisterController(app, NewStatusApi(env, baseUrl, cfg.MongoDbKey))
}
