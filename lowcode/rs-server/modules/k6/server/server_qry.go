package server

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
)

type QueryFunc = func(tx context.Context) *common.Result[any]

func (e *Server) DoQuery(ctx *WebContext, fun QueryFunc, opts ...restapp.DoOptions) *common.Result[any] {
	data, _, err := restapp.DoQuery(ctx.ictx, ctx.GetTenantId(), func(ctx context.Context) (interface{}, bool, error) {
		return fun(ctx).GetFoundResults()
	}, opts...)
	return common.NewResult[any](data, err)
}

func (e *Server) DoQueryOne(ctx *WebContext, fun QueryFunc, opts ...restapp.DoOptions) *common.Result[any] {
	data, _, err := restapp.DoQueryOne(ctx.ictx, ctx.GetTenantId(), func(ctx context.Context) (interface{}, bool, error) {
		return fun(ctx).GetFoundResults()
	}, opts...)
	return common.NewResult[any](data, err)
}

func (e *Server) DoRequest(ctx *WebContext, fun QueryFunc, opts ...restapp.DoOptions) *common.Result[any] {
	err := restapp.Do(ctx.ictx, ctx.GetTenantId(), func(ctx context.Context) error {
		return fun(ctx).Error
	}, opts...)
	return common.NewResult[any](nil, err)
}
