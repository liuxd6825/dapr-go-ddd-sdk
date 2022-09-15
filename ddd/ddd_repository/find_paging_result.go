package ddd_repository

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
)

type FindPagingResult[T ddd.Entity] struct {
	Data        []T    `json:"data"`
	TotalRows   *int64 `json:"totalRows"`
	TotalPages  *int64 `json:"totalPages"`
	PageNum     int64  `json:"pageNum"`
	PageSize    int64  `json:"pageSize"`
	Filter      string `json:"filter"`
	Fields      string `json:"fields"`
	Sort        string `json:"sort"`
	IsFound     bool   `json:"isFound"`
	IsTotalRows bool   `json:"isTotalRows"`
	Error       error  `json:"-"`
}

type FindPagingResultOptions[T ddd.Entity] struct {
	Data        *[]T
	TotalRows   int64
	TotalPages  int64
	PageNum     int64
	PageSize    int64
	Filter      string
	Fields      string
	Sort        string
	IsFound     bool
	IsTotalRows bool
	Error       error
}

func NewFindPagingResult[T ddd.Entity](data []T, totalRows *int64, query FindPagingQuery, err error) *FindPagingResult[T] {
	res := &FindPagingResult[T]{
		Data:        data,
		TotalRows:   nil,
		TotalPages:  nil,
		PageNum:     0,
		PageSize:    0,
		Sort:        "",
		Filter:      "",
		IsFound:     false,
		IsTotalRows: false,
		Error:       err,
	}
	if data != nil {
		res.Data = data
		res.IsFound = len(data) > 0
	}
	if totalRows != nil {
		res.TotalRows = totalRows
		res.TotalPages = getTotalPage(totalRows, query.GetPageSize())
	}
	if query != nil {
		res.PageNum = query.GetPageNum()
		res.PageSize = query.GetPageSize()
		res.Sort = query.GetSort()
		res.Filter = query.GetFilter()
		res.Fields = query.GetFields()
		res.IsTotalRows = query.GetIsTotalRows()
	}
	return res
}

func NewFindPagingResultOptions[T ddd.Entity]() *FindPagingResultOptions[T] {
	return &FindPagingResultOptions[T]{}
}

func NewFindPagingResultWithError[T ddd.Entity](err error) *FindPagingResult[T] {
	return &FindPagingResult[T]{
		Data:    nil,
		IsFound: false,
		Error:   err,
	}
}

func (f *FindPagingResult[T]) GetData() []T {
	return f.Data
}

func (f *FindPagingResult[T]) GetAnyData() any {
	return f.Data
}

func (f *FindPagingResult[T]) GetTotalRows() *int64 {
	return f.TotalRows
}

func (f *FindPagingResult[T]) GetTotalPages() *int64 {
	return f.TotalPages
}

func (f *FindPagingResult[T]) GetPageNum() int64 {
	return f.PageNum
}

func (f *FindPagingResult[T]) GetPageSize() int64 {
	return f.PageSize
}

func (f *FindPagingResult[T]) GetFilter() string {
	return f.Filter
}

func (f *FindPagingResult[T]) GetFields() string {
	return f.Fields
}

func (f *FindPagingResult[T]) GetSort() string {
	return f.Sort
}

func (f *FindPagingResult[T]) GetIsFound() bool {
	return f.IsFound
}

func (f *FindPagingResult[T]) GetIsTotalRows() bool {
	return f.IsTotalRows
}

func (f *FindPagingResult[T]) GetError() error {
	return f.Error
}

func (f *FindPagingResult[T]) Result() (*FindPagingResult[T], bool, error) {
	return f, f.IsFound, f.Error
}

func (f *FindPagingResult[T]) OnError(onErr OnError) *FindPagingResult[T] {
	if f.Error != nil && onErr != nil {
		f.Error = onErr(f.Error)
	}
	return f
}

func (f *FindPagingResult[T]) OnNotFond(fond OnIsFond) *FindPagingResult[T] {
	if f.Error == nil && !f.IsFound && fond != nil {
		f.Error = fond()
	}
	return f
}

func (f *FindPagingResult[T]) OnSuccess(success OnSuccessList[T]) *FindPagingResult[T] {
	if f.Error == nil && success != nil && f.IsFound {
		f.Error = success(f.Data)
	}
	return f
}

func (f *FindPagingResultOptions[T]) SetData(data *[]T) *FindPagingResultOptions[T] {
	f.Data = data
	return f
}

func (f *FindPagingResultOptions[T]) SetTotalRows(totalRows int64) *FindPagingResultOptions[T] {
	f.TotalRows = totalRows
	return f
}

func (f *FindPagingResultOptions[T]) SetTotalPages(totalPages int64) *FindPagingResultOptions[T] {
	f.TotalPages = totalPages
	return f
}

func (f *FindPagingResultOptions[T]) SetPageNum(pageNum int64) *FindPagingResultOptions[T] {
	f.PageNum = pageNum
	return f
}

func (f *FindPagingResultOptions[T]) SetPageSize(pageSize int64) *FindPagingResultOptions[T] {
	f.PageSize = pageSize
	return f
}

func (f *FindPagingResultOptions[T]) SetFilter(filter string) *FindPagingResultOptions[T] {
	f.Filter = filter
	return f
}

func (f *FindPagingResultOptions[T]) SetSort(sort string) *FindPagingResultOptions[T] {
	f.Sort = sort
	return f
}

func (f *FindPagingResultOptions[T]) SetError(err error) *FindPagingResultOptions[T] {
	f.Error = err
	return f
}

func (f *FindPagingResultOptions[T]) SetIsFound(isFound bool) *FindPagingResultOptions[T] {
	f.IsFound = isFound
	return f
}

func (f *FindPagingResultOptions[T]) SetIsTotalRows(v bool) *FindPagingResultOptions[T] {
	f.IsTotalRows = v
	return f
}

func getTotalPage(totalRows *int64, pageSize int64) *int64 {
	if totalRows == nil {
		return nil
	}
	if pageSize == 0 {
		return nil
	}
	rows := *totalRows
	totalPage := rows / pageSize
	if rows%pageSize > 1 {
		totalPage++
	}
	return &totalPage
}
