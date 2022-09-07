package errors

import (
	"errors"
	"fmt"
	"runtime"
)

func New(text string) error {
	name := runFuncName(2)
	return errors.New(fmt.Sprintf("%v() error:%v", name, text))
}

func NewMethod(packName, methodName string, msg string) error {
	return errors.New(fmt.Sprintf("%v.%v(), error:%v", packName, methodName, msg))
}

func ErrorOf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

func runFuncName(skip int) string {
	pc := make([]uintptr, 1)
	runtime.Callers(skip+1, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
