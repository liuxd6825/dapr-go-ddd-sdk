package ddd

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_errors"
	"reflect"
)

func CallMethod(object interface{}, methodName string, ps ...interface{}) (err error) {
	at := reflect.TypeOf(object).Elem()
	a := reflect.ValueOf(object)
	typeName := at.Name()
	method := a.MethodByName(methodName)

	defer func() {
		if e := recover(); e != nil {
			message := fmt.Sprintf("%v", e)
			err = ddd_errors.NewMethodCallError(typeName, methodName, message)
		}
	}()

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
				return ddd_errors.NewMethodCallError(at.Name(), methodName, err.Error())
			}
		}
		return nil
	}

	return ddd_errors.NewMethodNotExistError(typeName, methodName)
}
