package hserver

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
)

func catchError(ctx iris.Context, e error, recover any) error {
	var err error
	if e != nil {
		err = e
	} else if recover != nil {
		if ve, ok := recover.(error); ok {
			err = ve
		} else if obj, ok := recover.(*goja.Object); ok {
			eObj := obj.Export()
			err = fmt.Errorf("unknown error %s", eObj)
		} else {
			err = fmt.Errorf("unknown error %s", recover)
		}
	}
	if err != nil && ctx != nil {
		logs.Error(ctx, "", nil, err.Error())
		setError(ctx, err)
	}
	return nil
}

func setError(ctx iris.Context, err error) {
	if err != nil && ctx != nil {
		ctx.SetErr(err)
		ctx.StatusCode(httptest.StatusInternalServerError)
		_, _ = ctx.WriteString(err.Error())
	}
}
