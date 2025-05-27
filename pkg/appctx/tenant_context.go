package appctx

import "context"

type tenantValue struct {
	tenantId string
}
type tenantKey struct {
}

var tenantCtxKey = tenantKey{}

func NewTenantContext(parent context.Context, tenantId string) context.Context {
	tenVal := &tenantValue{tenantId: tenantId}
	return context.WithValue(parent, tenantCtxKey, tenVal)
}

func SetTenantContext(parent context.Context, tenantId string) context.Context {
	var tenVal *tenantValue
	val := parent.Value(tenantCtxKey)
	if val != nil {
		tenVal = val.(*tenantValue)
		tenVal.tenantId = tenantId
		return parent
	}

	tenVal = &tenantValue{tenantId: tenantId}
	return context.WithValue(parent, tenantCtxKey, tenVal)
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
		if user, ok := GetAuthUser(ctx); ok {
			return user.GetTenantId(), true
		}
		return "", false
	}
	tenVal, ok := val.(*tenantValue)
	return tenVal.tenantId, ok
}

// GetTenantId2
//
//	@Description: 根据上下文取得租户ID
//	@param ctx
//	@return string
//	@return bool
func GetTenantId2(ctx context.Context) string {
	if ctx == nil {
		panic("nil context")
	}
	val := ctx.Value(tenantCtxKey)
	if val == nil {
		if user, ok := GetAuthUser(ctx); ok {
			return user.GetTenantId()
		}
		panic("token error in context")
	}
	tenVal, ok := val.(*tenantValue)
	if !ok {
		panic("wrong type")
	}
	return tenVal.tenantId
}
