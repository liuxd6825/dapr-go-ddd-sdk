package errors

import (
	"encoding/json"
)

type VerifyError struct {
	Message string       `json:"message"`
	Errors  []FieldError `json:"errors"`
}

func NewVerifyError() *VerifyError {
	return &VerifyError{
		Message: "数据验证错误",
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

func (v *VerifyError) GetFieldError(fieldName string) *FieldError {
	for _, e := range v.Errors {
		if e.Field == fieldName {
			return &e
		}
	}
	return nil
}

func (v *VerifyError) GetFieldErrors() []FieldError {
	return v.Errors
}

func (v *VerifyError) AppendField(fieldName string, msg string, title ...string) {
	fieldError := NewFieldError(fieldName, msg, title...)
	v.Errors = append(v.Errors, *fieldError)
}

func (v *VerifyError) Count() int {
	return len(v.Errors)
}

func (v *VerifyError) IsHasError() bool {
	return len(v.Errors) > 0
}
func (v *VerifyError) HasError() bool {
	return v.IsHasError()
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

func (v *VerifyError) MergeWithPrefix(prefix string, ve *VerifyError) {
	if ve == nil {
		return
	}
	for _, e := range ve.Errors {
		newField := prefix + "." + e.Field
		v.Errors = append(v.Errors, FieldError{
			Field:   newField,
			Title:   e.Title,
			Message: e.Message,
		})
	}
}
