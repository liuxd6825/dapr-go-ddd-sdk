package dbevent

import "github.com/liuxd6825/dapr-go-ddd-sdk/types/times"

/*
type Outbox struct {
	Id          string              `gorm:"id;primaryKey" bson:"id" json:"id"`
	TopicName   string              `gorm:"topic_name" bson:"topic_name" json:"topicName,omitempty"`
	TenantId    string              `gorm:"tenant_id" bson:"tenant_id" json:"tenantId,omitempty"` // 租户ID
	AppId       string              `gorm:"app_name" bson:"app_id" json:"appName,omitempty"`      // 应用ID
	EventId     string              `gorm:"event_id" bson:"event_id" json:"eventId,omitempty"`    // 事件ID
	EventType   string              `gorm:"event_type" json:"eventType,omitempty"`                // 事件类型
	EventVer    string              `gorm:"event_ver" json:"eventVer,omitempty"`                  // 事件版本
	CommandId   string              `gorm:"command_id" json:"commandId,omitempty"`                // 命令ID
	AggId       string              `gorm:"agg_id" json:"aggId,omitempty"`                        // 聚合ID
	AggType     string              `gorm:"agg_type" json:"aggType,omitempty"`                    // 聚合类型
	CreatedTime *time.Time          `gorm:"created_time;index:,sort:desc," json:"createdTime"`    // 发生事件
	Data        any                 `gorm:"type:json;serializer:json" json:"data,omitempty"`      // 事件数据
	Metadata    map[string][]string `gorm:"type:json;serializer:json" json:"context,omitempty"`   // 上下文数据
}*/

type OutboxEvent struct {
	Id          string      `gorm:"id;primaryKey" bson:"id" json:"id"`
	TenantId    string      `gorm:"tenant_id" bson:"tenant_id" json:"tenant_id"` // 租户ID
	AppId       string      `gorm:"app_id" bson:"app_id" json:"app_id"`
	Topic       string      `gorm:"topic" bson:"topic" json:"topic"`
	Data        string      `gorm:"data"  bson:"data"   json:"data"`                                         // 事件数据
	Meta        string      `gorm:"meta" bson:"meta"  json:"meta"`                                           // 上下文数据
	CreatedTime *times.Time `gorm:"created_time;index:,sort:desc,"  bson:"created_time"  json:"createdTime"` // 发生事件
}

type OutboxStatus string
