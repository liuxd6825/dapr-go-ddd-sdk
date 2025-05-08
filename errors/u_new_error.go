package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type Error struct {
	msg string
	err error
}

type Errors struct {
	details []error `json:"details"`
}

type ParamsError struct {
	funName string
	Errors
}

func NewErrors() *Errors {
	return &Errors{
		details: make([]error, 0),
	}
}

func NewParamsError(funName string) *ParamsError {
	return &ParamsError{
		funName: funName,
	}
}

func News(errs ...error) error {
	var errsList []error
	for _, err := range errs {
		if err != nil {
			errsList = append(errsList, err)
		}
	}
	return &Errors{
		details: errsList,
	}
}

func (e *ParamsError) AddNil(msg string, data any) {
	if data == nil {
		e.AddError(msg, "不能为nil")
	} else if v, ok := data.(string); ok {
		if v == "" {
			e.AddError(msg, "不能为nil")
		}
	} else if v, ok := data.(bool); ok {
		if v {
			e.AddError(msg, "条件不满足")
		}
	}
}

func (e *ParamsError) AddEmpty(msg string, data string) {
	if data == "" {
		e.AddError(msg, "不能为空")
	}

}

func (e *ParamsError) AddBool(msg string, condition bool) {
	if condition {
		e.AddError(msg, "条件不满足")
	}
}

func (e *ParamsError) AddError(name string, reason string) {
	err := errors.New(fmt.Sprintf("参数 \"%s\" 传入错误: %s。", name, reason))
	e.details = append(e.details, err)
}

func NewMethod(packName, methodName string, msg string) error {
	return errors.New(fmt.Sprintf("%v.%v(), error:%v", packName, methodName, msg))
}

func ErrorOf(format string, args ...any) error {
	return New(format, args...)
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

func (e *Errors) HasError() bool {
	return len(e.details) > 0
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

func NewErr(err error, formatOrText string, text ...any) error {
	/*
		count := len(text)
		if count == 0 {
			return errors.New(formatOrText)
		}*/

	str := fmt.Sprintf(formatOrText, text...)
	return &Error{
		msg: str,
		err: err,
	}
}

func NewFunc(skip int, formatOrText string, text ...any) error {
	funName := runFuncName(skip + 1)
	return New("funName:"+funName+"() "+formatOrText, text...)
}

func runFuncName(skip int) string {
	pc := make([]uintptr, 1)
	runtime.Callers(skip+1, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func (e *Error) Error() string {
	if e.err == nil {
		return e.msg
	}
	return e.msg + ": " + e.err.Error()
}
