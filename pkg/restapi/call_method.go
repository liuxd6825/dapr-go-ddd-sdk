package restapi

import (
	context2 "context"
	"github.com/kataras/iris/v12/context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"reflect"
)

type CallMethod struct {
	Method        reflect.Value
	InCtx         int
	InICtx        int
	InParams      int
	CloseInParams bool // 是否关闭InParams参数
	OutData       int
	OutError      int
}

func NewCallMethod(object any, methodName string) (*CallMethod, error) {
	method := reflect.ValueOf(object).MethodByName(methodName)
	if !method.IsValid() {
		return nil, errors.New("handler not found: %s ", methodName)
	}

	callMethod := CallMethod{
		Method:   method,
		InCtx:    -1,
		InICtx:   -1,
		InParams: -1,
		OutData:  -1,
		OutError: -1,
	}

	inCount := method.Type().NumIn()
	for i := 0; i < inCount; i++ {
		paramType := method.Type().In(i)
		if paramType == reflect.TypeOf((*context2.Context)(nil)).Elem() {
			callMethod.InCtx = i
		} else if paramType == reflect.TypeOf((*context.Context)(nil)).Elem() {
			callMethod.InICtx = i
		} else if paramType.Kind() == reflect.Ptr && paramType.Elem().Kind() == reflect.Struct {
			callMethod.InParams = i
		} else if paramType == reflect.TypeOf((*store.FindPagingQuery)(nil)).Elem() {
			// 允许any类型作为参数
			callMethod.InParams = i
		} else {
			return nil, errors.New("unsupported parameter type: %s", paramType.String())
		}
	}

	outCount := method.Type().NumOut()
	for i := 0; i < outCount; i++ {
		outType := method.Type().Out(i)
		if outType == reflect.TypeOf((*any)(nil)).Elem() {
			callMethod.OutData = i
		} else if outType == reflect.TypeOf((*error)(nil)).Elem() {
			callMethod.OutError = i
		}
	}

	return &callMethod, nil
}

func (c *CallMethod) Call(ctx context2.Context, ictx *context.Context, params any) (any, error) {
	in := make([]reflect.Value, c.Method.Type().NumIn())
	if c.InCtx >= 0 {
		in[c.InCtx] = reflect.ValueOf(ctx)
	}
	if c.InICtx >= 0 {
		in[c.InICtx] = reflect.ValueOf(ictx)
	}
	if c.InParams >= 0 {
		in[c.InParams] = reflect.ValueOf(params)
	}

	out := c.Method.Call(in)

	var resData any
	var err error
	if c.OutData >= 0 {
		resData = out[c.OutData].Interface()
	}
	if c.OutError >= 0 {
		errValue := out[c.OutError]
		if !errValue.IsNil() {
			err = errValue.Interface().(error)
		}
	}

	return resData, err
}
