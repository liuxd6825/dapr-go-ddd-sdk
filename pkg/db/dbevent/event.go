package dbevent

import "time"

type Event struct {
	TenantId    string    `json:"tenantId,omitempty"`
	TopicName   string    `json:"topicName,omitempty"`
	AppId       string    `json:"appId"`
	EventId     string    `json:"eventId,omitempty"`
	EventType   string    `json:"eventType,omitempty"`
	EventVer    string    `json:"eventVersion,omitempty"`
	CommandId   string    `json:"commandId,omitempty"`
	AggId       string    `json:"aggId,omitempty"`
	AggType     string    `json:"aggType,omitempty"`
	CreatedTime time.Time `json:"createdTime"`
	IsValidOnly bool      `json:"isValidOnly,omitempty"`
	Data        any       `json:"data,omitempty"`
}

func NewEvent() *Event {
	return &Event{}
}

func (c *Event) GetAggId() string {
	return c.AggId
}

func (c *Event) GetAggType() string {
	return c.AggType
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

func (c *Event) GetEventVer() string {
	return c.EventVer
}
