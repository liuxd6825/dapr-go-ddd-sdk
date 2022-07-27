package types

type Result[T any] struct {
	data  T
	error error
}

func NewResult[T any](data T, err error) *Result[T] {
	return &Result[T]{
		data:  data,
		error: err,
	}
}

func (r *Result[T]) OnSuccess(doSuccess func(data T) error) *Result[T] {
	if r.error == nil {
		r.error = doSuccess(r.data)
	}
	return r
}

func (r *Result[T]) OnError(doError func(err error)) *Result[T] {
	if r.error != nil {
		doError(r.error)
	}
	return r
}

func (r *Result[T]) Data() T {
	return r.data
}

func (r *Result[T]) Error() error {
	return r.error
}
