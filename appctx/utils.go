package appctx

import (
	"context"
	"fmt"
	"strings"
)

const Authorization = "Authorization"

func NewContext(parent context.Context, tenantId string, token string, header map[string][]string) (ctx context.Context) {
	ctx = NewTenantContext(parent, tenantId)
	ctx = NewHeaderContext(ctx, header)
	if ctx1, err1 := NewAuthContext(ctx, getAuthorization(token, header)); err1 == nil {
		ctx = ctx1
	}
	return ctx
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
