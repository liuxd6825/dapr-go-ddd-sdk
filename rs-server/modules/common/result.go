package common

type Result[T any] struct {
	Data  T     `json:"data"`
	Error error `json:"error"`
}

func NewResult[T any](data T, err error) *Result[T] {
	return &Result[T]{Data: data, Error: err}
}

func (r *Result[T]) DoSuccess(fun func(data T) error) *Result[T] {
	if r.Error == nil {
		err := fun(r.Data)
		if err != nil {
			r.Error = err
		}
	}
	return r
}

func (r *Result[T]) DoError(fun func(err error)) *Result[T] {
	if r.Error != nil {
		fun(r.Error)
	}
	return r
}

func (r *Result[T]) GetResults() (T, error) {
	return r.Data, r.Error
}
