package dbevent

import "time"

type Entity struct {
	TenantId     string         `json:"tenantId"`
	AggId        string         `json:"aggId"`
	AggType      string         `json:"aggType"`
	AggVersion   string         `json:"aggVersion"`
	CommandId    string         `json:"commandId"`
	EventId      string         `json:"eventId"`
	EventType    string         `json:"eventType"`
	EventVersion string         `json:"eventVersion"`
	CreatedTime  time.Time      `json:"createdTime"`
	IsValidOnly  bool           `json:"isValidOnly"`
	Data         map[string]any `json:"data"`
}
