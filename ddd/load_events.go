package ddd

import (
	"encoding/json"
)

type LoadEventsRequest struct {
	TenantId    string `json:"tenantId"`
	AggregateId string `json:"AggregateId"`
}

type LoadEventsResponse struct {
	TenantId      string         `json:"tenantId"`
	AggregateId   string         `json:"AggregateId"`
	AggregateType string         `json:"aggregateType"`
	Snapshot      *Snapshot      `json:"snapshot"`
	EventRecords  *[]EventRecord `json:"events"`
}

type Snapshot struct {
	AggregateData     map[string]interface{} `json:"aggregateData"`
	AggregateRevision string                 `json:"aggregateRevision"`
	SequenceNumber    int64                  `json:"sequenceNumber"`
	Metadata          map[string]interface{} `json:"metadata"`
}

type EventRecord struct {
	EventId        string                 `json:"eventId"`
	EventData      map[string]interface{} `json:"eventData"`
	EventType      string                 `json:"eventType"`
	EventRevision  string                 `json:"eventRevision"`
	SequenceNumber int64                  `json:"sequenceNumber"`
}

func (e *EventRecord) Marshal(domainEvent interface{}) error {
	jsonEvent, err := json.Marshal(e.EventData)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonEvent, domainEvent)
	if err != nil {
		return err
	}
	return nil
}

func (e *EventRecord) SetFields(key string, set func(value interface{})) {
	v, ok := e.EventData[key]
	if ok {
		set(v)
	}
}
