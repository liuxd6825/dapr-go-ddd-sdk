package dbevent

import (
	"time"
)

type Outbox struct {
	Id          string              `gorm:"id;primaryKey" json:"id"`
	TenantId    string              `gorm:"tenant_id" json:"tenantId,omitempty"`                // 租户ID
	AppId       string              `gorm:"app_name" json:"appName,omitempty"`                  // 应用ID
	EventId     string              `gorm:"event_id" json:"eventId,omitempty"`                  // 事件ID
	EventType   string              `gorm:"event_type" json:"eventType,omitempty"`              // 事件类型
	EventVer    string              `gorm:"event_ver" json:"eventVer,omitempty"`                // 事件版本
	CommandId   string              `gorm:"command_id" json:"commandId,omitempty"`              // 命令ID
	AggId       string              `gorm:"agg_id" json:"aggId,omitempty"`                      // 聚合ID
	AggType     string              `gorm:"agg_type" json:"aggType,omitempty"`                  // 聚合类型
	CreatedTime *time.Time          `gorm:"index:,sort:desc,created_time" json:"createdTime"`   // 发生事件
	Data        any                 `gorm:"type:json;serializer:json" json:"data,omitempty"`    // 事件数据
	Metadata    map[string][]string `gorm:"type:json;serializer:json" json:"context,omitempty"` // 上下文数据
}

type OutboxStatus string
