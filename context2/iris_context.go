package context2

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_context"
)

// GetIrisContext
//
//	@Description: 获取Iris上下文
//	@param ctx
//	@return iris.Context
func GetIrisContext(ctx context.Context) iris.Context {
	v := ddd_context.GetServerContext(ctx)
	if s, ok := v.(*serverContext); ok {
		return s.ctx
	}
	return nil
}
