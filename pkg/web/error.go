package web

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

func SetError(ctx iris.Context, err error) {
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.SetErr(err)
		if errors.Is(err, &errors.VerifyError{}) {
			verr := err.(*errors.VerifyError)
			_ = ctx.JSON(verr)
		}
	}
}

func SetData(ctx iris.Context, data any) error {
	ctx.StatusCode(iris.StatusOK)
	return ctx.JSON(data)
}
