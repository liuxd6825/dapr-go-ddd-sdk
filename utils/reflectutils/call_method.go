package reflectutils

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"reflect"
)

// CallMethod
//
//	@Description: 动态调用对象方法
//	@param object 被调用的对象
//	@param methodName 要执行的方法名称
//	@param ps 参数值
//	@return err 错误
func CallMethod(object interface{}, methodName string, ps ...interface{}) (err error) {
	at := reflect.TypeOf(object).Elem()
	a := reflect.ValueOf(object)
	typeName := at.Name()
	method := a.MethodByName(methodName)

	defer func() {
		if e := recover(); e != nil {
			message := fmt.Sprintf("%v", e)
			err = errors.NewMethodCallError(typeName, methodName, message)
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
				return errors.NewMethodCallError(at.Name(), methodName, err.Error())
			}
		}
		return nil
	}

	return errors.NewMethodNotExistError(typeName, methodName)
}
