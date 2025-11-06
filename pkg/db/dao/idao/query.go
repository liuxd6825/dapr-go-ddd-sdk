package idao

import (
	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

type FindPagingQuery = store2.FindPagingQuery
type FindPagingQueryBuilder = store2.FindPagingQueryBuilder
type FindPagingQueryRequest = store2.FindPagingQueryRequest
type FindPagingByCaseIdQuery = store2.FindPagingByCaseIdQuery
type FindPagingByCaseIdQueryRequest = store2.FindPagingByCaseIdQueryRequest
type FindByIdQueryRequest = store2.FindByIdQueryRequest

func NewFindPagingQueryBuilder() FindPagingQueryBuilder {
	return store2.NewFindPagingQueryBuilder()
}

func NewFindPagingQuery() FindPagingQuery {
	return store2.NewFindPagingQuery()
}
