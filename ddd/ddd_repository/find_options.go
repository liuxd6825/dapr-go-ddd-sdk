package ddd_repository

type FindOptions struct {
	err       error
	data      interface{}
	isFind    bool
	OnSuccess OnSuccess
	OnError   OnError
	OnNotFond OnIsFond
}

type FindOption func(options *FindOptions)

func FindOnSuccess(success OnSuccess) FindOption {
	return func(opt *FindOptions) {
		opt.OnSuccess = success
	}
}

func FindOnError(onError OnError) FindOption {
	return func(opt *FindOptions) {
		opt.OnError = onError
	}
}

func NewFindOptions() *FindOptions {
	return &FindOptions{
		OnSuccess: onSuccessDefault,
		OnError:   onErrorDefault,
		OnNotFond: onNotFondDefault,
	}
}

func (f *FindOptions) Init(opts ...FindOption) *FindOptions {
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func (f *FindOptions) SetResult(data interface{}, isFind bool, err error) *FindOptions {
	f.data = data
	f.isFind = isFind
	f.err = err
	return f
}

func (f *FindOptions) Error() error {
	return f.err
}

func (f *FindOptions) Data() interface{} {
	return f.data
}

func (f *FindOptions) IsFound() bool {
	return f.isFind
}
