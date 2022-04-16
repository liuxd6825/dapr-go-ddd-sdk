package assert

func Nil(a Assert, v interface{}, opts ...*Options) bool {
	options := mergeOptions(opts...)
	b := v == nil
	setValue(a, options, b)
	return b
}

func NotNil(a Assert, v interface{}, opts ...*Options) bool {
	options := mergeOptions(opts...)
	b := v != nil
	setValue(a, options, b)
	return b
}

func Equal(a Assert, v1 interface{}, v2 interface{}, opts ...*Options) bool {
	options := mergeOptions(opts...)
	b := v1 == v2
	setValue(a, options, b)
	return b
}

func NotEqual(a Assert, v1 interface{}, v2 interface{}, opts ...*Options) bool {
	options := mergeOptions(opts...)
	b := v1 != v2
	setValue(a, options, b)
	return b
}
