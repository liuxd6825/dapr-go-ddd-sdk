package server

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rs-server/modules/common"
)

type CmdFunc func(ctx context.Context, command Command) error

func (e *Server) DoCmd(rctx *RContext, fun CmdFunc, opts ...restapp.DoOptions) *common.Result[any] {
	err := restapp.DoCmd(rctx.ictx, rctx.GetTenantId(), func(ctx context.Context) error {
		return rctx.NewCommand().DoSuccess(func(cmd Command) error {
			return fun(ctx, cmd)
		}).Error
	}, opts...)
	return common.NewResult[any](nil, err)
}

func (e *Server) DoCmdAndQueryOne(rctx *RContext, queryAppId string, cmdFun CmdFunc, queryFun QueryFunc, opts ...restapp.CmdAndQueryOption) *common.Result[any] {
	cmd, err := rctx.NewCommand().GetResults()
	if err != nil {
		return common.NewResult[any](nil, err)
	}

	data, _, err := restapp.DoCmdAndQueryOne(rctx.ictx, rctx.GetTenantId(), queryAppId, cmd, func(ctx context.Context) error {
		return cmdFun(ctx, rctx.NewCommand().Data)
	}, func(ctx context.Context) (any, error) {
		data, err := queryFun(ctx).GetResults()
		return data, err
	}, opts...)
	return common.NewResult[any](data, err)
}

func (e *Server) DoCmdAndQueryList(ctx *RContext, queryAppId string, cmd Command, cmdFun CmdFunc, queryFun QueryFunc, opts ...restapp.CmdAndQueryOption) *common.Result[any] {
	data, _, err := restapp.DoCmdAndQueryList(ctx.ictx, ctx.GetTenantId(), queryAppId, cmd, func(ctx context.Context) error {
		err := cmdFun(ctx, cmd)
		return err
	}, func(ctx context.Context) (any, error) {
		data, err := queryFun(ctx).GetResults()
		return data, err
	}, opts...)
	return common.NewResult[any](data, err)
}
