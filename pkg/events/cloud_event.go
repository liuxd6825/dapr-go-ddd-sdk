package events

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

// CloudEvent 消息事件的标准结构
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

// EventPayload 事件的载荷对象 对应发送端的OutboxEvent
type EventPayload struct {
	PID         string          `json:"_id"`
	ID          string          `json:"id"`
	AppID       string          `json:"app_id"`
	CreatedTime time.Time       `json:"created_time"`
	Data        json.RawMessage `json:"data"` // 使用 json.RawMessage 来接收原始的 JSON 数据
	Meta        json.RawMessage `json:"meta"`
	EventType   string          `json:"event_type"`
	TenantId    string          `json:"tenant_id"`
}

// Event 事件本体
type Event[T any] struct {
	EventId    string    `json:"eventId" bson:"eventId" validate:"required" `
	OccurredOn time.Time `json:"occurredOn" bson:"occurredOn" validate:"required" `
	Data       T         `json:"data"  bson:"data" validate:"required"`
}

const (
	TENANT_ID = "tenantId"
	AUTH_USER = "authUser"
	HEADER    = "header"
)

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
