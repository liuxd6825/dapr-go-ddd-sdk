package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
)

type FindByIdRequest = store.FindByIdQueryRequest
type FindPagingRequest = store.FindPagingQueryRequest
type FindDistinctRequest = store.FindDistinctQueryRequest
type FindByIdsRequest = store.FindByIdsQueryRequest
type FindAllRequest = store.FindAllQueryRequest
type FindAutoCompleteRequest = store.FindAutoCompleteQueryRequest
type FindPagingByCaseIdRequest = store.FindPagingByCaseIdQueryRequest

var restAssembler = RestAssembler{}

func GetFindByIdRequest(ictx iris.Context) (*FindByIdRequest, error) {
	return restAssembler.AsFindByIdRequest(ictx)
}

func GetFindByIdsRequest(ictx iris.Context) (*FindByIdsRequest, error) {
	return restAssembler.AsFindByIdsRequest(ictx)
}

func GetFindAllRequest(ictx iris.Context) (*FindAllRequest, error) {
	return restAssembler.AsFindAllRequest(ictx)
}

func GetFindAutoCompleteRequest(ictx iris.Context) (*FindAutoCompleteRequest, error) {
	return restAssembler.AsFindAutoCompleteRequest(ictx)
}

func GetDistinctRequest(ictx iris.Context) (*FindDistinctRequest, error) {
	return restAssembler.AsDistinctRequest(ictx)
}

func GetFindPagingRequest(ictx iris.Context) (*FindPagingRequest, error) {
	return restAssembler.AsFindPagingRequest(ictx)
}
