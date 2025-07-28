package restapi

import (
	"context"
	"encoding/base64"
	"github.com/goccy/go-json"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"time"
)

type CloudEvent struct {
	ID              string `json:"id"`
	SpecVersion     string `json:"specversion"`
	DataContentType string `json:"datacontenttype"`
	Source          string `json:"source"`
	Type            string `json:"type"`
	Topic           string `json:"topic"`
	PubsubName      string `json:"pubsubname"`
	DataBase64      string `json:"data_base64"`
	data            string `json:"-"`
}

type Event struct {
	PID         string          `json:"_id"`
	ID          string          `json:"id"`
	AppID       string          `json:"app_id"`
	CreatedTime time.Time       `json:"created_time"`
	Data        json.RawMessage `json:"data"` // 使用 json.RawMessage 来接收原始的 JSON 数据
	Meta        json.RawMessage `json:"meta"`
	TenantID    string          `json:"tenant_id"`
	Topic       string          `json:"topic"`
}

func (c *CloudEvent) GetData() (string, error) {
	if c.data == "" {
		dataBytes, err := base64.StdEncoding.DecodeString(c.DataBase64)
		if err != nil {
			return "", err
		}
		c.data = string(dataBytes)
	}
	return c.data, nil
}

func GetEventParams(ictx iris.Context, target interface{}, removeNames ...string) (res any, rctx context.Context, err error) {
	request := ictx.Request()
	if request.Method != iris.MethodPost {
		return nil, nil, errors.New("request body is null")
	}

	if request.ContentLength == 0 {
		return nil, nil, errors.New("request body is null")
	}
	var cloudEvent CloudEvent
	err = ictx.ReadJSON(&cloudEvent)
	if err != nil {
		return nil, nil, err
	}

	cloudData, err := cloudEvent.GetData()
	if err != nil {
		return nil, nil, err
	}

	var event Event
	err = json.Unmarshal([]byte(cloudData), &event)
	if err != nil {
		return nil, nil, err
	}
	rctx = context.Background()
	if event.Data == nil {
		return target, rctx, nil
	}
	err = json.Unmarshal(event.Data, target)
	if err != nil {
		return nil, nil, err
	}

	return target, rctx, err
}
