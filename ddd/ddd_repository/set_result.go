package ddd_repository

type SetResult[T interface{}] struct {
	Error error `json:"error"`
	Data  T     `json:"data"`
}

func NewSetResult[T interface{}](data T, err error) *SetResult[T] {
	return &SetResult[T]{
		Data:  data,
		Error: err,
	}
}

func NewSetResultError[T interface{}](err error) *SetResult[T] {
	return &SetResult[T]{
		Error: err,
	}
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
