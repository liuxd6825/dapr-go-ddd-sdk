package store

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
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
	IsFound     bool   `json:"IsFound"`              // 是否找到数据
	IsTotalRows bool   `json:"isTotalRows"`          // 是否统计总记录数
}

type FindPagingResult[T any] interface {
	GetData() []T
	SetData(data []T)

	GetSumData() []T
	SetSumData(data []T)

	GetTotalRows() int64
	SetTotalRows(totalRows int64)

	GetTotalPages() int64
	SetTotalPages(totalPages int64)

	GetPageNum() int64
	SetPageNum(pageNum int64)

	GetPageSize() int64
	SetPageSize(pageSize int64)

	GetFilter() string
	SetFilter(filter string)

	GetFields() string
	SetFields(val string)

	GetSort() string
	SetSort(sort string)

	GetIsFound() bool
	SetIsFound(isFound bool)

	GetIsTotalRows() bool
	SetIsTotalRows(v bool)

	GetIsSum() bool
	SetIsSum(val bool)

	GetError() error
	SetError(err error)

	GetDataLength() int64
}

type FindPagingResultStruct[T any] struct {
	Data        []T    `json:"data"`
	SumData     []T    `json:"sumData,omitempty"`
	TotalRows   int64  `json:"totalRows,omitempty"`
	TotalPages  int64  `json:"totalPages,omitempty"`
	PageNum     int64  `json:"pageNum,omitempty"`
	PageSize    int64  `json:"pageSize,omitempty"`
	Filter      string `json:"filter"`
	Fields      string `json:"fields"`
	Sort        string `json:"sort"`
	IsFound     bool   `json:"isFound"`
	IsTotalRows bool   `json:"isTotalRows"`
	IsSum       bool   `json:"isSum"`
	Error       error  `json:"error,omitempty"`
}

func (f *FindPagingResultStruct[T]) SetIsSum(val bool) {
	f.IsSum = val
}

func (f *FindPagingResultStruct[T]) SetSumData(data []T) {
	f.SumData = data
}

func (f *FindPagingResultStruct[T]) SetTotalRows(totalRows int64) {
	f.TotalRows = totalRows
}

func (f *FindPagingResultStruct[T]) SetPageNum(pageNum int64) {
	f.PageNum = pageNum
}

func (f *FindPagingResultStruct[T]) SetPageSize(pageSize int64) {
	f.PageSize = pageSize
}

func (f *FindPagingResultStruct[T]) SetFilter(filter string) {
	f.Filter = filter
}

func (f *FindPagingResultStruct[T]) SetFields(val string) {
	f.Fields = val
}

func (f *FindPagingResultStruct[T]) SetSort(sort string) {
	f.Sort = sort
}

func (f *FindPagingResultStruct[T]) SetIsFound(isFound bool) {
	f.IsFound = isFound
}

func (f *FindPagingResultStruct[T]) SetIsTotalRows(v bool) {
	f.IsTotalRows = v
}

type FindPagingResultOptions[T interface{}] struct {
	Data        *[]T   `json:"data"`
	SumData     *[]T   `json:"sumData"`
	TotalRows   int64  `json:"totalRows"`
	TotalPages  int64  `json:"totalPages"`
	PageNum     int64  `json:"pageNum"`
	PageSize    int64  `json:"pageSize"`
	Filter      string `json:"filter"`
	Fields      string `json:"fields"`
	Sort        string `json:"sort"`
	IsFound     bool   `json:"isFound"`
	IsTotalRows bool   `json:"isTotalRows"`
	IsSum       bool   `json:"isSum"`
	Error       error  `json:"error,omitempty"`
}

func NewFindPagingSumResult[T any](data []T, sumData []T, totalRows *int64, query FindPagingQuery, err error, sumErr error) FindPagingResult[T] {
	var total int64
	if totalRows != nil {
		total = *totalRows
	}
	res := NewFindPagingResult[T](data, total, query, err)
	res.SetSumData(sumData)
	if err == nil && sumErr != nil {
		res.SetError(sumErr)
	}
	return res
}

func NewFindPagingResultEmpty[T any]() FindPagingResult[T] {
	res := &FindPagingResultStruct[T]{
		Data:        nil,
		TotalRows:   0,
		TotalPages:  0,
		PageNum:     0,
		PageSize:    0,
		Sort:        "",
		Filter:      "",
		IsFound:     false,
		IsTotalRows: false,
		Error:       nil,
	}
	return res
}

func NewFindPagingResult[T any](data []T, totalRows int64, query FindPagingQuery, err error) FindPagingResult[T] {
	return NewFindPagingResultStruct[T](data, totalRows, query, err)
}

func NewFindPagingResultStruct[T any](data []T, totalRows int64, query FindPagingQuery, err error) *FindPagingResultStruct[T] {
	res := &FindPagingResultStruct[T]{
		Data:        data,
		TotalRows:   0,
		TotalPages:  0,
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

	res.TotalRows = totalRows

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

func NewFindPagingResultOptions[T any]() *FindPagingResultOptions[T] {
	return &FindPagingResultOptions[T]{}
}

func NewFindPagingResultWithError[T any](err ...error) FindPagingResult[T] {
	return &FindPagingResultStruct[T]{
		Data:    []T{},
		IsFound: false,
		Error:   errors.News(err...),
	}
}

func (f *FindPagingResultStruct[T]) GetIsSum() bool {
	return f.IsSum
}

func (f *FindPagingResultStruct[T]) GetDataLength() int64 {
	var data []T = f.Data
	v := len(data)
	return int64(v)
}

func (f *FindPagingResultStruct[T]) GetSumDataLength() int64 {
	var data []T = f.SumData
	v := len(data)
	return int64(v)
}

func (f *FindPagingResultStruct[T]) GetData() []T {
	return f.Data
}

func (f *FindPagingResultStruct[T]) GetSumData() []T {
	return f.SumData
}

func (f *FindPagingResultStruct[T]) GetAnyData() any {
	return f.Data
}

func (f *FindPagingResultStruct[T]) GetTotalRows() int64 {
	return f.TotalRows
}

func (f *FindPagingResultStruct[T]) GetTotalPages() int64 {
	return f.TotalPages
}

func (f *FindPagingResultStruct[T]) GetPageNum() int64 {
	return f.PageNum
}

func (f *FindPagingResultStruct[T]) GetPageSize() int64 {
	return f.PageSize
}

func (f *FindPagingResultStruct[T]) GetFilter() string {
	return f.Filter
}

func (f *FindPagingResultStruct[T]) GetFields() string {
	return f.Fields
}

func (f *FindPagingResultStruct[T]) GetSort() string {
	return f.Sort
}

func (f *FindPagingResultStruct[T]) GetIsFound() bool {
	return f.IsFound
}

func (f *FindPagingResultStruct[T]) GetIsTotalRows() bool {
	return f.IsTotalRows
}

func (f *FindPagingResultStruct[T]) GetError() error {
	return f.Error
}

func (f *FindPagingResultStruct[T]) SetError(err error) {
	f.Error = err
}

func (f *FindPagingResultStruct[T]) SetData(data []T) {
	f.Data = data
}

func (f *FindPagingResultStruct[T]) SetTotalPages(val int64) {
	f.TotalPages = val
}

func (f *FindPagingResultStruct[T]) SetTotalRow(val int64) {
	f.TotalRows = val
}

func (f *FindPagingResultStruct[T]) SetSum(isSum bool, sumData []T, err error) {
	f.IsSum = isSum
	f.SumData = sumData
	if f.Error == nil && err != nil {
		f.Error = err
	}
}

func (f *FindPagingResultStruct[T]) Result() (*FindPagingResultStruct[T], bool, error) {
	return f, f.IsFound, f.Error
}

func (f *FindPagingResultStruct[T]) DataResult() ([]T, bool, error) {
	return f.Data, f.IsFound, f.Error
}

func (f *FindPagingResultStruct[T]) OnError(onErr OnError) *FindPagingResultStruct[T] {
	if f.Error != nil && onErr != nil {
		f.Error = onErr(f.Error)
	}
	return f
}

func (f *FindPagingResultStruct[T]) OnNotFond(fond OnIsFond) *FindPagingResultStruct[T] {
	if f.Error == nil && !f.IsFound && fond != nil {
		f.Error = fond()
	}
	return f
}

func (f *FindPagingResultStruct[T]) OnSuccess(success OnSuccessList[T]) *FindPagingResultStruct[T] {
	if f.Error == nil && success != nil && f.IsFound {
		f.Error = success(f.Data)
	}
	return f
}

func (f *FindPagingResultOptions[T]) SetData(data *[]T) {
	f.Data = data
}

func (f *FindPagingResultOptions[T]) SetTotalRows(totalRows int64) {
	f.TotalRows = totalRows
}

func (f *FindPagingResultOptions[T]) SetTotalPages(totalPages int64) {
	f.TotalPages = totalPages
}

func (f *FindPagingResultOptions[T]) SetPageNum(pageNum int64) {
	f.PageNum = pageNum
}

func (f *FindPagingResultOptions[T]) SetPageSize(pageSize int64) {
	f.PageSize = pageSize
}

func (f *FindPagingResultOptions[T]) SetFilter(filter string) {
	f.Filter = filter
}

func (f *FindPagingResultOptions[T]) SetSort(sort string) {
	f.Sort = sort
}

func (f *FindPagingResultOptions[T]) SetError(err error) {
	f.Error = err
}

func (f *FindPagingResultOptions[T]) SetIsFound(isFound bool) {
	f.IsFound = isFound
}

func (f *FindPagingResultOptions[T]) SetIsTotalRows(v bool) {
	f.IsTotalRows = v
}

func getTotalPage(totalRows int64, pageSize int64) int64 {
	if totalRows == 0 {
		return 0
	}
	if pageSize == 0 {
		return 0
	}
	rows := totalRows
	totalPage := rows / pageSize
	if rows%pageSize >= 1 {
		totalPage++
	}
	return totalPage
}
