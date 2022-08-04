package errors

import (
	"encoding/json"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewFieldError(fieldName string, message string) *FieldError {
	return &FieldError{
		Field:   fieldName,
		Message: message,
	}
}

func (e *FieldError) Error() string {
	bytes, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
