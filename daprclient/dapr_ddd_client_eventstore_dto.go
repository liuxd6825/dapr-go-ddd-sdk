package daprclient

import (
	"context"
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/appctx"
	"time"
)

type RequestHeader struct {
	Values map[string]string `json:"values"`
}

type CommitRequest struct {
	Headers   RequestHeader `json:"headers"`
	CompName  string        `json:"compName"`
	TenantId  string        `json:"tenantId"`
	SessionId string        `json:"sessionId"`
}

type CommitResponse struct {
	Headers *ResponseHeaders `json:"headers"`
}

type RollbackRequest struct {
	Headers   RequestHeader `json:"headers"`
	CompName  string        `json:"compName"`
	TenantId  string        `json:"tenantId"`
	SessionId string        `json:"sessionId"`
}

type RollbackResponse struct {
	Headers *ResponseHeaders `json:"headers"`
}

type ApplyEventRequest struct {
	Headers       RequestHeader `json:"headers"`
	CompName      string        `json:"compName"`
	SessionId     string        `json:"sessionId"`
	TenantId      string        `json:"tenantId"`
	AggregateId   string        `json:"aggregateId"`
	AggregateType string        `json:"aggregateType"`
	Events        []*EventDto   `json:"events"`
}

type ApplyEventResponse struct {
	Headers *ResponseHeaders `json:"headers"`
}

type CreateEventRequest struct {
	Headers       RequestHeader `json:"headers"`
	CompName      string        `json:"compName"`
	SessionId     string        `json:"sessionId"`
	TenantId      string        `json:"tenantId"`
	AggregateId   string        `json:"aggregateId"`
	AggregateType string        `json:"aggregateType"`
	Events        []*EventDto   `json:"events"`
}

type CreateEventResponse struct {
	Headers *ResponseHeaders `json:"headers"`
}

type DeleteEventRequest struct {
	Headers       RequestHeader `json:"headers"`
	CompName      string        `json:"compName"`
	SessionId     string        `json:"sessionId"`
	TenantId      string        `json:"tenantId"`
	AggregateId   string        `json:"aggregateId"`
	AggregateType string        `json:"aggregateType"`
	Event         *EventDto     `json:"event"`
}

type DeleteEventResponse struct {
	Headers *ResponseHeaders `json:"headers"`
}

type EventDto struct {
	ApplyType    string            `json:"applyType"`
	EventId      string            `json:"eventId"`
	CommandId    string            `json:"commandId"`
	EventData    interface{}       `json:"eventData"`
	EventType    string            `json:"eventType"`
	EventVersion string            `json:"eventVersion"`
	Metadata     map[string]string `json:"metadata"`
	PubsubName   string            `json:"pubsubName"`
	Topic        string            `json:"topic"`
	Relations    map[string]string `json:"relations"` // 聚合关系
	IsSourcing   bool              `json:"isSourcing"`
}

type ExistAggregateResponse struct {
	Headers *ResponseHeaders `json:"headers"`
	IsExist bool             `json:"isExist"`
}

type LoadEventsRequest struct {
	Headers       RequestHeader `json:"headers"`
	CompName      string        `json:"compName"`
	TenantId      string        `json:"tenantId"`
	AggregateId   string        `json:"aggregateId"`
	AggregateType string        `json:"aggregateType"`
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
	TenantId          string                 `json:"tenantId"`
	AggregateData     map[string]interface{} `json:"aggregateData"`
	AggregateRevision string                 `json:"aggregateRevision"`
	SequenceNumber    uint64                 `json:"sequenceNumber"`
	Metadata          map[string]string      `json:"metadata"`
}

type EventRecord struct {
	TenantId       string                 `json:"tenantId"`
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

type SaveSnapshotRequest struct {
	Headers          RequestHeader     `json:"headers"`
	CompName         string            `json:"compName"`
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
	Headers       RequestHeader `json:"headers"`
	CompName      string        `json:"compName"`
	TenantId      string        `json:"tenantId"`
	AggregateType string        `json:"aggregateType"`
	Filter        string        `json:"filter"`
	Sort          string        `json:"sort"`
	PageNum       uint64        `json:"pageNum"`
	PageSize      uint64        `json:"pageSize"`
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
	RelName     string `json:"relName"`
	RelValue    string `json:"relValue"`
}

type ResponseHeaders struct {
	Values  map[string]string `json:"values"`
	Status  ResponseStatus    `json:"status"`
	Message string            `json:"message"`
}

type GetEventsRequest struct {
	Headers       RequestHeader `json:"headers"`
	CompName      string        `json:"compName"`
	TenantId      string        `json:"tenantId"`
	AggregateType string        `json:"aggregateType"`
	Filter        string        `json:"filter"`
	Sort          string        `json:"sort"`
	PageNum       uint64        `json:"pageNum"`
	PageSize      uint64        `json:"pageSize"`
}

type GetEventsResponse struct {
	Headers    *ResponseHeaders `json:"headers"`
	Data       []*GetEventsItem `json:"data"`
	TotalRows  uint64           `json:"totalRows"`
	TotalPages uint64           `json:"totalPages"`
	PageNum    uint64           `json:"pageNum"`
	PageSize   uint64           `json:"pageSize"`
	Filter     string           `json:"filter"`
	Sort       string           `json:"sort"`
	Error      string           `json:"error"`
	IsFound    bool             `json:"isFound"`
}

type GetEventsItem struct {
	EventId      string
	CommandId    string
	EventData    map[string]interface{}
	EventType    string
	EventVersion string
	EventTime    *time.Time
	PubsubName   string
	Topic        string
	Metadata     map[string]string
	IsSourcing   bool
}

type ResponseStatus int64

const (
	ResponseStatusSuccess        ResponseStatus = iota // 执行成功
	ResponseStatusError                                // 执行错误
	ResponseStatusEventDuplicate                       // 事件已经存在，被重复执行
)

func NewResponseHeadersNil() *ResponseHeaders {
	return &ResponseHeaders{}
}

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

func (r *ResponseHeaders) SetError(err error) {
	if err != nil {
		r.Message = err.Error()
	}
}

func (r *ResponseHeaders) SetMessage(v string) {
	r.Message = v
}

func (r *ResponseHeaders) SetStatus(v int32) {
	r.Status = ResponseStatus(v)
}

func (r *ResponseHeaders) SetValues(v map[string]string) {
	r.Values = v
}

func (r *EventRecordJsonMarshalResult) OnSuccess(ctx context.Context, doSuccess func(ctx context.Context, record *EventRecord) error) *EventRecordJsonMarshalResult {
	if r.err == nil {
		newCtx := appctx.NewTenantContext(ctx, r.eventRecord.TenantId)
		r.err = doSuccess(newCtx, r.eventRecord)
	}
	return r
}

func (r *EventRecordJsonMarshalResult) OnError(doError func(err error, eventRecord *EventRecord)) *EventRecordJsonMarshalResult {
	if r.err != nil {
		doError(r.err, r.eventRecord)
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
