package ddd_repository

type SetManyResult[T interface{}] struct {
	Error error `json:"error"`
	Count int64 `json:"count"`
	Data  []T   `json:"data"`
}

type SetManyCountResult struct {
	Error         error       `json:"error"`
	MatchedCount  int64       `json:"matchedCount"`  // The number of documents matched by the filter.
	ModifiedCount int64       `json:"modifiedCount"` // The number of documents modified by the operation.
	UpsertedCount int64       `json:"upsertedCount"` // The number of documents upserted by the operation.
	UpsertedID    interface{} `json:"upsertedId"`    // The _id field of the upserted document, or nil if no upsert was done.
}

type UpdateResult interface {
	GetMatchedCount() int64     // The number of documents matched by the filter.
	GetModifiedCount() int64    // The number of documents modified by the operation.
	GetUpsertedCount() int64    // The number of documents upserted by the operation.
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
		res.ModifiedCount = updateRes.GetModifiedCount()
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
func (r *SetManyCountResult) SetModifiedCount(v int64) *SetManyCountResult {
	r.ModifiedCount = v
	return r
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

func (s *SetManyResult[T]) GetData() []T {
	return s.Data
}

func (s *SetManyResult[T]) Result() ([]T, error) {
	return s.Data, s.Error
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
