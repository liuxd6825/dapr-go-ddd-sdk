package idao

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"

type Result struct {
	RowsAffected int64 `json:"rowsAffected"`    // 影响行数
	Error        error `json:"error,omitempty"` // 错误信息
}

type FindPagingResult[T any] interface {
	store.FindPagingResult[T]
}

type RowsAffected interface {
	GetRowsAffected() int64
}

type RowsError interface {
	GetError() error
}

/*func NewFindPagingResult[T any]() *FindPagingResult[T] {
	return &FindPagingResult[T]{}
}*/

func NewResult(rows RowsAffected) *Result {
	if rows != nil {
		return &Result{RowsAffected: rows.GetRowsAffected()}
	}
	return &Result{}
}

func (r *Result) GetRowsAffected() int64 {
	return r.RowsAffected
}

func (r *Result) SetRowsAffected(val int64) *Result {
	r.RowsAffected = val
	return r
}

func (r *Result) GetError() error {
	return r.Error
}

func (r *Result) SetError(val error) *Result {
	r.Error = val
	return r
}
