package ddd_query

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"

type FindPagingQuery = ddd_repository.FindPagingQueryRequest
type FindPagingResult = ddd_repository.FindPagingResult[any]

func NewFindPagingQuery() ddd_repository.FindPagingQuery {
	return ddd_repository.NewFindPagingQuery()
}

func NewFindPagingQueryDTO() *ddd_repository.FindPagingQueryDTO {
	return ddd_repository.NewFindPagingQueryDTO()
}
