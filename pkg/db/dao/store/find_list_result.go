package store

type FindListResult[T interface{}] struct {
	Error   error `json:"error"`
	Data    []T   `json:"data"`
	IsFound bool  `json:"isFound"`
}

func NewFindListResult[T interface{}](data []T, isFound bool, err error) *FindListResult[T] {
	return &FindListResult[T]{
		Data:    data,
		IsFound: isFound,
		Error:   err,
	}
}

func NewFindListResultEmpty[T interface{}]() *FindListResult[T] {
	return &FindListResult[T]{}
}

func NewFindListResultError[T any](err error) *FindListResult[T] {
	return &FindListResult[T]{
		Error: err,
	}
}

func (f *FindListResult[T]) GetError() error {
	return f.Error
}

func (f *FindListResult[T]) SetError(err error) *FindListResult[T] {
	f.Error = err
	return f
}

func (f *FindListResult[T]) GetData() []T {
	return f.Data
}

func (f *FindListResult[T]) SetData(list []T) *FindListResult[T] {
	f.Data = list
	if len(list) > 0 {
		f.IsFound = true
	}
	return f
}

func (f *FindListResult[T]) GetIsFound() bool {
	return f.IsFound
}

func (f *FindListResult[T]) SetIsFound(val bool) *FindListResult[T] {
	f.IsFound = val
	return f
}

func (f *FindListResult[T]) Result() ([]T, bool, error) {
	return f.Data, f.IsFound, f.Error
}

func (f *FindListResult[T]) OnSuccess(success OnSuccessList[T]) *FindListResult[T] {
	if f.Error == nil && success != nil && f.IsFound {
		f.Error = success(f.Data)
	}
	return f
}

func (f *FindListResult[T]) OnError(onErr OnError) *FindListResult[T] {
	if f.Error != nil && onErr != nil {
		f.Error = onErr(f.Error)
	}
	return f
}

func (f *FindListResult[T]) OnNotFond(fond OnIsFond) *FindListResult[T] {
	if f.Error == nil && !f.IsFound && fond != nil {
		f.Error = fond()
	}
	return f
}
