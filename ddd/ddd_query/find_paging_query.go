package ddd_query

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"

type FindPagingQuery = ddd_repository.FindPagingQueryRequest
type FindPagingResult = ddd_repository.FindPagingResult[any]
type FindPagingByCaseIdQuery = ddd_repository.FindPagingByCaseIdQueryRequest

func NewFindPagingQuery() ddd_repository.FindPagingQuery {
	return ddd_repository.NewFindPagingQuery()
}

func NewFindPagingByCaseIdQuery(paging *FindPagingQuery, caseId string) *FindPagingByCaseIdQuery {
	return ddd_repository.NewFindPagingByCaseIdQuery(paging, caseId)
}

func NewFindPagingQueryDTO() *ddd_repository.FindPagingQueryDTO {
	return ddd_repository.NewFindPagingQueryDTO()
}
