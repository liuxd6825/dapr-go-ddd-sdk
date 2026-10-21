package sys

import (
	"context"

	"github.com/kataras/iris/v12"
	bank_api "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/bank/restapi"
	card_api "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/restapi"
	case_api "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/case/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/configs"
	currency_api "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/currency/restapi"
	dict_api "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/dict/restapi"
	protal_api "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/restapi"
	status_api "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/status/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *envs.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "sys"}, func() error {
		sysconfigs.Init(env)
		cfg := sysconfigs.GetConfig()
		case_api.RegisterAllApi(app, baseUrl, env, cfg)
		currency_api.RegisterAllApi(app, baseUrl, env, cfg)
		bank_api.RegisterAllApi(app, baseUrl, env, cfg)
		card_api.RegisterAllApi(app, baseUrl, env, cfg)
		dict_api.RegisterAllApi(app, baseUrl, env, cfg)
		protal_api.RegisterAllApi(app, baseUrl, env, cfg)
		status_api.RegisterAllApi(app, baseUrl, env, cfg)

		return nil
	})
	if err != nil {
		panic(err)
	}
}
