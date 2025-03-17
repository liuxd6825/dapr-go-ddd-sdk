package restapp

import (
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/mvc"
)

func Handle(b mvc.BeforeActivation, httpMethod, path, funcName string, middleware ...context.Handler) *router.Route {
	return b.Handle(httpMethod, path, funcName, middleware...)
}
