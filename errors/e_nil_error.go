package errors

type NullError struct {
}

func NewNullError() *NullError {
	return &NullError{}
}
func (e *NullError) Error() string {
	return "null error"
}
