package ddd_query

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"

type FindPagingQuery = ddd_repository.FindPagingQueryRequest

type FindPagingResult[T any] = ddd_repository.FindPagingResult[any]

func NewFindPagingQuery() *FindPagingQuery {
	return &FindPagingQuery{}
}
