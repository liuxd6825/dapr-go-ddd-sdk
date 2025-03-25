package ddd_query

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"

type FindPagingQuery = store.FindPagingQueryRequest
type FindPagingResult = store.FindPagingResult[any]
type FindPagingByCaseIdQuery = store.FindPagingByCaseIdQueryRequest

func NewFindPagingQuery() store.FindPagingQuery {
	return store.NewFindPagingQuery()
}

func NewFindPagingByCaseIdQuery(paging *FindPagingQuery, caseId string) *FindPagingByCaseIdQuery {
	return store.NewFindPagingByCaseIdQuery(paging, caseId)
}

func NewFindPagingQueryDTO() *store.FindPagingQueryDTO {
	return store.NewFindPagingQueryDTO()
}
