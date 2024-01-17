package ddd_service

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
)

type Finder interface {
	Init(queryAppId, resourceName, apiVersion string)
	FindById(ctx context.Context, tenantId, id string, resData interface{}) (isFound bool, err error)
	FindByIds(ctx context.Context, tenantId string, ids []string, resData interface{}) (isFound bool, err error)
}

type daprFinder struct {
	queryAppId   string
	resourceName string
	apiVersion   string
}

func NewFinder(queryAppId, queryResourceName, queryApiVersion string) Finder {
	finder := &daprFinder{}
	finder.Init(queryAppId, queryResourceName, queryApiVersion)
	return finder
}

func (s *daprFinder) Init(queryAppId, resourceName, apiVersion string) {
	s.queryAppId = queryAppId
	s.resourceName = resourceName
	s.apiVersion = apiVersion
}

func (s *daprFinder) FindById(ctx context.Context, tenantId, id string, resData interface{}) (isFound bool, err error) {
	return s.find(ctx, tenantId, "/"+id, nil, resData)
}

func (s *daprFinder) FindByIds(ctx context.Context, tenantId string, ids []string, resData interface{}) (isFound bool, err error) {
	idParams := ""
	count := len(ids)
	for i, id := range ids {
		idParams = idParams + fmt.Sprintf("id=%v", id)
		if i < count-2 {
			idParams += "&"
		}
	}
	methodName := fmt.Sprintf(":getById?%v", idParams)
	return s.find(ctx, tenantId, methodName, nil, resData)
}

func (s *daprFinder) find(ctx context.Context, tenantId, methodName string, req interface{}, resData interface{}) (isFound bool, err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()
	methodNameUrl := fmt.Sprintf("/api/%s/tenants/%s/%s%s", s.apiVersion, tenantId, s.resourceName, methodName)
	_, err = daprclient.GetDaprDDDClient().InvokeService(ctx, s.queryAppId, methodNameUrl, "get", req, resData)
	if err == nil {
		isFound = true
	}
	return isFound, err
}
