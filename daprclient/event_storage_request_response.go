package daprclient

import "encoding/json"

type ApplyEventRequest struct {
	TenantId      string            `json:"tenantId"`
	CommandId     string            `json:"commandId"`
	EventId       string            `json:"eventId"`
	EventData     interface{}       `json:"eventData"`
	EventType     string            `json:"eventType"`
	EventRevision string            `json:"eventRevision"`
	AggregateId   string            `json:"AggregateId"`
	AggregateType string            `json:"aggregateType"`
	Metadata      map[string]string `json:"metadata"`
	PubsubName    string            `json:"pubsubName"`
	Topic         string            `json:"topic"`
}

type ApplyEventsResponse struct {
}

type ExistAggregateResponse struct {
	IsExist bool `json:"isExist"`
}

type LoadEventsRequest struct {
	TenantId      string `json:"tenantId"`
	AggregateId   string `json:"aggregateId"`
	AggregateType string `json:"aggregateType"`
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
	SequenceNumber    uint64                 `json:"sequenceNumber"`
	Metadata          map[string]string      `json:"metadata"`
}

type EventRecord struct {
	EventId        string                 `json:"eventId"`
	EventData      map[string]interface{} `json:"eventData"`
	EventType      string                 `json:"eventType"`
	EventRevision  string                 `json:"eventRevision"`
	SequenceNumber uint64                 `json:"sequenceNumber"`
}

// NewEventRecordByJsonBytes 通过json反序列化EventRecord
func NewEventRecordByJsonBytes(data []byte) *EventRecordJsonMarshalResult {
	eventRecord := &EventRecord{}
	err := json.Unmarshal(data, &eventRecord)
	return &EventRecordJsonMarshalResult{
		eventRecord: eventRecord,
		err:         err,
	}
}

// EventRecordJsonMarshalResult EventRecord反序列化返回值
type EventRecordJsonMarshalResult struct {
	eventRecord *EventRecord
	err         error
}

func (r *EventRecordJsonMarshalResult) OnSuccess(doSuccess func(eventRecord *EventRecord) error) *EventRecordJsonMarshalResult {
	if r.err == nil {
		r.err = doSuccess(r.eventRecord)
	}
	return r
}
func (r *EventRecordJsonMarshalResult) OnError(doError func(err error)) *EventRecordJsonMarshalResult {
	if r.err != nil {
		doError(r.err)
	}
	return r
}

func (r *EventRecordJsonMarshalResult) GetError() error {
	return r.err
}

func (r *EventRecordJsonMarshalResult) GetEventRecord() *EventRecord {
	return r.eventRecord
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

type SaveSnapshotRequest struct {
	TenantId          string            `json:"tenantId"`
	AggregateId       string            `json:"AggregateId"`
	AggregateType     string            `json:"aggregateType"`
	AggregateData     interface{}       `json:"aggregateData"`
	AggregateRevision string            `json:"aggregateRevision"`
	Metadata          map[string]string `json:"metadata"`
	SequenceNumber    uint64            `json:"sequenceNumber"`
}

type SaveSnapshotResponse struct {
}
