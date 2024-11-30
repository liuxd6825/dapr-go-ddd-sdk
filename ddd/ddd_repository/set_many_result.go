package ddd_repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

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

func NewSetManyCountResultError(err error) *SetManyCountResult {
	return &SetManyCountResult{
		Error: err,
	}
}
func NewSetManyCountResult(updateRes *mongo.UpdateResult, err error) *SetManyCountResult {
	res := &SetManyCountResult{
		Error: err,
	}
	if updateRes != nil {
		res.UpsertedCount = updateRes.UpsertedCount
		res.ModifiedCount = updateRes.ModifiedCount
		res.MatchedCount = updateRes.MatchedCount
		res.UpsertedID = updateRes.UpsertedID
	}
	return res
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

func (s *SetManyCountResult) GetError() error {
	return s.Error
}
