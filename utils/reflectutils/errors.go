package reflectutils

import "fmt"

type methodNotExistError struct {
	typeName   string
	methodName string
}

func NewNewMethodNotExistError(typeName string, methodName string) error {
	return &methodNotExistError{
		typeName:   typeName,
		methodName: methodName,
	}
}

func (m *methodNotExistError) Error() string {
	return fmt.Sprintf("%s.%s()方法不存在.", m.typeName, m.methodName)
}

type methodCallError struct {
	typeName   string
	methodName string
	message    string
}

func NewMethodCallError(typeName, methodName, message string) error {
	return &methodCallError{
		typeName:   typeName,
		methodName: methodName,
		message:    message,
	}
}

func (e *methodCallError) Error() string {
	return fmt.Sprintf("%s.%s() doing error, %s.", e.typeName, e.methodName, e.message)
}
