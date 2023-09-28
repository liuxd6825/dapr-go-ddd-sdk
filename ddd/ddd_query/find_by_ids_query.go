package ddd_query

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"

type FindByIdsQuery = ddd_repository.FindByIdsQuery

func NewFindByIdsQuery() *FindByIdsQuery {
	return &FindByIdsQuery{}
}
