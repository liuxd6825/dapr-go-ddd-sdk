package restapi

import (
	errors2 "errors"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"strings"
)

var notFoundError = errors2.New("Not Found")

func NotFoundError() error {
	return notFoundError
}

func SetError(ctx iris.Context, err error) {
	if err != nil {
		req := ctx.Request()
		logs.Error(ctx, logs.Fields{"method": req.Method, "uri": req.RequestURI, "error": err.Error()})
		ctx.SetErr(err)
		if _, ok := err.(*errors.VerifyError); ok {
			ctx.StatusCode(iris.StatusBadRequest)
		} else if isNotFoundError(err) {
			ctx.StatusCode(iris.StatusNotFound)
		} else if err.Error() == "" {
			ctx.StatusCode(iris.StatusInternalServerError)
		}
	}
}

func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	if err == notFoundError {
		return true
	}
	if strings.Contains(strings.ToLower(err.Error()), "not found") {
		return true
	}
	return false
}

func SetData2(ictx iris.Context, data any) error {
	ictx.StatusCode(iris.StatusOK)
	return ictx.JSON(data)
}

func SetData(ictx iris.Context, data any) error {
	return SetData2(ictx, data)
	jsonData, err := jsonutils.MarshalBytes(data)
	if err != nil {
		return err
	}
	ictx.ContentType("application/json")
	_, err = ictx.Write(jsonData)
	return err
}
