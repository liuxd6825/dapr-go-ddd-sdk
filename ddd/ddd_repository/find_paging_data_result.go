package ddd_repository

type PagingData struct {
	Data      interface{} `json:"data"`
	Count     int64       `json:"count"`
	Page      int64       `json:"page"`
	TotalPage int64       `json:"totalPage"`
	Size      int64       `json:"size"`
	Filter    string      `json:"filter"`
	Sort      string      `json:"sort"`
}

type FindPagingDataResult struct {
	FindResult
}

func NewFindPagingDataResult(data *PagingData, isFound bool, err error) *FindPagingDataResult {
	return &FindPagingDataResult{
		FindResult{
			data:    data,
			isFound: isFound,
			err:     err,
		},
	}
}

func (f *FindPagingDataResult) OnSuccess(success OnSuccessPaging) *FindPagingDataResult {
	if f.err == nil && success != nil && f.isFound {
		data := f.data.(*PagingData)
		f.err = success(data)
	}
	return f
}
