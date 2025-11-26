package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/events"
)

func GetEventParams(ictx iris.Context, event interface{}, removeNames ...string) (res any, rctx context.Context, err error) {
	request := ictx.Request()
	if request.Method != iris.MethodPost {
		return nil, nil, errors.New("request body is null")
	}

	if request.ContentLength == 0 {
		return nil, nil, errors.New("request body is null")
	}
	var cloudEvent events.CloudEvent
	err = ictx.ReadJSON(&cloudEvent)
	if err != nil {
		return nil, nil, err
	}
	event, rctx, err = events.LoadEvent(cloudEvent, event)
	return event, rctx, err
}
