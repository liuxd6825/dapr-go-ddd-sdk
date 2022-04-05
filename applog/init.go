package applog

import (
	"context"
	"github.com/google/uuid"
	"github.com/liuxd6825/dapr-go-ddd-sdk/httpclient"
	"time"
)

type DoAction func() error

type Event interface {
	GetTenantId() string
	GetCommandId() string
	GetEventId() string
}

type EventHandler interface {
	GetClassName() string
}

var log Logger
var appId string

func Init(httpClient *httpclient.HttpClient, aAppId string) {
	log = NewLogger(httpClient)
	appId = aAppId
}

func DoEvent(handler EventHandler, event Event, funcName string, actionFunc DoAction) error {
	err := actionFunc()
	if err == nil {
		_, _ = InfoEvent(event.GetTenantId(), handler.GetClassName(), funcName, "success", event.GetEventId(), event.GetCommandId(), "")
	} else {
		_, _ = ErrorEvent(event.GetTenantId(), handler.GetClassName(), funcName, "error", event.GetEventId(), event.GetCommandId(), "")
	}
	return nil
}

func Debug(tenantId, className, funcName, message string) (string, error) {
	return writeLog(context.Background(), tenantId, className, funcName, DEBUG, message)
}

func Info(tenantId, className, funcName, message string) (string, error) {
	return writeLog(context.Background(), tenantId, className, funcName, INFO, message)
}

func Warn(tenantId, className, funcName, message string) (string, error) {
	return writeLog(context.Background(), tenantId, className, funcName, WARN, message)
}

func Error(tenantId, className, funcName, message string) (string, error) {
	return writeLog(context.Background(), tenantId, className, funcName, ERROR, message)
}

func Fatal(tenantId, className, funcName, message string) (string, error) {
	return writeLog(context.Background(), tenantId, className, funcName, FATAL, message)
}

func InfoEvent(tenantId, className, funcName, message, eventId, commandId, pubAppId string) (string, error) {
	return writeEventLog(context.Background(), tenantId, className, funcName, INFO, message, eventId, commandId, pubAppId, false)
}

func ErrorEvent(tenantId, className, funcName, message, eventId, commandId, pubAppId string) (string, error) {
	return writeEventLog(context.Background(), tenantId, className, funcName, ERROR, message, eventId, commandId, pubAppId, false)
}

func GetEventInfo(tenantId, commandId string) (*[]EventLogDto, error) {
	return getEventLogByCommandId(context.Background(), tenantId, commandId)
}

func writeEventLog(ctx context.Context, tenantId, className, funcName string, level Level, message, eventId, commandId, pubAppId string, status bool) (string, error) {
	uid := uuid.New().String()
	timeNow := time.Now()
	req := &WriteEventLogRequest{
		Id:        uid,
		TenantId:  tenantId,
		AppId:     appId,
		Class:     className,
		Func:      funcName,
		Level:     level.ToString(),
		Time:      &timeNow,
		Status:    true,
		Message:   message,
		EventId:   eventId,
		CommandId: commandId,
		PubAppId:  pubAppId,
	}
	_, err := log.WriteEventLog(ctx, req)
	return req.Id, err
}

func updateEventLog(ctx context.Context, tenantId, id, className, funcName string, level Level, message, eventId, commandId, pubAppId string) (string, error) {
	timeNow := time.Now()
	req := &WriteEventLogRequest{
		Id:       id,
		TenantId: tenantId,
		AppId:    appId,
		Class:    className,
		Func:     funcName,
		Level:    level.ToString(),
		Time:     &timeNow,
		Status:   true,
		Message:  message,

		EventId:   eventId,
		CommandId: commandId,
		PubAppId:  pubAppId,
	}
	_, err := log.WriteEventLog(ctx, req)
	return req.Id, err
}

func getEventLogByCommandId(ctx context.Context, tenantId, commandId string) (*[]EventLogDto, error) {
	req := &GetEventLogByCommandIdRequest{
		TenantId:  tenantId,
		AppId:     appId,
		CommandId: commandId,
	}
	resp, err := log.GetEventLogByCommandId(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	return resp.Data, nil
}

func writeLog(ctx context.Context, tenantId, className, funcName string, level Level, message string) (string, error) {
	if level > log.GetLevel() {
		return "", nil
	}
	uid := uuid.New().String()
	timeNow := time.Now()
	req := &WriteAppLogRequest{
		Id:       uid,
		TenantId: tenantId,
		AppId:    appId,
		Class:    className,
		Func:     funcName,
		Level:    level.ToString(),
		Time:     &timeNow,
		Status:   true,
		Message:  message,
	}
	_, err := log.WriteAppLog(ctx, req)
	return req.Id, err
}

func updateLog(ctx context.Context, tenantId, id, className, funcName string, level Level, message string) (*UpdateAppLogRequest, error) {
	timeNow := time.Now()
	req := &UpdateAppLogRequest{
		Id:       id,
		TenantId: tenantId,
		AppId:    appId,
		Class:    className,
		Func:     funcName,
		Level:    level.ToString(),
		Time:     &timeNow,
		Status:   true,
		Message:  message,
	}
	_, err := log.UpdateAppLog(ctx, req)
	return req, err
}
