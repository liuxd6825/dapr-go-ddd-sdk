package goplus

import "github.com/liuxd6825/dapr-go-ddd-sdk/errors"

// const string with rethrow message
const gotry_rethrow = "----> Founded an Exception!!!\n"

// Exception interface
type Exception error

// GoTry object
type GoTry struct {
	catch   func(Exception)
	finally func()
	Error   Exception
}

// Throw function (return or rethrow an exception)
func Throw(e Exception) {
	if e == nil {
		panic(gotry_rethrow)
	} else {
		panic(e)
	}
}

// Try this function
func Try(funcToTry func() error) (o *GoTry) {
	o = &GoTry{nil, nil, nil}
	// catch throw in try
	defer func() {
		if o.Error == nil {
			o.Error = errors.GetRecoverError(o.Error, recover())
		}
	}()
	// do the func
	o.Error = funcToTry()
	return
}

// Catch function
func (o *GoTry) Catch(funcCatched func(err Exception)) *GoTry {
	o.catch = funcCatched
	if o.Error != nil {
		defer func() {
			// call finally
			if o.finally != nil {
				o.finally()
			}

			if err := recover(); err != nil {
				if err == gotry_rethrow {
					err = o.Error
				}
				panic(err)
			}
		}()
		o.catch(o.Error)
	} else if o.finally != nil {
		o.finally()
	}
	return o
}

// Finally function
func (o *GoTry) Finally(finallyFunc func()) {
	if o.finally != nil {
		panic("Finally Function by default !!")
	} else {
		o.finally = finallyFunc
	}
	defer o.finally()
	return
}
