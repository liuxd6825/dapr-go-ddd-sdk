package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type Errors struct {
	details []error `json:"details"`
}

func NewErrors() *Errors {
	return &Errors{
		details: make([]error, 0),
	}
}

func (e *Errors) Error() string {
	sb := strings.Builder{}
	for _, er := range e.details {
		if er != nil {
			sb.WriteString(er.Error())
			sb.WriteString(";")
		}
	}
	return sb.String()
}

func (e *Errors) AddError(err error) {
	e.details = append(e.details, err)
}

func (e *Errors) AddString(str string) {
	e.AddError(errors.New(str))
}

func (e *Errors) AddFormat(format string, obj ...any) {
	e.AddError(errors.New(fmt.Sprintf(format, obj...)))
}

func (e *Errors) Len() int {
	return len(e.details)
}

func (e *Errors) IsEmpty() bool {
	return len(e.details) == 0
}

func (e *Errors) List() []error {
	return e.details
}

func (e *Errors) NewError() error {
	if len(e.details) == 0 {
		return nil
	}
	return errors.New(e.Error())
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

func NewFunc(skip int, formatOrText string, text ...any) error {
	funName := runFuncName(skip + 1)
	return New("funName:"+funName+"() "+formatOrText, text...)
}

func News(errs ...error) error {
	return &Errors{
		details: errs,
	}
}

func NewMethod(packName, methodName string, msg string) error {
	return errors.New(fmt.Sprintf("%v.%v(), error:%v", packName, methodName, msg))
}

func ErrorOf(format string, args ...any) error {
	return New(format, args...)
}

func runFuncName(skip int) string {
	pc := make([]uintptr, 1)
	runtime.Callers(skip+1, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
