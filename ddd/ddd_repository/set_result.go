package ddd_repository

type SetResult[T interface{}] struct {
	Data         T     `json:"data"`
	Error        error `json:"error"`
	RowsAffected int64 `json:"rowsAffected"`
}

func NewSetResult[T interface{}](data T, err error) *SetResult[T] {
	return &SetResult[T]{
		Data:         data,
		Error:        err,
		RowsAffected: 0,
	}
}

func NewSetResultError[T interface{}](err error) *SetResult[T] {
	return &SetResult[T]{
		Error: err,
	}
}

func (s *SetResult[T]) SetError(err error) *SetResult[T] {
	s.Error = err
	return s
}

func (s *SetResult[T]) GetError() error {
	return s.Error
}

func (s *SetResult[T]) GetData() T {
	return s.Data
}

func (s *SetResult[T]) Result() (T, error) {
	return s.Data, s.Error
}

func (s *SetResult[T]) OnSuccess(success OnSuccess[T]) *SetResult[T] {
	if s.Error == nil && success != nil {
		s.Error = success(s.Data)
	}
	return s
}

func (s *SetResult[T]) OnError(err OnError) *SetResult[T] {
	if s.Error != nil && err != nil {
		s.Error = err(s.Error)
	}
	return s
}

func (s *SetResult[T]) SetRowsAffected(count int64) *SetResult[T] {
	s.RowsAffected = count
	return s
}

func (s *SetResult[T]) GetRowsAffected() int64 {
	return s.RowsAffected
}
