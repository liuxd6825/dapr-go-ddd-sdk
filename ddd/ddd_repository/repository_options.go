package ddd_repository

import "time"

type FindOptions struct {
	MaxTime *time.Duration
}

type FindOneOptions struct {
	MaxTime *time.Duration
}

type SetOptions struct {
	MaxTime *time.Duration
}

func MergeFindOptions(opts ...*FindOptions) *FindOptions {
	res := &FindOptions{}
	for _, o := range opts {
		if o.MaxTime != nil {
			res.MaxTime = o.MaxTime
		}
	}
	return res
}

func MergeSetOptions(opts ...*SetOptions) *SetOptions {
	res := &SetOptions{}
	for _, o := range opts {
		if o.MaxTime != nil {
			res.MaxTime = o.MaxTime
		}
	}
	return res
}

/*
type FindOptions struct {
	Error       error
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

func (f *FindOptions) init(opts ...FindOption) *FindOptions {
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func (f *FindOptions) SetResult(data interface{}, isFind bool, Error error) *FindOptions {
	f.data = data
	f.isFind = isFind
	f.Error = Error
	return f
}

func (f *FindOptions) Error() error {
	return f.Error
}

func (f *FindOptions) Data() interface{} {
	return f.data
}

func (f *FindOptions) isFound() bool {
	return f.isFind
}
*/
