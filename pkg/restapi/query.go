package restapi

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

type FindPagingQueryRequest = store.FindPagingQueryRequest
type FindByIdQueryRequest struct {
	Id string `json:"id" param:"id" validate:"required"`
}
