package ddd_query

import (
	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

type FindPagingQuery = store2.FindPagingQueryRequest
type FindPagingResult = store2.FindPagingResult[any]
type FindPagingByCaseIdQuery = store2.FindPagingByCaseIdQueryRequest

func NewFindPagingQuery() store2.FindPagingQuery {
	return store2.NewFindPagingQuery()
}

func NewFindPagingByCaseIdQuery(paging *FindPagingQuery, caseId string) *FindPagingByCaseIdQuery {
	return store2.NewFindPagingByCaseIdQuery(paging, caseId)
}

func NewFindPagingQueryDTO() *store2.FindPagingQueryDTO {
	return store2.NewFindPagingQueryDTO()
}
