package rest

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_errors"
	"net/http"
)

const (
	ContentTypeApplicationJson = "application/json"
	ContentTypeTextPlain       = "text/plain"
)

type GetOneFunc = func(ctx iris.Context) (interface{}, bool, error)
type GetFunc = func(ctx iris.Context) (interface{}, error)

func ErrorNotFond(ctx iris.Context) {
	ctx.SetErr(iris.ErrNotFound)
	ctx.StatusCode(http.StatusNotFound)
	ctx.ContentType(ContentTypeTextPlain)
}

func ErrorInternalServerError(ctx iris.Context, err error) {
	ctx.SetErr(err)
	ctx.StatusCode(http.StatusInternalServerError)
	ctx.ContentType(ContentTypeTextPlain)
}

func ErrorVerifyError(ctx iris.Context, err error) {
	bytes, _ := json.Marshal(err)
	_, _ = ctx.Write(bytes)
	ctx.StatusCode(http.StatusInternalServerError)
	ctx.ContentType(ContentTypeTextPlain)
}

func SetError(ctx iris.Context, err error) {
	switch err.(type) {
	case *ddd_errors.NullError:
		ErrorNotFond(ctx)
		break
	case *ddd_errors.VerifyError:
		ErrorVerifyError(ctx, err)
		break
	default:
		ErrorInternalServerError(ctx, err)
		break
	}
}

func ResultOneData(ctx iris.Context, getAction GetOneFunc) {
	data, ok, err := getAction(ctx)
	if err != nil {
		SetError(ctx, err)
		return
	}
	if !ok {
		ErrorNotFond(ctx)
		return
	}
	if data == nil {
		ErrorNotFond(ctx)
		return
	}
	SetData(ctx, data)
}

func Result(ctx iris.Context, getAction GetFunc) {
	data, err := getAction(ctx)
	if err != nil {
		SetError(ctx, err)
		return
	}
	SetData(ctx, data)
}

func SetData(ctx iris.Context, data interface{}) {
	_, err := ctx.JSON(data)
	if err != nil {
		SetError(ctx, err)
		return
	}
}
