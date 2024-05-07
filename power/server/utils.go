package server

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func Do(ctx iris.Context, fun func(ctx iris.Context) error) {
	var err error
	defer func() {
		if e := recover(); e != nil {
			if err := e.(error); err != nil {
				setError(ctx, err)
			}
		}
	}()

	err = fun(ctx)
	if err != nil {
		setError(ctx, err)
	}
}

func setError(ctx iris.Context, err error) {
	ctx.SetErr(err)
	ctx.StatusCode(httptest.StatusInternalServerError)
	_, _ = ctx.WriteString(err.Error())
}

func ctxRecover(ctx iris.Context, recover any) {
	if recover != nil {
		if err := recover.(error); err != nil {
			setError(ctx, err)
		}
	}
}
