package ddd_repository

import "github.com/dapr/dapr-go-ddd-sdk/ddd"

type FindListResult[T ddd.Entity] struct {
	err     error
	data    *[]T
	isFound bool
}

func NewFindListResult[T ddd.Entity](data *[]T, isFound bool, err error) *FindListResult[T] {
	return &FindListResult[T]{
		data:    data,
		isFound: isFound,
		err:     err,
	}
}

func (f *FindListResult[T]) GetError() error {
	return f.err
}

func (f *FindListResult[T]) GetData() *[]T {
	return f.data
}

func (f *FindListResult[T]) IsFound() bool {
	return f.isFound
}

func (f *FindListResult[T]) Result() (*[]T, bool, error) {
	return f.data, f.isFound, f.err
}

func (f *FindListResult[T]) OnSuccess(success OnSuccessList[T]) *FindListResult[T] {
	if f.err == nil && success != nil && f.isFound {
		f.err = success(f.data)
	}
	return f
}

func (f *FindListResult[T]) OnError(onErr OnError) *FindListResult[T] {
	if f.err != nil && onErr != nil {
		f.err = onErr(f.err)
	}
	return f
}

func (f *FindListResult[T]) OnNotFond(fond OnIsFond) *FindListResult[T] {
	if f.err == nil && !f.isFound && fond != nil {
		f.err = fond()
	}
	return f
}
