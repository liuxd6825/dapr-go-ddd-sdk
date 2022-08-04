package errors

import (
	"errors"
	"fmt"
)

func New(text string) error {
	return errors.New(text)
}

func ErrorOf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}
