package logs_pkg

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

type Logs struct {
}

func New() *Logs {
	return &Logs{}
}

func (l *Logs) Trace(ctx context.Context, fields logs.Fields, args ...interface{}) {
	logs.Trace(ctx, getTenantId(ctx), fields, args...)
}

func (l *Logs) Print(ctx context.Context, fields logs.Fields) {
	logs.Print(ctx, getTenantId(ctx), fields)
}

func (l *Logs) Printf(ctx context.Context, fields logs.Fields, fmt string, args ...any) {
	logs.Printf(ctx, getTenantId(ctx), fields, fmt, args...)
}

func (l *Logs) Println(ctx context.Context, fields logs.Fields, fmt string, args ...any) {
	logs.Println(ctx, getTenantId(ctx), fields, fmt, args...)
}

func (l *Logs) Debug(ctx context.Context, fields logs.Fields) {
	logs.Debug(ctx, getTenantId(ctx), fields)
}

func (l *Logs) DebugMsg(ctx context.Context, args ...any) {
	logs.DebugMsg(ctx, getTenantId(ctx), args)
}

func (l *Logs) Debugf(ctx context.Context, fields logs.Fields, fmt string, args ...interface{}) {
	logs.Debugf(ctx, getTenantId(ctx), fields, fmt, args...)
}

func (l *Logs) DebugEvent(ctx context.Context, event logs.Event, funcName string) {
	logs.DebugEvent(ctx, event, funcName)
}

func (l *Logs) Debugfmt(ctx context.Context, fmt string, args ...interface{}) {
	logs.Debugfmt(ctx, getTenantId(ctx), fmt, args...)
}

func (l *Logs) Info(ctx context.Context, fields logs.Fields, args ...interface{}) {
	logs.Info(ctx, getTenantId(ctx), fields, args...)
}

func (l *Logs) Infof(ctx context.Context, fields logs.Fields, fmt string, args ...interface{}) {
	logs.Infof(ctx, getTenantId(ctx), fields, fmt, args...)
}

func (l *Logs) InfoMsg(ctx context.Context, args ...any) {
	logs.InfoMsg(ctx, getTenantId(ctx), args)
}

func (l *Logs) Warn(ctx context.Context, fields logs.Fields) {
	logs.Warn(ctx, getTenantId(ctx), fields)
}

func (l *Logs) Warnf(ctx context.Context, fields logs.Fields, fmt string, args ...interface{}) {
	logs.Warnf(ctx, getTenantId(ctx), fields, fmt, args...)
}

func (l *Logs) WarnMsg(ctx context.Context, args ...any) {
	logs.WarnMsg(ctx, getTenantId(ctx), args)
}

func (l *Logs) Warning(ctx context.Context, fields logs.Fields, args ...interface{}) {
	logs.Warning(ctx, getTenantId(ctx), fields, args...)
}

func (l *Logs) Error(ctx context.Context, fields logs.Fields, args ...interface{}) {
	logs.Error(ctx, getTenantId(ctx), fields, args...)
}

func (l *Logs) ErrorErr(ctx context.Context, err error) {
	logs.ErrorErr(ctx, getTenantId(ctx), err)
}

func (l *Logs) Errorf(ctx context.Context, fields logs.Fields, fmt string, args ...interface{}) {
	logs.Errorf(ctx, getTenantId(ctx), fields, fmt, args...)
}

func (l *Logs) Errorfmt(ctx context.Context, fmt string, args ...interface{}) {
	logs.Errorfmt(ctx, getTenantId(ctx), fmt, args...)
}

func (l *Logs) ErrorMsg(ctx context.Context, args ...any) {
	logs.ErrorMsg(ctx, getTenantId(ctx), args)
}

func (l *Logs) Panic(ctx context.Context, fields logs.Fields, args ...interface{}) {
	logs.Panic(ctx, getTenantId(ctx), fields, args...)
}

func (l *Logs) Panicf(ctx context.Context, fields logs.Fields, fmt string, args ...interface{}) {
	logs.Panicf(ctx, getTenantId(ctx), fields, fmt, args...)
}

func (l *Logs) PanicError(ctx context.Context, err error) {
	logs.PanicError(ctx, getTenantId(ctx), err)
}

func (l *Logs) Fatal(ctx context.Context, fields logs.Fields, args ...interface{}) {
	logs.Fatal(ctx, getTenantId(ctx), fields, args...)
}

func (l *Logs) FatalMsg(ctx context.Context, args ...interface{}) {
	logs.FatalMsg(ctx, getTenantId(ctx), args)
}

// DebugStart
//
//	@Description:
//	@param ctx
//	@param tenantId
//	@param fields
//	@param fun
//	@param format
//	@param args
//	@return err
func (l *Logs) DebugStart(ctx context.Context, fields logs.Fields, fun func() error) (err error) {
	return logs.DebugStart(ctx, getTenantId(ctx), fields, fun)
}

func getTenantId(ctx context.Context) string {
	authToken, ok := appctx.GetAuthToken(ctx)
	if !ok {
		return ""
	}
	return authToken.GetUser().GetTenantId()
}
