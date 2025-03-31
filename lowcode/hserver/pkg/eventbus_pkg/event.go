package eventbus_pkg

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"time"
)

type ApplyEventOptions = ddd.ApplyEventOptions
type DomainEvent struct {
	CommandId   string         `json:"commandId" validate:"required"` // 关联命令ID
	EventId     string         `json:"eventId" validate:"required"`   // 领域事件ID
	EventType   string         `json:"eventType" validate:"required"` // 事件类型
	EventVer    string         `json:"eventVer" validate:"required"`
	AggId       string         `json:"aggId" validate:"required"`
	AggType     string         `json:"aggType" validate:"required"`
	CreatedTime time.Time      `json:"time" validate:"required"` // 事件创建时间
	Data        map[string]any `json:"data" validate:"required"` // 业务字段项
}

func (d *DomainEvent) GetTenantId() string {
	tenantId, ok := d.Data["tenantId"]
	if !ok {
		return ""
	}
	return tenantId.(string)
}

func (d *DomainEvent) GetCommandId() string {
	return d.CommandId
}

func (d *DomainEvent) GetEventId() string {
	return d.EventId
}

func (d *DomainEvent) GetEventType() string {
	return d.EventType
}

func (d *DomainEvent) GetEventVer() string {
	return d.EventVer
}

func (d *DomainEvent) GetAggId() string {
	return d.AggId
}

func (d *DomainEvent) GetAggType() string {
	return d.AggType
}

func (d *DomainEvent) GetCreatedTime() time.Time {
	return d.CreatedTime
}

func (d *DomainEvent) GetData() interface{} {
	return d.Data
}
