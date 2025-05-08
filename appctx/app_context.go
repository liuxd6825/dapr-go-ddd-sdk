package appctx

import (
	"context"
)

type appKey struct {
}

type AppInfo interface {
	GetAppId() string
	GetAppName() string
}

type app struct {
	appId   string
	appName string
}

var appCtxKey = appKey{}

func newAppContext(parent context.Context, appId, appName string) context.Context {
	newCtx := context.WithValue(parent, appCtxKey, &app{
		appId:   appId,
		appName: appName,
	})
	return newCtx
}

func GetAppInfo(ctx context.Context) (AppInfo, bool) {
	if ctx == nil {
		return nil, false
	}
	val := ctx.Value(appCtxKey)
	if val == nil {
		return nil, false
	}
	appInfo, ok := val.(AppInfo)
	if !ok {
		return nil, false
	}
	return appInfo, true
}

func (a *app) GetAppId() string {
	return a.appId
}

func (a *app) GetAppName() string {
	return a.appName
}
