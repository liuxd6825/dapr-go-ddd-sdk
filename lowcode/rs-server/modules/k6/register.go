package k6

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/server"
)

func NewModules(app *iris.Application, data map[string]any, cfg common.IEnvConfig) map[string]any {
	value := make(map[string]any)
	value["server"] = server.New(app, data, cfg)
	return value
}
