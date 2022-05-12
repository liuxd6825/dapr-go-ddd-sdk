package ddd_repository

type FindOneResult[T any] struct {
	err     error
	data    T
	isFound bool
}

func NewFindOneResult[T any](data T, isFound bool, err error) *FindOneResult[T] {
	return &FindOneResult[T]{
		data:    data,
		isFound: isFound,
		err:     err,
	}
}

func (f *FindOneResult[T]) GetError() error {
	return f.err
}

func (f *FindOneResult[T]) GetData() T {
	return f.data
}

func (f *FindOneResult[T]) IsFound() bool {
	return f.isFound
}

func (f *FindOneResult[T]) Result() (T, bool, error) {
	return f.data, f.isFound, f.err
}

func (f *FindOneResult[T]) OnSuccess(success OnSuccess[T]) *FindOneResult[T] {
	if f.err == nil && success != nil && f.isFound {
		f.err = success(f.data)
	}
	return f
}

func (f *FindOneResult[T]) OnError(onErr OnError) *FindOneResult[T] {
	if f.err != nil && onErr != nil {
		f.err = onErr(f.err)
	}
	return f
}

func (f *FindOneResult[T]) OnNotFond(fond OnIsFond) *FindOneResult[T] {
	if f.err == nil && !f.isFound && fond != nil {
		f.err = fond()
	}
	return f
}
