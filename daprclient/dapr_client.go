package daprclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/dapr/go-sdk/client"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_errors"
)

func (c *DaprClient) InvokeMethod(ctx context.Context, appID, methodName, verb string, request interface{}, response interface{}) (interface{}, error) {
	var err error
	defer func() {
		if e := ddd_errors.GetRecoverError(recover()); e != nil {
			err = e
		}
	}()
	var respBytes []byte

	if request != nil {
		reqBytes, err := json.Marshal(request)
		if err != nil {
			return nil, newAppError(appID, err)
		}
		content := &sdk.DataContent{
			ContentType: "application/json",
			Data:        reqBytes,
		}
		respBytes, err = c.grpcClient.InvokeMethodWithContent(ctx, appID, methodName, verb, content)
	} else {
		respBytes, err = c.grpcClient.InvokeMethod(ctx, appID, methodName, verb)
	}
	if err != nil {
		return nil, newAppError(appID, err)
	}
	if len(respBytes) > 0 {
		err = json.Unmarshal(respBytes, response)
		if err != nil {
			return nil, newAppError(appID, err)
		}
		return response, nil
	}
	return nil, nil
}

func newAppError(appID string, err error) error {
	msg := fmt.Sprintf("AppId is %s , %s", appID, err.Error())
	return errors.New(msg)
}
