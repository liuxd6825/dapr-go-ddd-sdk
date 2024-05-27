package event

import (
	"fmt"
	"github.com/liuxu6825/duxm-graph-service/pkg/cmd-service/domain/graph/field"
	"time"
)

//
// AccountNodeCreateEvent
// @Description: 领域事件
//
type {{.Name}} struct {
	EventId     string                        `json:"eventId" validate:"required"`   // 领域事件ID
	CommandId   string                        `json:"commandId" validate:"required"` // 关联命令ID
	CreatedTime time.Time                     `json:"time" validate:"required"`      // 事件创建时间
	Version     string                        `json:"version" validate:"required"`   // 事件版本
	EventType   string                        `json:"eventType" validate:"required"` // 事件类型
	Data        field.{{.Name}}Fields `json:"data" validate:"required"`      // 业务字段项
}

func New{{.Name}}(commandId string) *{{.Name}}Event {
	return &AccountNodeCreateEvent{
		EventId:     fmt.Sprintf("%s(%s)", commandId, {{.Name}}EventType.String()),
		CommandId:   commandId,
		Version:     "v1.0",
		EventType:   {{.Name}}EventType.String(),
		CreatedTime: time.Now(),
	}
}

func (e *{{.Name}}) GetEventId() string {
	return e.EventId
}

func (e *{{.Name}}) GetEventType() string {
	return e.EventType
}

func (e *{{.Name}}) GetEventVersion() string {
	return e.Version
}

func (e *{{.Name}}) GetCommandId() string {
	return e.CommandId
}

func (e *{{.Name}}) GetTenantId() string {
	return e.Data.TenantId
}

func (e *{{.Name}}) GetAggregateId() string {
	return e.Data.GraphId
}

func (e *{{.Name}}) GetCreatedTime() time.Time {
	return e.CreatedTime
}

func (e *{{.Name}}) GetData() interface{} {
	return e.Data
}
