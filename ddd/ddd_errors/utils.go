package ddd_errors

import "errors"

func GetRecoverError() (err error) {
	if m := recover(); m != nil {
		if msg, ok := m.(string); !ok {
			err = errors.New(msg)
		} else {
			err = errors.New(" any error.")
		}
	}
	return
}
