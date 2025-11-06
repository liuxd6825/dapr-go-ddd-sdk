package gp

func Recover(err error, rec any) error {
	if err != nil {
		return err
	}
	if rec != nil {
		e, ok := rec.(error)
		if e != nil && ok {
			return e
		}
	}
	return nil
}
