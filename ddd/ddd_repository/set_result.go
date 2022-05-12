package ddd_repository

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd"

type SetResult[T ddd.Entity] struct {
	err  error
	data T
}

func NewSetResult[T ddd.Entity](data T, err error) *SetResult[T] {
	return &SetResult[T]{
		data: data,
		err:  err,
	}
}

func (s *SetResult[T]) GetError() error {
	return s.err
}

func (s *SetResult[T]) GetData() T {
	return s.data
}

func (s *SetResult[T]) Result() (T, error) {
	return s.data, s.err
}

func (s *SetResult[T]) OnSuccess(success OnSuccess[T]) *SetResult[T] {
	if s.err == nil && success != nil {
		s.err = success(s.data)
	}
	return s
}

func (s *SetResult[T]) OnError(err OnError) *SetResult[T] {
	if s.err != nil && err != nil {
		s.err = err(s.err)
	}
	return s
}
