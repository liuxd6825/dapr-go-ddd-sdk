package restapi

import (
	"github.com/kataras/iris/v12"
	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

type FindByIdRequest struct {
	Id string `json:"id" path:"id" required:"true"`
}

type FindByIdsRequest struct {
	Ids []string `json:"ids" path:"ids" required:"true"` // 聚合根Id列表
}

type FindPagingRequest = store2.FindPagingQueryRequest
type FindDistinctRequest = store2.FindDistinctQueryRequest
type FindAutoCompleteRequest = store2.FindAutoCompleteQueryRequest
type FindPagingByCaseIdRequest = store2.FindPagingByCaseIdQueryRequest

var restAssembler = RestAssembler{}

func GetFindAutoCompleteRequest(ictx iris.Context) (*FindAutoCompleteRequest, error) {
	return restAssembler.AsFindAutoCompleteRequest(ictx)
}

func GetDistinctRequest(ictx iris.Context) (*FindDistinctRequest, error) {
	return restAssembler.AsDistinctRequest(ictx)
}

func GetFindPagingRequest(ictx iris.Context) (*FindPagingRequest, error) {
	return restAssembler.AsFindPagingRequest(ictx)
}
