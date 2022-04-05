package daprsdk

import (
	"context"
	"encoding/json"
	sdk "github.com/dapr/go-sdk/client"
)

func InvokeMethod(ctx context.Context, client sdk.Client, appID, methodName, verb string, request interface{}, response interface{}) (interface{}, error) {
	var respBytes []byte
	var err error
	if request != nil {
		reqBytes, err := json.Marshal(request)
		if err != nil {
			return nil, err
		}
		content := &sdk.DataContent{
			ContentType: "application/json",
			Data:        reqBytes,
		}
		respBytes, err = client.InvokeMethodWithContent(ctx, appID, methodName, verb, content)
	} else {
		respBytes, err = client.InvokeMethod(ctx, appID, methodName, verb)
	}
	if err != nil {
		return nil, err
	}
	if len(respBytes) > 0 {
		err = json.Unmarshal(respBytes, response)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
	return nil, nil
}
