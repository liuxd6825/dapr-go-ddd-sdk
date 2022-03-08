package ddd_errors

import (
	"encoding/json"
)

type VerifyError struct {
	Message string       `json:"message"`
	Errors  []FieldError `json:"errors"`
}

func NewVerifyError() *VerifyError {
	return &VerifyError{
		Message: "数据认证错误",
		Errors:  make([]FieldError, 0),
	}
}

func (v *VerifyError) Appends(errs *VerifyError) {
	if errs == nil {
		return
	}
	for _, e := range errs.Errors {
		v.Errors = append(v.Errors, e)
	}
}

func (v *VerifyError) AppendField(fieldName string, msg string) {
	fieldError := NewFieldError(fieldName, msg)
	v.Errors = append(v.Errors, *fieldError)
}

func (v *VerifyError) Count() int {
	return len(v.Errors)
}

func (v *VerifyError) IsHasError() bool {
	return len(v.Errors) > 0
}

func (v *VerifyError) Error() string {
	b, _ := json.Marshal(v)
	return string(b)
}

func (v *VerifyError) GetError() error {
	if v.IsHasError() {
		return v
	}
	return nil
}
