package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	tagconfigs "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/configs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func Init(app *iris.Application, baseUrl string, env *envs.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service-name ": "tag"}, func() error {
		tagconfigs.Init(env)
		RegisterTagApi(app, baseUrl, env)
		RegisterTagTypeApi(app, baseUrl, env)
		RegisterTagRelationApi(app, baseUrl, env)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterTagApi(app *iris.Application, baseUrl string, env *envs.Env) {
	restapi.RegisterController(app, NewTagAPI(env, baseUrl, tagconfigs.GetMongoDbKey()))
}

func RegisterTagTypeApi(app *iris.Application, baseUrl string, env *envs.Env) {
	restapi.RegisterController(app, NewTagTypeAPI(env, baseUrl, tagconfigs.GetMongoDbKey()))
}

func RegisterTagRelationApi(app *iris.Application, baseUrl string, env *envs.Env) {
	restapi.RegisterController(app, NewTagRelationAPI(env, baseUrl, tagconfigs.GetMongoDbKey()))
}
