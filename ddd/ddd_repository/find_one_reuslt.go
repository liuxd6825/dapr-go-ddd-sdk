package ddd_repository

type FindOneResult[T any] struct {
	Err     error `json:"err" js:"err"`
	Data    T     `json:"data" js:"data"`
	IsFound bool  `json:"isFound" js:"isFound"`
}

func NewFindOneResult[T any](data T, isFound bool, err error) *FindOneResult[T] {
	return &FindOneResult[T]{
		Data:    data,
		IsFound: isFound,
		Err:     err,
	}
}

func (f *FindOneResult[T]) GetError() error {
	return f.Err
}

func (f *FindOneResult[T]) GetData() T {
	return f.Data
}

func (f *FindOneResult[T]) GetIsFound() bool {
	return f.IsFound
}

func (f *FindOneResult[T]) Result() (T, bool, error) {
	return f.Data, f.IsFound, f.Err
}

func (f *FindOneResult[T]) OnSuccess(success OnSuccess[T]) *FindOneResult[T] {
	if f.Err == nil && success != nil && f.IsFound {
		f.Err = success(f.Data)
	}
	return f
}

func (f *FindOneResult[T]) OnError(onErr OnError) *FindOneResult[T] {
	if f.Err != nil && onErr != nil {
		f.Err = onErr(f.Err)
	}
	return f
}

func (f *FindOneResult[T]) OnNotFond(fond OnIsFond) *FindOneResult[T] {
	if f.Err == nil && !f.IsFound && fond != nil {
		f.Err = fond()
	}
	return f
}
