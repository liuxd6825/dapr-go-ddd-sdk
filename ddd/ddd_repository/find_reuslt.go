package ddd_repository

type FindResult struct {
	err     error
	data    interface{}
	isFound bool
}

func NewFindResult(data interface{}, isFound bool, err error) *FindResult {
	return &FindResult{
		data:    data,
		isFound: isFound,
		err:     err,
	}
}

func (f *FindResult) GetError() error {
	return f.err
}

func (f *FindResult) GetData() interface{} {
	return f.data
}

func (f *FindResult) IsFound() interface{} {
	return f.isFound
}

func (f *FindResult) Result() (interface{}, bool, error) {
	return f.data, f.isFound, f.err
}

func (f *FindResult) OnSuccess(success OnSuccess) *FindResult {
	if f.err == nil && success != nil && f.isFound {
		f.err = success(f.data)
	}
	return f
}

func (f *FindResult) OnError(onErr OnError) *FindResult {
	if f.err != nil && onErr != nil {
		f.err = onErr(f.err)
	}
	return f
}

func (f *FindResult) OnNotFond(fond OnIsFond) *FindResult {
	if f.err == nil && !f.isFound && fond != nil {
		f.err = fond()
	}
	return f
}
