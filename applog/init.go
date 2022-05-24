package applog

import (
	"context"
	json2 "encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
	"runtime"
	"strings"
	"time"

	logrus "github.com/sirupsen/logrus"
)

var log Logger
var appId string

type DoAction func() error
type DoFunc func() (interface{}, error)

type Event interface {
	GetTenantId() string
	GetCommandId() string
	GetEventId() string
}

//
// Init
// @Description: 初始化日期
// @param daprClient DaprDddClient
// @param aAppId Darp Appliation Id
// @param level 日志级别
//
func Init(daprClient daprclient.DaprDddClient, aAppId string, level Level) {
	log = NewLogger(daprClient)
	log.SetLevel(level)
	appId = aAppId
}

//
// DoEventLog
// @Description: 执行日志记录
// @param ctx Context
// @param handler 事件Handler
// @param event 事件
// @param funcName 方法名称
// @param method 执行函数
// @return error 错误
//
func DoEventLog(ctx context.Context, structNameFunc func() string, event Event, funcName string, fun DoAction) error {
	err := fun()
	if err == nil {
		_, _ = InfoEvent(event.GetTenantId(), structNameFunc(), funcName, "success", event.GetEventId(), event.GetCommandId(), "")
	} else {
		_, _ = ErrorEvent(event.GetTenantId(), structNameFunc(), funcName, "error", event.GetEventId(), event.GetCommandId(), "")
	}
	return nil
}

//
// DoAppLog
// @Description:
// @param ctx
// @param info
// @param fun
// @return error
//
func DoAppLog(ctx context.Context, info *LogInfo, fun DoFunc) error {
	if err := assert.NotNil(info, assert.WidthOptionsError("info is nil")); err != nil {
		return err
	}
	if err := assert.NotNil(ctx, assert.WidthOptionsError("ctx is nil")); err != nil {
		return err
	}
	if err := assert.NotNil(fun, assert.WidthOptionsError("fun is nil")); err != nil {
		return err
	}
	resp, err := fun()
	_, _ = writeAppLog(ctx, info.TenantId, info.ClassName, info.FuncName, info.Level, info.Message)

	if log.GetLevel() <= INFO {
		bs, _ := json2.Marshal(resp)
		println(fmt.Sprintf("Result:%v \r\n", string(bs)))
	}

	if err != nil {
		_, _ = writeAppLog(ctx, info.TenantId, info.ClassName, info.FuncName, ERROR, err.Error())
	}
	return err
}

//
// Debug
// @Description:  写调试级日志
// @param tenantId
// @param className
// @param funcName
// @param message
// @return string
// @return error
//
func Debug(tenantId, className, funcName, message string) (string, error) {
	return writeAppLog(context.Background(), tenantId, className, funcName, DEBUG, message)
}

//
// Info
// @Description: 写信息级日志
// @param tenantId
// @param className
// @param funcName
// @param message
// @return string
// @return error
//
func Info(tenantId, className, funcName, message string) (string, error) {
	return writeAppLog(context.Background(), tenantId, className, funcName, INFO, message)
}

//
// Warn
//  @Description: 写警告级日志
//  @param tenantId
//  @param className
//  @param funcName
//  @param message
//  @return string
//  @return error
//
func Warn(tenantId, className, funcName, message string) (string, error) {
	return writeAppLog(context.Background(), tenantId, className, funcName, WARN, message)
}

//
// Error
//  @Description: 写错误级日志
//  @param tenantId
//  @param className
//  @param funcName
//  @param message
//  @return string
//  @return error
//
func Error(tenantId, className, funcName, message string) (string, error) {
	return writeAppLog(context.Background(), tenantId, className, funcName, ERROR, message)
}

//
// Fatal
//  @Description:  写致命级日志
//  @param tenantId
//  @param className
//  @param funcName
//  @param message
//  @return string
//  @return error
//
func Fatal(tenantId, className, funcName, message string) (string, error) {
	return writeAppLog(context.Background(), tenantId, className, funcName, FATAL, message)
}

//
// InfoEvent
//  @Description: 写事件处理日志
//  @param tenantId
//  @param className
//  @param funcName
//  @param message
//  @param eventId
//  @param commandId
//  @param pubAppId
//  @return string
//  @return error
//
func InfoEvent(tenantId, structName, funcName, message, eventId, commandId, pubAppId string) (string, error) {
	return writeEventLog(context.Background(), tenantId, structName, funcName, INFO, message, eventId, commandId, pubAppId, false)
}

//
// ErrorEvent
//  @Description: 写错误级日志
//  @param tenantId
//  @param className
//  @param funcName
//  @param message
//  @param eventId
//  @param commandId
//  @param pubAppId
//  @return string
//  @return error
//
func ErrorEvent(tenantId, className, funcName, message, eventId, commandId, pubAppId string) (string, error) {
	return writeEventLog(context.Background(), tenantId, className, funcName, ERROR, message, eventId, commandId, pubAppId, false)
}

//
// GetEventLogByAppIdAndCommandId
//  @Description:  按CommandId获取日志
//  @param tenantId
//  @param appId
//  @param commandId
//  @return *[]EventLogDto
//  @return error
//
func GetEventLogByAppIdAndCommandId(tenantId, appId, commandId string) (*[]EventLogDto, error) {
	return getEventLogByAppIdAndCommandId(context.Background(), tenantId, appId, commandId)
}

//
//  writeEventLog
//  @Description: 写事件日志
//  @param ctx
//  @param tenantId
//  @param structName 被记录的结构名称
//  @param funcName 被记录的方法名称
//  @param level 日志级别
//  @param message
//  @param eventId
//  @param commandId
//  @param pubAppId
//  @param status
//  @return string
//  @return error
//
func writeEventLog(ctx context.Context, tenantId, structName, funcName string, level Level, message, eventId, commandId, pubAppId string, status bool) (string, error) {
	uid := uuid.New().String()
	timeNow := time.Now()
	req := &WriteEventLogRequest{
		Id:        uid,
		TenantId:  tenantId,
		AppId:     appId,
		Class:     structName,
		Func:      funcName,
		Level:     level.ToString(),
		Time:      &timeNow,
		Status:    true,
		Message:   message,
		EventId:   eventId,
		CommandId: commandId,
		PubAppId:  pubAppId,
	}

	logrus.WithFields(logrus.Fields{
		"id":        uid,
		"tenantId":  tenantId,
		"appId":     appId,
		"class":     structName,
		"func":      funcName,
		"level":     level.ToString(),
		"time":      &timeNow,
		"status":    true,
		"message":   message,
		"eventId":   eventId,
		"commandId": commandId,
		"pubAppId":  pubAppId,
	}).Infoln(fmt.Sprintf("EVENT LOG %s", structName))

	_, err := log.WriteEventLog(ctx, req)
	return req.Id, err
}

//
//  updateEventLog
//  @Description: 更新事件日志
//  @param ctx
//  @param tenantId
//  @param id
//  @param className
//  @param funcName
//  @param level
//  @param message
//  @param eventId
//  @param commandId
//  @param pubAppId
//  @return string
//  @return error
//
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

//
//  getEventLogByAppIdAndCommandId
//  @Description: 按AppId与CommandId获取事件日志
//  @param ctx
//  @param tenantId
//  @param appId
//  @param commandId
//  @return *[]EventLogDto
//  @return error
//
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

//
//  writeAppLog
//  @Description: 写日志
//  @param ctx  上下文
//  @param tenantId 租户id
//  @param id 日志id
//  @param className 结构名称
//  @param funcName 方法名称
//  @param level 日志级别
//  @param message 日志内容
//  @return string 日志id
//  @return error 错误
//
func writeAppLog(ctx context.Context, tenantId, className, funcName string, level Level, message string) (string, error) {
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

	logrus.WithFields(logrus.Fields{
		"id":       uid,
		"tenantId": tenantId,
		"appId":    appId,
		"class":    className,
		"func":     funcName,
		"level":    level.ToString(),
		"time":     &timeNow,
		"status":   true,
		"message":  message,
	}).Infoln("AppLog")

	fmt.Printf("[%s] appid=%s; class=%s; func=%s; msg=%s; status=%t; time=%s;\n", strings.ToUpper(req.Level), req.AppId, req.Class, req.Func, req.Message, req.Status, req.Time)

	_, err := log.WriteAppLog(ctx, req)
	return req.Id, err
}

//
//  updateAppLog
//  @Description: 更新日志
//  @param ctx  上下文
//  @param tenantId 租户id
//  @param id 日志id
//  @param className 结构名称
//  @param funcName 方法名称
//  @param level 日志级别
//  @param message 日志内容
//  @return *UpdateAppLogRequest 更新结果
//  @return error 错误
//
func updateAppLog(ctx context.Context, tenantId, id, className, funcName string, level Level, message string) (*UpdateAppLogRequest, error) {
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

// RunFuncName
//  @Description: 取得当前运行的方法名称
//  @param skip 获取第几级方法名称
//  @return string 方法名称
//
func RunFuncName(skip int) string {
	pc := make([]uintptr, 1)
	runtime.Callers(skip+1, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
