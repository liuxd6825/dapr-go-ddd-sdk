package applog

import (
	"context"
	json2 "encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/liuxd6825/dapr-go-ddd-sdk/httpclient"
	"time"
)

var log Logger
var appId string

type DoAction func() error
type DoMethod func() (interface{}, error)

type Event interface {
	GetTenantId() string
	GetCommandId() string
	GetEventId() string
}

type EventHandler interface {
	GetClassName() string
}

func Init(httpClient *httpclient.HttpClient, aAppId string, level Level) {
	log = NewLogger(httpClient)
	log.SetLevel(level)
	appId = aAppId
}

func DoEventLog(ctx context.Context, handler EventHandler, event Event, funcName string, method DoAction) error {
	err := method()
	if err == nil {
		_, _ = InfoEvent(event.GetTenantId(), handler.GetClassName(), funcName, "success", event.GetEventId(), event.GetCommandId(), "")
	} else {
		_, _ = ErrorEvent(event.GetTenantId(), handler.GetClassName(), funcName, "error", event.GetEventId(), event.GetCommandId(), "")
	}
	return nil
}

func DoAppLog(ctx context.Context, info *LogInfo, method DoMethod) error {
	resp, err := method()
	_, _ = writeLog(ctx, info.TenantId, info.ClassName, info.FuncName, info.Level, info.Message)

	if log.GetLevel() <= INFO {
		bs, _ := json2.Marshal(resp)
		println(fmt.Sprintf("Result:%v \r\n", string(bs)))
	}

	if err != nil {
		_, _ = writeLog(ctx, info.TenantId, info.ClassName, info.FuncName, ERROR, err.Error())
	}
	return err
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

func GetEventLogByAppIdAndCommandId(tenantId, appId, commandId string) (*[]EventLogDto, error) {
	return getEventLogByAppIdAndCommandId(context.Background(), tenantId, appId, commandId)
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

func getEventLogByAppIdAndCommandId(ctx context.Context, tenantId, appId, commandId string) (*[]EventLogDto, error) {
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
	if level < log.GetLevel() {
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

	fmt.Printf("[%s] appid=%s; class=%s; func=%s; msg=%s; status=%t; time=%s;\n", req.Level, req.AppId, req.Class, req.Func, req.Message, req.Status, req.Time)

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
