package errorutils

import "errors"

func GetMessage(e any) (res string, ok bool) {
	switch e.(type) {
	case string:
		ok = true
		res, _ = e.(string)
		break
	case *string:
		ok = true
		spur, _ := e.(*string)
		res = *spur
		break
	case error:
		ok = true
		err, _ := e.(error)
		res = err.Error()
	}
	return res, ok
}

func GetError(re any) (err error) {
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
