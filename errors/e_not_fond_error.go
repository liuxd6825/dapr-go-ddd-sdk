package errors

type NotFondError struct {
}

func NewNotFondError() error {
	return &NotFondError{}
}

func (e *NotFondError) Error() string {
	return "Not Fond"
}
