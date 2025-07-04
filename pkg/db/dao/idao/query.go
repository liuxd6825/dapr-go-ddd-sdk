package idao

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"

type FindPagingQuery = store.FindPagingQuery
type FindPagingQueryBuilder = store.FindPagingQueryBuilder

func NewFindPagingQueryBuilder() FindPagingQueryBuilder {
	return store.NewFindPagingQueryBuilder()
}

func NewFindPagingQuery() FindPagingQuery {
	return store.NewFindPagingQuery()
}
