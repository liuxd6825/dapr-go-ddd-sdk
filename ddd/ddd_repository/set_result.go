package ddd_repository

type SetResult struct {
	err  error
	data interface{}
}

func NewSetResult(data interface{}, err error) *SetResult {
	return &SetResult{
		data: data,
		err:  err,
	}
}

func (s *SetResult) GetError() error {
	return s.err
}

func (s *SetResult) GetData() interface{} {
	return s.data
}
