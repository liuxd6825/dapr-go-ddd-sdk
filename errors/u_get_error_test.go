package errors

import (
	"errors"
	"fmt"
	"testing"
)

func Test_GetRecoverError(t *testing.T) {
	err := getRecoverError()
	if err != nil {
		t.Error(err)
	}
}

func getRecoverError() (resErr error) {
	defer func() {
		resErr = GetRecoverError(resErr, recover())
		fmt.Print(resErr)
	}()

	err := errors.New("error")
	return err
}
