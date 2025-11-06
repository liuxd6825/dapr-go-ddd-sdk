package dbevent

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types/times"
)

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
	Id          string         `json:"id" gorm:"id;primaryKey" bson:"id" json:"id"`
	TenantId    string         `json:"tenantId" gorm:"tenant_id" bson:"tenant_id" json:"tenantId"` // 租户ID
	AppId       string         `json:"appId" gorm:"app_id" bson:"app_id" json:"appId"`
	EventType   string         `json:"eventType" gorm:"event_type" bson:"event_type" json:"eventType"`
	Data        map[string]any `json:"data" gorm:"data;json" bson:"data" json:"data"`                                            // 事件数据
	Meta        map[string]any `json:"meta" gorm:"meta;json" bson:"meta" json:"meta"`                                            // 上下文数据
	CreatedTime *times.Time    `json:"createdTime" gorm:"created_time;index:,sort:desc," bson:"created_time" json:"createdTime"` // 发生事件
	//Version     string         `gorm:"version" bson:"version"  json:"version"`                                // 事件数据版本
}

type OutboxStatus string
