package ddd_repository

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
)

// FindPagingResultDTO 分页查询结果
type FindPagingResultDTO struct {
	TotalRows   *int64 `json:"totalRows,omitempty"`  // 总记录数
	TotalPages  *int64 `json:"totalPages,omitempty"` // 总页数
	PageNum     int64  `json:"pageNum"`              // 当前页号
	PageSize    int64  `json:"pageSize"`             // 页大小
	Filter      string `json:"filter"`               // RSQL过滤条件
	Fields      string `json:"fields"`               // 字段值，多个用逗号分隔
	Sort        string `json:"sort"`                 // 排序条件
	Error       error  `json:"error"`                // 错误
	IsFound     bool   `json:"isFound"`              // 是否找到数据
	IsTotalRows bool   `json:"isTotalRows"`          // 是否统计总记录数
}

type FindPagingResult[T any] struct {
	Data        []T    `json:"data"`
	SumData     []T    `json:"sumData"`
	TotalRows   *int64 `json:"totalRows"`
	TotalPages  *int64 `json:"totalPages"`
	PageNum     int64  `json:"pageNum"`
	PageSize    int64  `json:"pageSize"`
	Filter      string `json:"filter"`
	Fields      string `json:"fields"`
	Sort        string `json:"sort"`
	IsFound     bool   `json:"isFound"`
	IsTotalRows bool   `json:"isTotalRows"`
	IsSum       bool   `json:"isSum"`
	Error       error  `json:"-"`
}

type FindPagingResultOptions[T interface{}] struct {
	Data        *[]T
	SumData     *[]T
	TotalRows   int64
	TotalPages  int64
	PageNum     int64
	PageSize    int64
	Filter      string
	Fields      string
	Sort        string
	IsFound     bool
	IsTotalRows bool
	IsSum       bool
	Error       error
}

func NewFindPagingSumResult[T ddd.Entity](data []T, sumData []T, totalRows *int64, query FindPagingQuery, err error, sumErr error) *FindPagingResult[T] {
	res := NewFindPagingResult(data, totalRows, query, err)
	res.SumData = sumData
	if err == nil && sumErr != nil {
		res.Error = sumErr
	}
	return res
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
	}

	if query != nil {
		res.TotalPages = getTotalPage(totalRows, query.GetPageSize())
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

func NewFindPagingResultWithError[T ddd.Entity](err ...error) *FindPagingResult[T] {
	return &FindPagingResult[T]{
		Data:    []T{},
		IsFound: false,
		Error:   errors.News(err...),
	}
}

func (f *FindPagingResult[T]) GetDataLength() int64 {
	var data []T = f.Data
	v := len(data)
	return int64(v)
}

func (f *FindPagingResult[T]) GetSumDataLength() int64 {
	var data []T = f.SumData
	v := len(data)
	return int64(v)
}

func (f *FindPagingResult[T]) GetData() []T {
	return f.Data
}

func (f *FindPagingResult[T]) GetSumData() []T {
	return f.SumData
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

func (f *FindPagingResult[T]) SetSum(isSum bool, sumData []T, err error) *FindPagingResult[T] {
	f.IsSum = isSum
	f.SumData = sumData
	if f.Error == nil && err != nil {
		f.Error = err
	}
	return f
}

func (f *FindPagingResult[T]) Result() (*FindPagingResult[T], bool, error) {
	return f, f.IsFound, f.Error
}

func (f *FindPagingResult[T]) DataResult() ([]T, bool, error) {
	return f.Data, f.IsFound, f.Error
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
