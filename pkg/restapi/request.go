package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_query"
)

var restAssembler = restapp.RestAssembler{}

func GetFindByIdRequest(ictx iris.Context) (*ddd_query.FindByIdQuery, error) {
	return restAssembler.AsFindByIdRequest(ictx)
}

func GetFindByIdsRequest(ictx iris.Context) (*ddd_query.FindByIdsQuery, error) {
	return restAssembler.AsFindByIdsRequest(ictx)
}

func GetFindAllRequest(ictx iris.Context) (*ddd_query.FindAllQuery, error) {
	return restAssembler.AsFindAllRequest(ictx)
}

func GetFindAutoCompleteRequest(ictx iris.Context) (ddd_query.FindAutoCompleteQuery, error) {
	return restAssembler.AsFindAutoCompleteRequest(ictx)
}

func GetDistinctRequest(ictx iris.Context) (ddd_query.FindDistinctQuery, error) {
	return restAssembler.AsDistinctRequest(ictx)
}

func GetFindPagingRequest(ictx iris.Context) (*ddd_query.FindPagingQuery, error) {
	return restAssembler.AsFindPagingRequest(ictx)
}
