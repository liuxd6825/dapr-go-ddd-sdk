package ddd_repository

type FindOneResult[T any] struct {
	Error   error `json:"error" js:"Error"`
	Data    T     `json:"data" js:"data"`
	IsFound bool  `json:"isFound" js:"isFound"`
}

func NewFindOneResult[T any](data T, isFound bool, err error) *FindOneResult[T] {
	return &FindOneResult[T]{
		Data:    data,
		IsFound: isFound,
		Error:   err,
	}
}

func (f *FindOneResult[T]) GetError() error {
	return f.Error
}

func (f *FindOneResult[T]) GetData() T {
	return f.Data
}

func (f *FindOneResult[T]) GetIsFound() bool {
	return f.IsFound
}

func (f *FindOneResult[T]) Result() (T, bool, error) {
	return f.Data, f.IsFound, f.Error
}

func (f *FindOneResult[T]) OnSuccess(success OnSuccess[T]) *FindOneResult[T] {
	if f.Error == nil && success != nil && f.IsFound {
		f.Error = success(f.Data)
	}
	return f
}

func (f *FindOneResult[T]) OnError(onErr OnError) *FindOneResult[T] {
	if f.Error != nil && onErr != nil {
		f.Error = onErr(f.Error)
	}
	return f
}

func (f *FindOneResult[T]) OnNotFond(fond OnIsFond) *FindOneResult[T] {
	if f.Error == nil && !f.IsFound && fond != nil {
		f.Error = fond()
	}
	return f
}
