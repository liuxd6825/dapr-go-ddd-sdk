package events

import (
	"context"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/reflectutils"
)

// EventEntity
// @Description:  数据存储用的事件实体类
type EventEntity struct {
	Id          string         `json:"id" gorm:"id;primaryKey" bson:"id" `                                    // ID
	TenantId    string         `json:"tenantId" gorm:"tenant_id" bson:"tenant_id"`                            // 租户ID
	AppId       string         `json:"appId" gorm:"app_id" bson:"app_id"`                                     // 应用ID
	EventType   string         `json:"eventType" gorm:"event_type" bson:"event_type" `                        // 事件类型
	Data        map[string]any `json:"data" gorm:"data;json" bson:"data"`                                     // 事件数据
	Meta        map[string]any `json:"meta" gorm:"meta;json" bson:"meta"`                                     // 上下文数据
	CreatedTime *time.Time     `json:"createdTime" gorm:"created_time;index:,sort:desc," bson:"created_time"` // 发生事件
	Version     string         `json:"version" gorm:"version" bson:"version"  `                               // 事件数据版本
}

type EventPublisher interface {
	Publish(ctx context.Context, data *EventEntity) error
	PublishList(ctx context.Context, data []*EventEntity) error
}

func (e *EventEntity) Verify() error {
	err := errors.NewVerifyError()
	if e.Id == "" {
		err.AppendField("id", "id is required")
	}
	if e.TenantId == "" {
		err.AppendField("tenantId", "tenantId is required")
	}
	if e.AppId == "" {
		err.AppendField("appId", "appId is required")
	}
	if e.EventType == "" {
		err.AppendField("eventType", "eventType is required")
	}
	if e.Data == nil {
		err.AppendField("data", "data is required")
	}
	if e.Meta == nil {
		err.AppendField("meta", "meta is required")
	}
	if e.CreatedTime == nil {
		err.AppendField("createdTime", "createdTime is required")
	}
	if err.HasError() {
		return err
	}
	return nil
}

func GetEventType(eventData any) string {
	_, eventType, _ := reflectutils.GetTypeDetails(eventData)

	return eventType
}
