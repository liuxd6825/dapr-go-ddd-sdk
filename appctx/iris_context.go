package appctx

import (
	"context"
	"github.com/kataras/iris/v12"
)

type irisCtxKey struct {
}

var irisKey = irisCtxKey{}

func NewIrisContext(parent context.Context, ictx iris.Context) context.Context {
	return context.WithValue(parent, irisKey, ictx)
}

func GetIrisContext(ctx context.Context) (iris.Context, bool) {
	if ctx == nil {
		return nil, false
	}
	if ictx, ok := ctx.(iris.Context); ok {
		return ictx, ok
	}
	val := ctx.Value(irisKey)
	if val == nil {
		return nil, false
	}
	ictx, ok := val.(iris.Context)
	return ictx, ok
}