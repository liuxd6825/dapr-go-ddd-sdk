package applog

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/sirupsen/logrus"
)

const (
	ApiWriteEventLog          = "v1.0/app-logger/%v/event-log/create"
	ApiUpdateEventLog         = "v1.0/app-logger/%v/event-log/update"
	ApiGetEventLogByCommandId = "v1.0/app-logger/%v/event-log/tenant-id/%s/app-id/%s/command-id/%s"

	ApiWriteAppLog   = "v1.0/app-logger/%v/app-log/create"
	ApiUpdateAppLog  = "v1.0/app-logger/%v/app-log/update"
	ApiGetAppLogById = "v1.0/app-logger/%v/app-log/tenant-id/%s/id/%s"
)

type Level = logrus.Level

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

type WriteEventLogRequest = dapr.WriteEventLogRequest

type WriteEventLogResponse = dapr.WriteEventLogResponse

type UpdateEventLogRequest = dapr.UpdateEventLogRequest

type UpdateEventLogResponse = dapr.UpdateEventLogResponse

type GetEventLogByCommandIdRequest = dapr.GetEventLogByCommandIdRequest

type GetEventLogByCommandIdResponse = dapr.GetEventLogByCommandIdResponse

type EventLogDto = dapr.EventLogDto

type WriteAppLogRequest = dapr.WriteAppLogRequest

type WriteAppLogResponse = dapr.WriteAppLogResponse

type UpdateAppLogRequest = dapr.UpdateAppLogRequest

type UpdateAppLogResponse = dapr.UpdateAppLogResponse

type GetAppLogByIdRequest = dapr.GetAppLogByIdRequest

type GetAppLogByIdResponse = dapr.GetAppLogByIdResponse
