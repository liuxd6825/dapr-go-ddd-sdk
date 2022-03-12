package ddd_repository

type FindResult struct {
	err    error
	data   interface{}
	isFind bool
}

func NewFindResult(data interface{}, isFind bool, err error) *FindResult {
	return &FindResult{
		data:   data,
		isFind: isFind,
		err:    err,
	}
}

func (f *FindResult) GetError() error {
	return f.err
}

func (f *FindResult) GetData() interface{} {
	return f.data
}

func (f *FindResult) GetIsFind() interface{} {
	return f.isFind
}
