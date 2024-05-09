package common

type Result[T any] struct {
	Data  T     `json:"data"`
	Ok    bool  `json:"ok"`
	Error error `json:"error"`
}

func NewResult[T any](data T, err error) *Result[T] {
	return &Result[T]{Data: data, Error: err, Ok: err == nil}
}

func (r *Result[T]) DoSuccess(fun func(data T)) *Result[T] {
	if r.Ok {
		fun(r.Data)
	}
	return r
}

func (r *Result[T]) DoError(fun func(err error)) *Result[T] {
	if !r.Ok {
		fun(r.Error)
	}
	return r
}

func (r *Result[T]) GetResults() (T, bool, error) {
	return r.Data, r.Ok, r.Error
}
