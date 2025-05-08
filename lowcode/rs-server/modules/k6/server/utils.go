package server

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
	"reflect"
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

func NewObject(runtime *goja.Runtime, data any) (*goja.Object, error) {
	if data == nil {
		obj := runtime.NewObject()
		return obj, nil
	}
	var err error
	obj := runtime.NewObject()
	if m, ok := data.(map[string]any); ok {
		for k, v := range m {
			if m, ok := v.(map[string]any); ok {
				if o, e := NewObject(runtime, m); e != nil {
					return nil, err
				} else {
					err = obj.Set(k, o)
				}
				continue
			} else if m, ok := v.(common.Object); ok {
				if o, e := NewObject(runtime, m); e != nil {
					return nil, err
				} else {
					err = obj.Set(k, o)
				}
				continue
			}
			t := reflect.ValueOf(v)
			if t.Kind() == reflect.Struct {
				if m, err := maputils.NewMap(data); err != nil {
					return nil, err
				} else {
					obj, err = NewObject(runtime, m)
				}
			} else if err = obj.Set(k, v); err != nil {
				fmt.Println(err)
			}
		}
	} else {
		if m, err = maputils.NewMap(data); err != nil {
			return nil, err
		} else {
			obj, err = NewObject(runtime, m)
		}

	}
	return obj, err
}
