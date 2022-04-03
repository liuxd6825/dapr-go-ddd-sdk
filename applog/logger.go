package applog

import (
	"context"
	"time"
)

const (
	ApiWriteEventLog          = "v1.0/logger/event-log/create"
	ApiUpdateEventLog         = "v1.0/logger/event-log/update"
	ApiGetEventLogByCommandId = "v1.0/logger/event-log/tenant-id/%s/app-id/%s/command-id/%s"

	ApiWriteAppLog   = "v1.0/logger/app-log/create"
	ApiUpdateAppLog  = "v1.0/logger/app-log/update"
	ApiGetAppLogById = "v1.0/logger/app-log/tenant-id/%s/id/%s"
)

type Logger interface {
	WriteEventLog(ctx context.Context, req *WriteEventLogRequest) (*WriteEventLogResponse, error)
	UpdateEventLog(ctx context.Context, req *UpdateEventLogRequest) (*UpdateEventLogResponse, error)
	GetEventLogByCommandId(ctx context.Context, req *GetEventLogByCommandIdRequest) (*GetEventLogByCommandIdResponse, error)

	WriteAppLog(ctx context.Context, req *WriteAppLogRequest) (*WriteAppLogResponse, error)
	UpdateAppLog(ctx context.Context, req *UpdateAppLogRequest) (*UpdateAppLogResponse, error)
	GetAppLogById(ctx context.Context, req *GetAppLogByIdRequest) (*GetAppLogByIdResponse, error)
}
type WriteEventLogRequest struct {
	Id       string     `json:"id"`
	TenantId string     `json:"tenantId"`
	AppId    string     `json:"appId"`
	Class    string     `json:"class"`
	Func     string     `json:"func"`
	Level    string     `json:"level"`
	Time     *time.Time `json:"time"`
	Status   bool       `json:"status"`
	Message  string     `json:"message"`

	PubAppId  string `json:"pubAppId"`
	EventId   string `json:"eventId"`
	CommandId string `json:"commandId"`
}

type WriteEventLogResponse struct {
}

type UpdateEventLogRequest struct {
	Id       string     `json:"id"`
	TenantId string     `json:"tenantId"`
	AppId    string     `json:"appId"`
	Class    string     `json:"class"`
	Func     string     `json:"func"`
	Level    string     `json:"level"`
	Time     *time.Time `json:"time"`
	Status   bool       `json:"status"`
	Message  string     `json:"message"`

	PubAppId  string `json:"pubAppId"`
	EventId   string `json:"eventId"`
	CommandId string `json:"commandId"`
}

type UpdateEventLogResponse struct {
}

// GetLogByCommandId

type GetEventLogByCommandIdRequest struct {
	TenantId  string `json:"tenantId"`
	AppId     string `json:"appId"`
	CommandId string `json:"commandId"`
}

type GetEventLogByCommandIdResponse struct {
	Data *[]EventLogDto `json:"data"`
}

type EventLogDto struct {
	Id       string     `json:"id"`
	TenantId string     `json:"tenantId"`
	AppId    string     `json:"appId"`
	Class    string     `json:"class"`
	Func     string     `json:"func"`
	Level    string     `json:"level"`
	Time     *time.Time `json:"time"`
	Status   bool       `json:"status"`
	Message  string     `json:"message"`

	PubAppId  string `json:"pubAppId"`
	EventId   string `json:"eventId"`
	CommandId string `json:"commandId"`
}

//

type WriteAppLogRequest struct {
	Id       string     `json:"id"`
	TenantId string     `json:"tenantId"`
	AppId    string     `json:"appId"`
	Class    string     `json:"class"`
	Func     string     `json:"func"`
	Level    string     `json:"level"`
	Time     *time.Time `json:"time"`
	Status   bool       `json:"status"`
	Message  string     `json:"message"`
}

type WriteAppLogResponse struct {
}

type UpdateAppLogRequest struct {
	Id       string     `json:"id"`
	TenantId string     `json:"tenantId"`
	AppId    string     `json:"appId"`
	Class    string     `json:"class"`
	Func     string     `json:"func"`
	Level    string     `json:"level"`
	Time     *time.Time `json:"time"`
	Status   bool       `json:"status"`
	Message  string     `json:"message"`
}

type UpdateAppLogResponse struct {
}

// GetLogByCommandId

type GetAppLogByIdRequest struct {
	TenantId string `json:"tenantId"`
	Id       string `json:"id"`
}

type GetAppLogByIdResponse struct {
	Id       string     `json:"id"`
	TenantId string     `json:"tenantId"`
	AppId    string     `json:"appId"`
	Class    string     `json:"class"`
	Func     string     `json:"func"`
	Level    string     `json:"level"`
	Time     *time.Time `json:"time"`
	Status   bool       `json:"status"`
	Message  string     `json:"message"`
}
