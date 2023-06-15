package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type Errors struct {
	Details []error `json:"details"`
}

func (e *Errors) Error() string {
	sb := strings.Builder{}
	for _, er := range e.Details {
		if er != nil {
			sb.WriteString(er.Error())
			sb.WriteString(";")
		}
	}
	return sb.String()
}

func New(formatOrText string, text ...any) error {
	//funName := runFuncName(2)
	count := len(text)
	if count == 0 {
		return errors.New(formatOrText)
	}

	str := fmt.Sprintf(formatOrText, text...)
	return errors.New(str)

}

func News(errs ...error) error {
	return &Errors{
		Details: errs,
	}
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
