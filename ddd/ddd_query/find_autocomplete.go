package ddd_query

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"

type FindAutoCompleteQuery = ddd_repository.FindAutoCompleteQuery
type FindAutoCompleteQueryDTO = ddd_repository.FindAutoCompleteQueryDTO

func NewFindAutoCompleteQuery() FindAutoCompleteQuery {
	return ddd_repository.NewFindAutoCompleteQuery()
}

func NewFindAutoCompleteQueryDTO() *FindAutoCompleteQueryDTO {
	return ddd_repository.NewFindAutoCompleteQueryDTO()
}
