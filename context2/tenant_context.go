package context2

import "context"

type tenantCtxKey struct {
}

func NewTenantContext(ctx context.Context, tenantId string) context.Context {
	return context.WithValue(ctx, tenantCtxKey{}, tenantId)
}

// GetTenantId
//
//	@Description: 根据上下文取得租户ID
//	@param ctx
//	@return string
//	@return bool
func GetTenantId(ctx context.Context) (string, bool) {
	val := ctx.Value(tenantCtxKey{})
	if val == nil {
		return "", false
	}
	tenantId, ok := val.(string)
	return tenantId, ok
}
