package errors

import (
	"errors"
	"fmt"
	"runtime"
)

func New(formatOrText string, text ...any) error {
	name := runFuncName(2)
	count := len(text)
	if count == 0 {
		return errors.New(fmt.Sprintf("%v() error: %v", name, formatOrText))
	}

	texts := make([]any, count+1)
	texts[0] = name
	for i, t := range text {
		texts[i+1] = t
	}
	if len(text) == 1 {
		return errors.New(fmt.Sprintf("%v() error:%v", name, text))
	}

	return errors.New(fmt.Sprintf("%v() error: "+formatOrText, texts...))

}

func NewMethod(packName, methodName string, msg string) error {
	return errors.New(fmt.Sprintf("%v.%v(), error:%v", packName, methodName, msg))
}

func ErrorOf(format string, args ...any) error {
	return New(format, args)
}

func runFuncName(skip int) string {
	pc := make([]uintptr, 1)
	runtime.Callers(skip+1, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
