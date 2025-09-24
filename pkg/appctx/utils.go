package appctx

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"strings"
)

const Authorization = "Authorization"

func NewWebContext(parent context.Context, ictx iris.Context) (ctx context.Context) {
	header := ictx.Request().Header
	token := ictx.Request().Header.Get("authorization")
	ctx = parent
	if ctx1, err1 := NewAuthContext(ctx, getAuthorization(token, header)); err1 == nil {
		ctx = ctx1
	}
	ctx = NewHeaderContext(ctx, ictx.Request().Header)
	//ctx = NewTenantContext(ctx, "test")
	return ctx
}

func NewContext(parent context.Context, tenantId string, token string, header map[string][]string) (ctx context.Context) {
	ctx = NewTenantContext(parent, tenantId)
	ctx = NewHeaderContext(ctx, header)
	if ctx1, err1 := NewAuthContext(ctx, getAuthorization(token, header)); err1 == nil {
		ctx = ctx1
	}
	return ctx
}

func NewContextWidthAuthToken(parent context.Context, tenantId string, authToken *AuthTokenEntity, header map[string][]string) (ctx context.Context) {
	newCtx := NewTenantContext(parent, tenantId)
	newCtx = NewHeaderContext(newCtx, header)
	if ctx1, err1 := NewAuthContextUser(newCtx, authToken); err1 == nil {
		newCtx = ctx1
	}
	return newCtx
}

func SetContext(parent context.Context, tenantId string, token string, header map[string][]string) (ctx context.Context) {
	ctx = SetTenantContext(parent, tenantId)
	ctx = SetHeaderContext(ctx, header)
	if ctx1, err1 := SetAuthContext(ctx, getAuthorization(token, header)); err1 == nil {
		ctx = ctx1
	}
	return ctx
}

func getAuthorization(token string, header map[string][]string) string {
	if token != "" {
		return token
	}
	if val, ok := header[Authorization]; ok {
		for _, s := range val {
			s = strings.Trim(s, " ")
			if s != "" {
				return s
			}
		}
	}
	return ""
}

func NewMapWithContext(ctx context.Context) map[string]any {
	data := make(map[string]any)
	if authToken, ok := GetAuthToken(ctx); ok {
		if bytes, err := json.Marshal(authToken); err == nil {
			data["authToken"] = string(bytes)
		}
	}
	if app, ok := GetAppInfo(ctx); ok {
		data["appId"] = app.GetAppId()
		data["appName"] = app.GetAppName()
	}
	return data
}

func NewContextWithMap(parent context.Context, data map[string]any) (ctx context.Context, err error) {
	ctx = parent
	if token, ok := data["authToken"]; ok {
		var authToken AuthTokenEntity
		bytes := []byte(token.(string))
		if err = json.Unmarshal(bytes, &authToken); err == nil {
			ctx = NewContextWidthAuthToken(ctx, authToken.GetUser().GetTenantId(), &authToken, map[string][]string{})
		}
	}
	if appId, ok := data["appId"]; ok {
		appName := data["appName"]
		ctx = newAppContext(ctx, appId.(string), appName.(string))
	}
	return ctx, err
}

func GetMessage(ctx context.Context) (res []string) {
	tenantId, tenOk := GetTenantId(ctx)
	res = append(res, fmt.Sprintf("tenantId=%v,ok=%v; ", tenantId, tenOk))

	head, headOk := GetHeader(ctx)
	res = append(res, fmt.Sprintf("header=%v,ok=%v; ", head, headOk))

	authToKen, atOK := GetAuthToken(ctx)
	if atOK {
		res = append(res, fmt.Sprintf("userName=%v,ok=%v; ", authToKen.GetUser().GetName(), atOK))
	} else {
		res = append(res, fmt.Sprintf("userName=nil,ok=%v; ", atOK))
	}
	return res
}
