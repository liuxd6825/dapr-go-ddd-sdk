package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

func SetError(ctx iris.Context, err error) {
	if err != nil {
		req := ctx.Request()
		logs.Error(ctx, logs.Fields{"method": req.Method, "uri": req.RequestURI, "error": err.Error()})
		ctx.SetErr(err)
		if _, ok := err.(*errors.VerifyError); ok {
			ctx.StatusCode(iris.StatusBadRequest)
		} else {
			ctx.StatusCode(iris.StatusInternalServerError)
		}
	}
}

func SetNotFoundError(ictx iris.Context, err error) {
	req := ictx.Request()
	if err != nil {
		logs.Error(ictx, logs.Fields{"method": req.Method, "uri": req.RequestURI, "error": err.Error()})
		ictx.SetErr(err)
	}
	ictx.StatusCode(iris.StatusFound)
}

func SetData(ctx iris.Context, data any) error {
	ctx.StatusCode(iris.StatusOK)
	return ctx.JSON(data)
}
