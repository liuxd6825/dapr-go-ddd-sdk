package sys

import (
	"context"

	"github.com/kataras/iris/v12"
	card_api "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/restapi"
	case_api "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/case/restapi"
	currency_api "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/currency/restapi"
	dict_api "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/dict/restapi"
	protal_api "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "sys"}, func() error {
		case_api.RegisterAllApi(app, baseUrl, env)
		currency_api.RegisterAllApi(app, baseUrl, env)
		card_api.RegisterAllApi(app, baseUrl, env)
		dict_api.RegisterAllApi(app, baseUrl, env)
		protal_api.RegisterAllApi(app, baseUrl, env)
		return nil
	})
	if err != nil {
		panic(err)
	}
}
