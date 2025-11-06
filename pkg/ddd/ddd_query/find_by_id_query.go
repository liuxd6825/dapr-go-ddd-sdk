package ddd_query

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

type FindByIdQuery = store.FindByIdQueryRequest

func NewFindByIdQuery() *FindByIdQuery {
	return &FindByIdQuery{}
}
