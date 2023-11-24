package daprclient

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
)

type MethodType int

type ApplicationClient struct {
	query *AppConfig
	cmd   *AppConfig
}

type AppConfig struct {
	AppId        string
	ApiVersion   string
	ResourceName string
	AppPort      int
	DaprHttpPort int
	DaprGrpcPort int
}

const (
	MethodTypeGet MethodType = iota
	MethodTypePost
	MethodTypePut
	MethodTypeDelete
	MethodTypePatch
)

func NewApplicationClient(cmd *AppConfig, query *AppConfig) *ApplicationClient {
	service := &ApplicationClient{
		cmd:   cmd,
		query: query,
	}
	return service
}

/*
func (s *ApplicationClient) HttpGet(ctx context.Context, tenantId, methodName string, paras ...string) (res *Response, err error) {
	defer func() {
		if e := errors.GetRecoverError(recover()); e != nil {
			err = e
		}
	}()

	url, err := s.getHttpUrl(s.query.DaprHttpPort, tenantId, methodName, paras...)
	if err != nil {
		return nil, err
	}
	return GetDaprDDDClient().HttpGet(ctx, url)
}

func (s *ApplicationClient) HttpPost(ctx context.Context, tenantId, methodName string, requestData interface{}, paras ...string) (res *Response, err error) {
	defer func() {
		if e := errors.GetRecoverError(recover()); e != nil {
			err = e
		}
	}()
	url, err := s.getHttpUrl(s.HttpCmdPort, tenantId, methodName, paras...)
	if err != nil {
		return nil, err
	}
	resp, err := GetDaprDDDClient().HttpPost(ctx, url, requestData)
	return resp, err
}

func (s *ApplicationClient) HttpPut(ctx context.Context, tenantId, methodName string, requestData interface{}, paras ...string) (res *Response, err error) {
	defer func() {
		if e := errors.GetRecoverError(recover()); e != nil {
			err = e
		}
	}()
	url, err := s.getHttpUrl(s.HttpCmdPort, tenantId, methodName, paras...)
	if err != nil {
		return nil, err
	}
	return GetDaprDDDClient().HttpPut(ctx, url, requestData)
}

func (s *ApplicationClient) HttpDelete(ctx context.Context, tenantId, methodName string, requestData interface{}, paras ...string) (res *Response, err error) {
	defer func() {
		if e := errors.GetRecoverError(recover()); e != nil {
			err = e
		}
	}()
	url, err := s.getHttpUrl(s.HttpCmdPort, tenantId, methodName, paras...)
	if err != nil {
		return nil, err
	}
	return GetDaprDDDClient().HttpDelete(ctx, url, requestData)
}
*/

func (s *ApplicationClient) DoCommand(ctx context.Context, tenantId string, methodName string, methodType MethodType, command any, response any) (any, error) {
	return s.InvokeService(ctx, s.cmd, tenantId, methodName, methodType, command, response)
}

func (s *ApplicationClient) QueryById(ctx context.Context, tenantId, id string, resData interface{}) (isFound bool, err error) {
	return s.QueryData(ctx, tenantId, "/"+id, nil, resData)
}

func (s *ApplicationClient) QueryByIds(ctx context.Context, tenantId string, ids []string, resData interface{}) (isFound bool, err error) {
	idParams := ""
	count := len(ids)
	for i, id := range ids {
		idParams = idParams + fmt.Sprintf("id=%v", id)
		if i < count-2 {
			idParams += "&"
		}
	}
	methodName := fmt.Sprintf(":getByIds?%v", idParams)
	return s.QueryData(ctx, tenantId, methodName, nil, resData)
}

func (s *ApplicationClient) QueryData(ctx context.Context, tenantId, methodName string, request interface{}, response interface{}) (isFound bool, err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()
	_, err = s.InvokeService(ctx, s.query, tenantId, methodName, MethodTypeGet, request, response)
	if err == nil {
		isFound = true
	}
	return isFound, err
}

func (s *ApplicationClient) getMethodName(methodName string, params ...string) (string, error) {
	var err error
	res := methodName
	count := len(params)
	for i := 0; i < count; i++ {
		res = fmt.Sprintf("%v=%v", params[i], params[i+1])
	}
	return res, err
}

func (s *ApplicationClient) InvokeService(ctx context.Context, config *AppConfig, tenantId, methodName string, methodType MethodType, request interface{}, response interface{}) (res any, err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()
	methodNameUrl := fmt.Sprintf("/api/%s/tenants/%s/%s%s", config.ApiVersion, tenantId, config.ResourceName, methodName)
	return GetDaprDDDClient().InvokeService(ctx, config.AppId, methodNameUrl, methodType.ToString(), request, response)
}

func (m MethodType) ToString() string {
	switch m {
	case MethodTypeGet:
		return "get"
	case MethodTypePost:
		return "post"
	case MethodTypePut:
		return "put"
	case MethodTypeDelete:
		return "delete"
	case MethodTypePatch:
		return "patch"
	}
	return ""
}
