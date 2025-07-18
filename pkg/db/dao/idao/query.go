package idao

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"

type FindPagingQuery = store.FindPagingQuery
type FindPagingQueryBuilder = store.FindPagingQueryBuilder
type FindPagingQueryRequest = store.FindPagingQueryRequest
type FindPagingByCaseIdQuery = store.FindPagingByCaseIdQuery
type FindPagingByCaseIdQueryRequest = store.FindPagingByCaseIdQueryRequest
type FindByIdQueryRequest = store.FindByIdQueryRequest

func NewFindPagingQueryBuilder() FindPagingQueryBuilder {
	return store.NewFindPagingQueryBuilder()
}

func NewFindPagingQuery() FindPagingQuery {
	return store.NewFindPagingQuery()
}
