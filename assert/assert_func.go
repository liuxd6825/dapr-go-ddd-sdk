package assert

func Nil(v interface{}, opts ...*Options) error {
	options := mergeOptions("assert.Nil() error", opts...)
	b := v == nil
	setValue(options, b)
	return options.Error()
}

func NotNil(v interface{}, opts ...*Options) error {
	options := mergeOptions("assert.NotNil() error", opts...)
	b := v != nil
	setValue(options, b)
	return options.Error()
}

func Ture(b bool, opts ...*Options) error {
	options := mergeOptions("assert.Ture() error", opts...)
	setValue(options, b)
	return options.Error()
}

func False(v bool, opts ...*Options) error {
	options := mergeOptions("assert.False() error", opts...)
	b := !v
	setValue(options, b)
	return options.Error()
}

func Empty(v string, opts ...*Options) error {
	options := mergeOptions("assert.Empty() error", opts...)
	b := v == ""
	setValue(options, b)
	return options.Error()
}

func NotEmpty(v string, opts ...*Options) error {
	options := mergeOptions("assert.NotEmpty() error", opts...)
	b := v != ""
	setValue(options, b)
	return options.Error()
}

func Equal(v1 interface{}, v2 interface{}, opts ...*Options) error {
	options := mergeOptions("assert.Equal() error ", opts...)
	b := v1 == v2
	setValue(options, b)
	return options.Error()
}

func NotEqual(v1 interface{}, v2 interface{}, opts ...*Options) error {
	options := mergeOptions("assert.NotEqual() error ", opts...)
	b := v1 != v2
	setValue(options, b)
	return options.Error()
}

func setValue(options *Options, bool2 bool) {
	if options != nil && options.assert != nil {
		options.assert.setValue(options, bool2)
	}
}
