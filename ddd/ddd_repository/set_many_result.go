package ddd_repository

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"go.mongodb.org/mongo-driver/mongo"
)

type SetManyResult[T interface{}] struct {
	err   error
	count int64
	data  []T
}

type SetManyCountResult struct {
	Error         error
	MatchedCount  int64       // The number of documents matched by the filter.
	ModifiedCount int64       // The number of documents modified by the operation.
	UpsertedCount int64       // The number of documents upserted by the operation.
	UpsertedID    interface{} // The _id field of the upserted document, or nil if no upsert was done.
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

func NewSetManyResult[T ddd.Entity](data []T, err error) *SetManyResult[T] {
	return &SetManyResult[T]{
		data: data,
		err:  err,
	}
}

func NewSetManyResultError[T ddd.Entity](err error) *SetManyResult[T] {
	return &SetManyResult[T]{
		err: err,
	}
}

func (s *SetManyResult[T]) GetError() error {
	return s.err
}

func (s *SetManyResult[T]) GetData() []T {
	return s.data
}

func (s *SetManyResult[T]) Result() ([]T, error) {
	return s.data, s.err
}

func (s *SetManyResult[T]) OnSuccess(success OnSuccessList[T]) *SetManyResult[T] {
	if s.err == nil && success != nil {
		s.err = success(s.data)
	}
	return s
}

func (s *SetManyResult[T]) OnError(err OnError) *SetManyResult[T] {
	if s.err != nil && err != nil {
		s.err = err(s.err)
	}
	return s
}

func (s *SetManyCountResult) GetError() error {
	return s.Error
}
