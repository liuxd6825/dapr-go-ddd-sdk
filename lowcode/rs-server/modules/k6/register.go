package k6

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/server"
)

func NewModules(app *iris.Application) map[string]any {
	value := make(map[string]any)
	value["server"] = server.New(app)
	return value
}
