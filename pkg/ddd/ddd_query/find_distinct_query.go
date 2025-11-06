package ddd_query

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

type FindDistinctQuery = store.FindDistinctQuery
type FindDistinctQueryDTO = store.FindDistinctQueryDTO

func NewFindDistinctQuery() FindDistinctQuery {
	return store.NewFindDistinctQuery()
}

func NewFindDistinctQueryDTO() *FindDistinctQueryDTO {
	return store.NewFindDistinctQueryDTO()
}
