package server

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
)

type CmdFunc func(ctx context.Context, command *Command) error

func (e *Server) DoCmd(rctx *WebContext, fun CmdFunc, opts ...restapp.DoOptions) {
	err := restapp.DoCmd(rctx.ictx, rctx.GetTenantId(), func(ctx context.Context) error {
		return nil
	}, opts...)
	if err != nil {
		panic(err)
	}
}

func (e *Server) DoCmdAndQueryOne(rctx *WebContext, queryAppId string, cmdFun CmdFunc, queryFun QueryFunc, opts ...restapp.CmdAndQueryOption) any {
	data, _, err := restapp.DoCmdAndQueryOne(rctx.ictx, rctx.GetTenantId(), queryAppId, nil, func(ctx context.Context) error {
		return nil
	}, func(ctx context.Context) (any, bool, error) {
		return queryFun(ctx).GetFoundResults()
	}, opts...)
	if err != nil {
		panic(err)
	}
	return data
}

func (e *Server) DoCmdAndQueryList(ctx *WebContext, queryAppId string, cmd *Command, cmdFun CmdFunc, queryFun QueryFunc, opts ...restapp.CmdAndQueryOption) any {
	data, _, err := restapp.DoCmdAndQueryList(ctx.ictx, ctx.GetTenantId(), queryAppId, cmd, func(ctx context.Context) error {
		err := cmdFun(ctx, cmd)
		return err
	}, func(ctx context.Context) (any, bool, error) {
		return queryFun(ctx).GetFoundResults()
	}, opts...)

	if err != nil {
		panic(err)
	}
	return data
}
