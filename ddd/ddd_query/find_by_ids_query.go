package ddd_query

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
)

type FindByIdsQuery = store.FindByIdsQueryRequest

func NewFindByIdsQuery() *FindByIdsQuery {
	return &FindByIdsQuery{}
}
