package appctx

import "context"

type tenantKey struct {
}

var tenantCtxKey = tenantKey{}

func NewTenantContext(parent context.Context, tenantId string) context.Context {
	return context.WithValue(parent, tenantCtxKey, tenantId)
}

// GetTenantId
//
//	@Description: 根据上下文取得租户ID
//	@param ctx
//	@return string
//	@return bool
func GetTenantId(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	val := ctx.Value(tenantCtxKey)
	if val == nil {
		return "", false
	}
	tenantId, ok := val.(string)
	return tenantId, ok
}
