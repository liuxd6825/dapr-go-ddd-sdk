package daprclient

import "encoding/json"

type ApplyEventRequest struct {
	TenantId      string      `json:"tenantId"`
	AggregateId   string      `json:"aggregateId"`
	AggregateType string      `json:"aggregateType"`
	Events        []*EventDto `json:"events"`
}

type ApplyEventResponse struct {
	Headers *ResponseHeaders `json:"headers"`
}

type CreateEventRequest struct {
	TenantId      string      `json:"tenantId"`
	AggregateId   string      `json:"aggregateId"`
	AggregateType string      `json:"aggregateType"`
	Events        []*EventDto `json:"events"`
}

type CreateEventResponse struct {
	Headers *ResponseHeaders `json:"headers"`
}

type DeleteEventRequest struct {
	TenantId      string    `json:"tenantId"`
	AggregateId   string    `json:"aggregateId"`
	AggregateType string    `json:"aggregateType"`
	Event         *EventDto `json:"event"`
}

type DeleteEventResponse struct {
	Headers *ResponseHeaders `json:"headers"`
}

type EventDto struct {
	EventId      string            `json:"eventId"`
	CommandId    string            `json:"commandId"`
	EventData    interface{}       `json:"eventData"`
	EventType    string            `json:"eventType"`
	EventVersion string            `json:"eventVersion"`
	Metadata     map[string]string `json:"metadata"`
	PubsubName   string            `json:"pubsubName"`
	Topic        string            `json:"topic"`
	Relations    map[string]string `json:"relations"` // 聚合关系
}

type ExistAggregateResponse struct {
	Headers *ResponseHeaders `json:"headers"`
	IsExist bool             `json:"isExist"`
}

type LoadEventsRequest struct {
	TenantId      string `json:"tenantId"`
	AggregateId   string `json:"aggregateId"`
	AggregateType string `json:"aggregateType"`
}

type LoadEventsResponse struct {
	Headers       *ResponseHeaders `json:"headers"`
	TenantId      string           `json:"tenantId"`
	AggregateId   string           `json:"aggregateId"`
	AggregateType string           `json:"aggregateType"`
	Snapshot      *Snapshot        `json:"snapshot"`
	EventRecords  *[]EventRecord   `json:"events"`
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
	EventVersion   string                 `json:"eventVersion"`
	SequenceNumber uint64                 `json:"sequenceNumber"`
}

// NewEventRecordByJsonBytes 通过json反序列化EventRecord
func NewEventRecordByJsonBytes(data []byte) *EventRecordJsonMarshalResult {
	eventRecord := &EventRecord{}
	err := json.Unmarshal(data, eventRecord)
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
	TenantId         string            `json:"tenantId"`
	AggregateId      string            `json:"AggregateId"`
	AggregateType    string            `json:"aggregateType"`
	AggregateData    interface{}       `json:"aggregateData"`
	AggregateVersion string            `json:"aggregateVersion"`
	Metadata         map[string]string `json:"metadata"`
	SequenceNumber   uint64            `json:"sequenceNumber"`
}

type SaveSnapshotResponse struct {
	Headers *ResponseHeaders `json:"headers"`
}

type GetRelationsRequest struct {
	TenantId      string `json:"tenantId"`
	AggregateType string `json:"aggregateType"`
	Filter        string `json:"filter"`
	Sort          string `json:"sort"`
	PageNum       uint64 `json:"pageNum"`
	PageSize      uint64 `json:"pageSize"`
}

type GetRelationsResponse struct {
	Headers    *ResponseHeaders `json:"headers"`
	Data       []*Relation      `json:"data"`
	TotalRows  uint64           `json:"totalRows"`
	TotalPages uint64           `json:"totalPages"`
	PageNum    uint64           `json:"pageNum"`
	PageSize   uint64           `json:"pageSize"`
	Filter     string           `json:"filter"`
	Sort       string           `json:"sort"`
	Error      string           `json:"error"`
	IsFound    bool             `json:"isFound"`
}

type Relation struct {
	Id          string `json:"id"`
	TenantId    string `json:"tenantId"`
	TableName   string `json:"tableName"`
	AggregateId string `json:"aggregateId"`
	IsDeleted   bool   `json:"isDeleted"`
	Items       map[string]string
}

type ResponseHeaders struct {
	Values  map[string]string `json:"values"`
	Status  ResponseStatus    `json:"status"`
	Message string            `json:"message"`
}

type ResponseStatus int32

const (
	ResponseStatusSuccess        ResponseStatus = iota // 执行成功
	ResponseStatusError                                // 执行错误
	ResponseStatusEventDuplicate                       // 事件已经存在，被重复执行
)

func NewResponseHeaders(status ResponseStatus, err error, values map[string]string) *ResponseHeaders {
	if values == nil {
		values = make(map[string]string)
	}
	if err != nil {
		return NewResponseHeadersError(err, values)
	}
	resp := &ResponseHeaders{
		Status:  status,
		Message: "Success",
		Values:  values,
	}
	return resp
}

func NewResponseHeadersError(err error, values map[string]string) *ResponseHeaders {
	if values == nil {
		values = make(map[string]string)
	}
	resp := &ResponseHeaders{
		Status:  ResponseStatusError,
		Message: err.Error(),
		Values:  values,
	}
	return resp
}
