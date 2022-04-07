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

type FindPagingResult struct {
	FindResult
}

func NewFindPagingListResult(data *PagingData, isFound bool, err error) *FindPagingResult {
	return &FindPagingResult{
		FindResult{
			data:    data,
			isFound: isFound,
			err:     err,
		},
	}
}

func (f *FindPagingResult) OnSuccess(success OnSuccessPaging) *FindPagingResult {
	if f.err == nil && success != nil && f.isFound {
		data := f.data.(*PagingData)
		f.err = success(data)
	}
	return f
}
