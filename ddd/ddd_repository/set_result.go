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

func (s *SetResult) Result() (interface{}, error) {
	return s.data, s.err
}

func (s *SetResult) OnSuccess(success OnSuccess) *SetResult {
	if s.err == nil && success != nil {
		s.err = success(s.data)
	}
	return s
}

func (s *SetResult) OnError(err OnError) *SetResult {
	if s.err != nil && err != nil {
		s.err = err(s.err)
	}
	return s
}
