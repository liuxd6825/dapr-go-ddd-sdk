package utils

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func RecoverError(e error, recover any) error {
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
	return err
}

func SetError(ctx iris.Context, err error) {
	if err != nil && ctx != nil {
		ctx.SetErr(err)
		ctx.StatusCode(httptest.StatusInternalServerError)
		_, _ = ctx.WriteString(err.Error())

	}
}
