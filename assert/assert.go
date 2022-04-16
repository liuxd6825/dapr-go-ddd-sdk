package assert

import "strings"

type Assert map[string]*Options

func NewAssert() Assert {
	return Assert{}
}

func (a Assert) setValue(opt *Options, value bool) {
	if opt != nil && opt.HasKey() {
		a[*opt.GetKey()] = opt
	}
}

func (a Assert) IsError() bool {
	for _, o := range a {
		if o.value == false {
			return false
		}
	}
	return true
}

func (a Assert) Error() error {
	err := &Error{}
	for _, item := range a {
		if !item.value {
			err.Add(*item.GetError())
		}
	}
	if err.hasError {
		return err
	}
	return nil
}

type Error struct {
	hasError bool
	build    strings.Builder
}

func (e *Error) Add(msg string) {
	e.hasError = true
	e.build.WriteString(msg)
}

func (e *Error) Error() string {
	return e.build.String()
}

type Options struct {
	key   *string
	error *string
	value bool
}

func (o *Options) HasKey() bool {
	return !(o.key == nil)
}

func (o *Options) GetKey() *string {
	return o.key
}

func (o *Options) SetKey(v string) {
	o.key = &v
}

func (o *Options) GetError() *string {
	return o.error
}

func setValue(a Assert, options *Options, bool2 bool) {
	if a != nil {
		a.setValue(options, bool2)
	}
}

func mergeOptions(opts ...*Options) *Options {
	opt := &Options{}
	for _, o := range opts {
		if o.key != nil {
			opt.key = o.key
		}
		if o.error != nil {
			opt.error = o.error
		}
	}
	return opt
}
