package web

import (
	"github.com/kataras/iris/v12"
)

func SetError(ctx iris.Context, err error) {
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.SetErr(err)
	}
}

func SetData(ctx iris.Context, data any) error {
	ctx.StatusCode(iris.StatusOK)
	return ctx.JSON(data)
}
