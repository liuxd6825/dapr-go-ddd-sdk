package ddd

import (
	"fmt"
	"reflect"
)

type MethodNotExistError struct {
	methodName string
}

func NewMethodError(methodName string) *MethodNotExistError {
	return &MethodNotExistError{
		methodName: methodName,
	}
}
func (e *MethodNotExistError) Error() string {
	return fmt.Sprintf(" %s() Method does not exist", e.methodName)
}

type MethodCallError struct {
	methodName string
	message    string
}

func NewMethodCallError(methodName string, message string) *MethodCallError {
	return &MethodCallError{
		methodName: methodName,
		message:    message,
	}
}
func (e *MethodCallError) Error() string {
	return fmt.Sprintf("reflect.Method.Call()  %s doing erro, %s", e.methodName, e.message)
}

func CallMethod(object interface{}, methodName string, ps ...interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			message := fmt.Sprintf("%v", e)
			err = NewMethodCallError(methodName, message)
		}
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
			if err, ok := v.Interface().(error); ok {
				return err
			}
		}
		return err
	}
	return NewMethodError(methodName)
}
