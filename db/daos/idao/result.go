package idao

type Result struct {
	RowsAffected int64 `json:"rowsAffected"` // 影响行数
}

type RowsAffected interface {
	GetRowsAffected() int64
}

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
