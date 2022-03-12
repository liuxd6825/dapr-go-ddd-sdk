package ddd_repository

type SetOptions struct {
	err       error
	data      interface{}
	OnSuccess OnSuccess
	OnError   OnError
}

type SetOption func(actions *SetOptions)

func SetOnSuccess(success OnSuccess) SetOption {
	return func(opt *SetOptions) {
		opt.OnSuccess = success
	}
}

func SetOnError(onError OnError) SetOption {
	return func(opt *SetOptions) {
		opt.OnError = onError
	}
}

func NewSetOptions() *SetOptions {
	return &SetOptions{
		OnSuccess: onSuccessDefault,
		OnError:   onErrorDefault,
	}
}

func (f *SetOptions) Init(opts ...SetOption) *SetOptions {
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func (f *SetOptions) Error() error {
	return f.err
}

func (f *SetOptions) Data() interface{} {
	return f.data
}

func (f *SetOptions) SetResult(data interface{}, err error) *SetOptions {
	f.data = data
	f.err = err
	return f
}
