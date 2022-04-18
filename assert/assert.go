package assert

import (
	"strings"
)

//
// AssertError
// @Description:  断言错误
//
type AssertError struct {
	msg string
}

func NewAssertError(msg string) *AssertError {
	return &AssertError{msg: msg}
}

func (e *AssertError) Error() string {
	return e.msg
}

// Context
// @Description:  断言上下文对象
type Context map[string]*Options

//
// NewContext
// @Description: 创建断言上下文对象
// @return Context
//
func NewContext() Context {
	return Context{}
}

func (a Context) setValue(opt *Options, value bool) {
	if opt != nil {
		opt.value = value
	}
	if opt != nil && opt.HasKey() {
		a[*opt.GetKey()] = opt
	}
}

func (a Context) IsError() bool {
	for _, o := range a {
		if o.value == false {
			return false
		}
	}
	return true
}

func (a Context) Error() error {
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
	assert *Context
	key    *string
	error  *string
	value  bool
}

func WidthOptionsError(error string) *Options {
	return NewOptionsBuilder().SetError(error).Build()
}
func WidthOptionsKey(key string) *Options {
	return NewOptionsBuilder().SetKey(key).Build()
}

func (o *Options) Assert(assert *Context) {
	o.assert = assert
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

func (o *Options) Error() error {
	if !o.value {
		return NewAssertError(*o.error)
	}
	return nil
}

type OptionsBuilder struct {
	assert *Context
	key    *string
	error  *string
}

func NewOptionsBuilder() *OptionsBuilder {
	return &OptionsBuilder{}
}

func (b *OptionsBuilder) SetAssert(assert *Context) *OptionsBuilder {
	b.assert = assert
	return b
}

func (b *OptionsBuilder) SetKey(key string) *OptionsBuilder {
	b.key = &key
	return b
}

func (b *OptionsBuilder) SetError(error string) *OptionsBuilder {
	b.error = &error
	return b
}

func (b *OptionsBuilder) Build() *Options {
	return &Options{
		assert: b.assert,
		key:    b.key,
		error:  b.error,
	}
}

func mergeOptions(defaultError string, opts ...*Options) *Options {
	opt := &Options{}
	if opts != nil {
		for _, o := range opts {
			if o != nil {
				if o.assert != nil {
					opt.assert = o.assert
				}
				if o.key != nil {
					opt.key = o.key
				}
				if o.error != nil {
					opt.error = o.error
				}
			}
		}
	}
	if opt.assert == nil {
		opt.assert = &Context{}
	}
	if opt.error == nil {
		opt.error = &defaultError
	}
	return opt
}
