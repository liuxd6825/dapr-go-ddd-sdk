package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	ragconfigs "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/configs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/subscribe"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func Init(app *iris.Application, baseUrl string, env *envs.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service-name ": "master"}, func() error {
		ragconfigs.Init(env)
		RegisterRagApi(app, baseUrl, env)
		RegisterChatApi(app, baseUrl, env)
		RegisterMessageApi(app, baseUrl, env)
		RegisterDocumentApi(app, baseUrl, env)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterRagApi(app *iris.Application, baseUrl string, env *envs.Env) {
	restapi.RegisterController(app, NewRagAPI(env, baseUrl, ragconfigs.GetConfig().Rag, ragconfigs.GetMongoDbKey()))
}

func RegisterMessageApi(app *iris.Application, baseUrl string, env *envs.Env) {
	restapi.RegisterController(app, NewMessageAPI(env, baseUrl, ragconfigs.GetMongoDbKey()))
}

func RegisterChatApi(app *iris.Application, baseUrl string, env *envs.Env) {
	restapi.RegisterController(app, NewChatAPI(env, baseUrl, ragconfigs.GetMongoDbKey()))
}

func RegisterDocumentApi(app *iris.Application, baseUrl string, env *envs.Env) {
	restapi.RegisterController(app, NewDocumentAPI(env, baseUrl, ragconfigs.GetMongoDbKey()))
	restapi.RegisterController(app, subscribe.NewDocumentEventHandler(env, baseUrl, ragconfigs.GetMongoDbKey()))

}
