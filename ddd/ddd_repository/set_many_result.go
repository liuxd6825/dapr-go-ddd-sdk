package ddd_repository

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd"

type SetManyResult[T interface{}] struct {
	err  error
	data []T
}

func NewSetManyResult[T ddd.Entity](data []T, err error) *SetManyResult[T] {
	return &SetManyResult[T]{
		data: data,
		err:  err,
	}
}

func NewSetManyResultError[T ddd.Entity](err error) *SetManyResult[T] {
	return &SetManyResult[T]{
		err: err,
	}
}

func (s *SetManyResult[T]) GetError() error {
	return s.err
}

func (s *SetManyResult[T]) GetData() []T {
	return s.data
}

func (s *SetManyResult[T]) Result() ([]T, error) {
	return s.data, s.err
}

func (s *SetManyResult[T]) OnSuccess(success OnSuccessList[T]) *SetManyResult[T] {
	if s.err == nil && success != nil {
		s.err = success(s.data)
	}
	return s
}

func (s *SetManyResult[T]) OnError(err OnError) *SetManyResult[T] {
	if s.err != nil && err != nil {
		s.err = err(s.err)
	}
	return s
}
