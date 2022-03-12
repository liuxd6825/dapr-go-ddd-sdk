package ddd

type ApplyEventRequest struct {
	TenantId      string      `json:"tenantId"`
	CommandId     string      `json:"commandId"`
	EventId       string      `json:"eventId"`
	EventData     interface{} `json:"eventData"`
	EventType     string      `json:"eventType"`
	EventRevision string      `json:"eventRevision"`
	AggregateId   string      `json:"AggregateId"`
	AggregateType string      `json:"aggregateType"`
	Metadata      interface{} `json:"metadata"`
	PubsubName    string      `json:"pubsubName"`
	Topic         string      `json:"topic"`
}

type ApplyEventsResponse struct {
}

type ExistAggregateResponse struct {
	IsExist bool `json:"isExist"`
}
