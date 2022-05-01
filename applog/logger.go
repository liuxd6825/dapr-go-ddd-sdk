package applog

import (
	"context"
	"errors"
	"fmt"
	"strings"
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

type Level int

const (
	TRACE Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	ALL
	OFF
)

func (l Level) ToString() string {
	switch l {
	case TRACE:
		return "trace"
	case DEBUG:
		return "debug"
	case INFO:
		return "info"
	case WARN:
		return "warn"
	case ERROR:
		return "error"
	case FATAL:
		return "fatal"
	case ALL:
		return "all"
	case OFF:
		return "off"
	}
	return "none"
}
func NewLevel(name string) (Level, error) {
	v := strings.ToLower(name)
	switch v {
	case "trace":
		return TRACE, nil
	case "debug":
		return DEBUG, nil
	case "info":
		return INFO, nil
	case "warn":
		return WARN, nil
	case "error":
		return ERROR, nil
	case "fatal":
		return FATAL, nil
	case "all":
		return ALL, nil
	case "off":
		return OFF, nil
	}
	return DEBUG, errors.New(fmt.Sprintf("%s as applog.Level error", name))
}

type Logger interface {
	WriteEventLog(ctx context.Context, req *WriteEventLogRequest) (*WriteEventLogResponse, error)
	UpdateEventLog(ctx context.Context, req *UpdateEventLogRequest) (*UpdateEventLogResponse, error)
	GetEventLogByCommandId(ctx context.Context, req *GetEventLogByCommandIdRequest) (*GetEventLogByCommandIdResponse, error)

	WriteAppLog(ctx context.Context, req *WriteAppLogRequest) (*WriteAppLogResponse, error)
	UpdateAppLog(ctx context.Context, req *UpdateAppLogRequest) (*UpdateAppLogResponse, error)
	GetAppLogById(ctx context.Context, req *GetAppLogByIdRequest) (*GetAppLogByIdResponse, error)

	SetLevel(level Level)
	GetLevel() Level
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

	PubAppId string `json:"pubAppId"`
	EventId  string `json:"eventId"`
	// EventType string `json:"eventType"`
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
