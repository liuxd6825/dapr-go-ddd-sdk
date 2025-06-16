package web

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

func SetError(ctx iris.Context, err error) {
	if err != nil {
		ctx.SetErr(err)
		if _, ok := err.(*errors.VerifyError); ok {
			ctx.StatusCode(iris.StatusBadRequest)
		} else {
			ctx.StatusCode(iris.StatusInternalServerError)
		}
	}
}

func SetData(ctx iris.Context, data any) error {
	ctx.StatusCode(iris.StatusOK)
	return ctx.JSON(data)
}
