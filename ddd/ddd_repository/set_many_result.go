package ddd_repository

type SetManyResult[T interface{}] struct {
	Error        error `json:"error"`
	RowsAffected int64 `json:"rowsAffected"`
	Data         []T   `json:"data"`
}

type SetManyCountResult struct {
	Error         error       `json:"error"`
	MatchedCount  int64       `json:"matchedCount"`  // 查询条件匹配到的文档数量。
	RowsAffected  int64       `json:"rowsAffected"`  // 表示实际被修改的文档数量。
	UpsertedCount int64       `json:"upsertedCount"` // 表示执行 upsert 操作时插入的新文档数量
	UpsertedID    interface{} `json:"upsertedId"`    // The _id field of the upserted document, or nil if no upsert was done.
}

type UpdateResult interface {
	GetMatchedCount() int64     // 查询条件匹配到的文档数量。
	GetRowsAffected() int64     // 表示实际被修改的文档数量。
	GetUpsertedCount() int64    // 表示执行 upsert 操作时插入的新文档数量
	GetUpsertedID() interface{} // The _id field of the upserted document, or nil if no upsert was done.
}

func NewSetManyCountResultError(err error) *SetManyCountResult {
	return &SetManyCountResult{
		Error: err,
	}
}
func NewSetManyCountResult() *SetManyCountResult {
	return &SetManyCountResult{}
}
func NewSetManyCountResultMongo(updateRes UpdateResult, err error) *SetManyCountResult {
	res := &SetManyCountResult{
		Error: err,
	}
	if updateRes != nil {
		res.UpsertedCount = updateRes.GetUpsertedCount()
		res.RowsAffected = updateRes.GetRowsAffected()
		res.MatchedCount = updateRes.GetMatchedCount()
		res.UpsertedID = updateRes.GetUpsertedID()
	}
	return res
}

func (r *SetManyCountResult) SetError(err error) *SetManyCountResult {
	r.Error = err
	return r
}
func (r *SetManyCountResult) SetMatchedCount(v int64) *SetManyCountResult {
	r.MatchedCount = v
	return r
}
func (r *SetManyCountResult) SetRowsAffected(v int64) *SetManyCountResult {
	r.RowsAffected = v
	return r
}
func (r *SetManyCountResult) GetRowsAffected() int64 {
	return r.RowsAffected
}

func (r *SetManyCountResult) SetUpsertedCount(v int64) *SetManyCountResult {
	r.UpsertedCount = v
	return r
}
func (r *SetManyCountResult) SetUpsertedID(v any) *SetManyCountResult {
	r.UpsertedID = v
	return r
}

func NewSetManyResult[T any](data []T, err error) *SetManyResult[T] {
	return &SetManyResult[T]{
		Data:  data,
		Error: err,
	}
}

func NewSetManyResultError[T any](err error) *SetManyResult[T] {
	return &SetManyResult[T]{
		Error: err,
	}
}

func (s *SetManyResult[T]) GetError() error {
	return s.Error
}

func (s *SetManyResult[T]) SetError(err error) *SetManyResult[T] {
	s.Error = err
	return s
}

func (s *SetManyResult[T]) GetData() []T {
	return s.Data
}

func (s *SetManyResult[T]) Result() ([]T, error) {
	return s.Data, s.Error
}

func (r *SetManyResult[T]) GetRowsAffected() int64 {
	return r.RowsAffected
}

func (r *SetManyResult[T]) SetRowsAffected(val int64) *SetManyResult[T] {
	r.RowsAffected = val
	return r
}

func (s *SetManyResult[T]) OnSuccess(success OnSuccessList[T]) *SetManyResult[T] {
	if s.Error == nil && success != nil {
		s.Error = success(s.Data)
	}
	return s
}

func (s *SetManyResult[T]) OnError(err OnError) *SetManyResult[T] {
	if s.Error != nil && err != nil {
		s.Error = err(s.Error)
	}
	return s
}

func (r *SetManyCountResult) GetError() error {
	return r.Error
}
