package ddd_query

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

type FindAutoCompleteQuery = store.FindAutoCompleteQuery
type FindAutoCompleteQueryDTO = store.FindAutoCompleteQueryDTO

func NewFindAutoCompleteQuery() FindAutoCompleteQuery {
	return store.NewFindAutoCompleteQuery()
}

func NewFindAutoCompleteQueryDTO() *FindAutoCompleteQueryDTO {
	return store.NewFindAutoCompleteQueryDTO()
}
