package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	docconfigs "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/configs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/subscribe"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func Init(app *iris.Application, baseUrl string, env *envs.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service-name ": "document"}, func() error {
		docconfigs.Init(env)
		if docconfigs.GetEnable() {
			dbKey := docconfigs.GetMongoDbKey()
			RegisterFolderApi(app, baseUrl, env, dbKey)
			RegisterFileApi(app, baseUrl, env, dbKey)
			RegisterDocumentApi(app, baseUrl, env, dbKey)
			RegisterDocumentMetaApi(app, baseUrl, env, dbKey)
			RegisterSub(app, baseUrl, env, dbKey)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterFolderApi(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewFolderAPI(env, baseUrl, dbKey))
}

func RegisterDocumentApi(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewDocumentAPI(env, baseUrl, docconfigs.GetConfig()))
}

func RegisterDocumentMetaApi(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewDocumentMetaAPI(env, baseUrl, dbKey))
}

func RegisterFileApi(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewFileAPI(env, baseUrl, dbKey))
}

func RegisterSub(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, subscribe.NewFolderEventSubHandler(env, baseUrl, dbKey))
}
