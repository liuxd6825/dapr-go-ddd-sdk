package ddd

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"reflect"
)

// CallMethod
// @Description: 动态调用方法
// @param object 被调用的结构
// @param methodName 方法名称
// @param ps 参数
// @return err
func CallMethod(object interface{}, methodName string, ps ...interface{}) (err error) {
	if object == nil {
		return errors.NewMethodCallError("ddd", "CallMethod", "object is nil")
	}
	at := reflect.TypeOf(object).Elem()
	typeName := at.Name()

	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()

	a := reflect.ValueOf(object)
	method := a.MethodByName(methodName)

	if method.IsValid() {
		args := make([]reflect.Value, len(ps))
		for i, p := range ps {
			args[i] = reflect.ValueOf(p)
		}
		resValues := method.Call(args)
		for _, v := range resValues {
			var err1 error
			switch v.Interface().(type) {
			case error:
				if e, ok := v.Interface().(error); ok {
					err1 = e
				}
				break
			}
			if err1 != nil {
				return err1
			}
		}
		return nil
	}

	return errors.NewMethodNotExistError(typeName, methodName)
}
