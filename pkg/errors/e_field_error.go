package errors

import (
	"encoding/json"
)

type FieldError struct {
	Field   string `json:"field"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

func NewFieldError(fieldName string, message string, title ...string) *FieldError {
	err := &FieldError{
		Field:   fieldName,
		Message: message,
	}
	if len(title) > 0 {
		err.Title = title[0]
	}
	return err
}

func (e *FieldError) Error() string {
	bytes, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
