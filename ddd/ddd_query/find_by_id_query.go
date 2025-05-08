package ddd_query

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"

type FindByIdQuery = ddd_repository.FindByIdQuery

func NewFindByIdQuery() *FindByIdQuery {
	return &FindByIdQuery{}
}
