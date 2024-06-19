package common

import (
	"fmt"
	"time"
)

type EventType string

const (
	EventType_Create EventType = "create"
	EventType_Update EventType = "update"
	EventType_Delete EventType = "delete"
)

var tableEventTypes = []string{
	EventType_Create.String(),
	EventType_Update.String(),
	EventType_Delete.String(),
}

func TableEventTypes() []string {
	return tableEventTypes
}

func (e EventType) String() string {
	return string(e)
}

// GetEventType
//
//	@Description: 获取事件类型
//	@param appId 应用ID
//	@param aggName 聚合根类型名称
//	@param operateType 操作类型 增加，更新，删除等
//	@param eventVersion 版本号
//	@return string
func GetEventType(appId string, aggName string, operateType string) string {
	return fmt.Sprintf("%s.%s.%s", appId, aggName, operateType)
}

type Event struct {
	TenantId     string         `json:"tenantId,omitempty"`
	EventId      string         `json:"eventId,omitempty"`
	EventType    string         `json:"eventType,omitempty"`
	EventVersion string         `json:"eventVersion,omitempty"`
	CommandId    string         `json:"commandId,omitempty"`
	AggregateId  string         `json:"aggregateId,omitempty"`
	CreatedTime  time.Time      `json:"createdTime"`
	IsValidOnly  bool           `json:"isValidOnly,omitempty"`
	Data         map[string]any `json:"data,omitempty"`
}

func NewEvent() *Event {
	return &Event{}
}

func (c *Event) GetAggregateId() string {
	return c.AggregateId
}

func (c *Event) GetEventId() string {
	return c.EventId
}

func (c *Event) SetEventId(val string) {
	c.EventId = val
}

func (c *Event) GetTenantId() string {
	return c.TenantId
}

func (c *Event) SetTenantId(val string) {
	c.TenantId = val
}

func (c *Event) GetEventType() string {
	return c.EventType
}

func (c *Event) SetEventType(val string) {
	c.EventType = val
}

func (c *Event) GetData() any {
	return c.Data
}

func (c *Event) SetData(val map[string]any) {
	c.Data = val
}

func (c *Event) GetIsValidOnly() bool {
	return c.IsValidOnly
}

func (c *Event) SetIsValidOnly(val bool) {
	c.IsValidOnly = val
}

func (c *Event) GetCreatedTime() time.Time {
	return c.CreatedTime
}

func (c *Event) SetCreatedTime(val time.Time) {
	c.CreatedTime = val
}

func (c *Event) GetCommandId() string {
	return c.CommandId
}

func (c *Event) SetCommandId(val string) {
	c.CommandId = val
}

func (c *Event) GetEventVersion() string {
	return c.EventVersion
}
