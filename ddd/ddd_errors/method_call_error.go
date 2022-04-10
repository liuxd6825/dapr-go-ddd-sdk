package ddd_errors

import "fmt"

type MethodCallError struct {
	typeName   string
	methodName string
	message    string
}

func NewMethodCallError(typeName, methodName, message string) *MethodCallError {
	return &MethodCallError{
		typeName:   typeName,
		methodName: methodName,
		message:    message,
	}
}

func (e *MethodCallError) Error() string {
	return fmt.Sprintf("%s.%s() doing error, %s.", e.typeName, e.methodName, e.message)
}
