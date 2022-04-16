package ddd_errors

import "errors"

func GetRecoverError(re any) (err error) {
	err = nil
	if re != nil {
		switch re.(type) {
		case string:
			{
				msg, _ := re.(string)
				err = errors.New(msg)
			}
		case error:
			{
				e, _ := re.(error)
				err = e
			}
		}
	}
	return
}
