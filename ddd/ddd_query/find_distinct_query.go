package ddd_query

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"

type FindDistinctQuery = ddd_repository.FindDistinctQuery
type FindDistinctQueryDTO = ddd_repository.FindDistinctQueryDTO

func NewFindDistinctQuery() FindDistinctQuery {
	return ddd_repository.NewFindDistinctQuery()
}

func NewFindDistinctQueryDTO() *FindDistinctQueryDTO {
	return ddd_repository.NewFindDistinctQueryDTO()
}
