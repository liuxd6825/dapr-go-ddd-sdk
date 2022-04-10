package ddd_errors

import "fmt"

type MethodNotExistError struct {
	typeName   string
	methodName string
}

func NewMethodNotExistError(typeName string, methodName string) *MethodNotExistError {
	return &MethodNotExistError{
		typeName:   typeName,
		methodName: methodName,
	}
}

func (e *MethodNotExistError) Error() string {
	return fmt.Sprintf(" %s.%s() Method does not exist", e.typeName, e.methodName)
}
