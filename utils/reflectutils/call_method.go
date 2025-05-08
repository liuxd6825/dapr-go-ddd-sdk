package reflectutils

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"reflect"
)

// CallMethod2
//
//	@Description: 动态调用对象方法
//	@param object 被调用的对象
//	@param methodName 要执行的方法名称
//	@param ps 参数值
//	@return err 错误
func CallMethod2(obj interface{}, methodName string, args ...interface{}) (err error) {
	if obj == nil {
		return errors.ErrorOf("CallMethod() object is null")
	}

	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()

	if err = checkArgsIsNil(args...); err != nil {
		return err
	}

	at := reflect.TypeOf(obj).Elem()
	a := reflect.ValueOf(obj)
	typeName := at.Name()
	method := a.MethodByName(methodName)

	if method.IsValid() {
		args := make([]reflect.Value, len(args))
		for i, p := range args {
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

func CallMethod(obj any, methodName string, args ...any) (err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()

	if obj == nil {
		return errors.ErrorOf("reflectutils.CallMethod() obj is null")
	}

	if err = checkArgsIsNil(args...); err != nil {
		return errors.ErrorOf("reflectutils.CallMethod() error:%s", err.Error())
	}

	refObj := NewRefObj(obj)
	method := refObj.Method(methodName)
	if !method.IsValid() {
		return NewNewMethodNotExistError(refObj.Type().Name(), methodName)
	}

	//取参数总量
	inCount := len(method.InTypes())
	//如果传入参数过少，退出
	if inCount > len(args) {
		return errors.ErrorOf("动态调用方法时，传入的参数个数不足。%s.%s()", refObj.Type().Name(), methodName)
	}
	//参方法参数个数，传入参数
	inArgs := make([]any, inCount)
	for i := 0; i < inCount; i++ {
		inArgs[i] = args[i]
	}

	//调用方法
	callResult, err := method.Call(inArgs...)
	if err != nil {
		return err
	}
	if callResult.IsError() {
		return callResult.Error
	}
	if err = getResultError(callResult.Result...); err != nil {
		return errors.NewMethodCallError(refObj.Type().Name(), methodName, err.Error())
	}
	return nil
}

func getResultError(values ...any) (err error) {
	for _, v := range values {
		switch v.(type) {
		case error:
			if e, ok := v.(error); ok {
				err = e
			}
			break
		}
	}
	return err
}

func checkArgsIsNil(args ...any) error {
	for i, v := range args {
		if v == nil {
			return errors.New("第%i参数不能为nil", i)
		}
	}
	return nil
}
