package session

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/dao/neo4j_dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
)

// StartNeo4jSession
//
//	@Description:
//	@param ctx
//	@param fun
//	@param options
//	@return error
func StartNeo4jSession(ctx context.Context, fun Func, options ...*ddd_repository.SessionOptions) error {
	opt := ddd_repository.NewSessionOptions(options...)
	session := neo4j_dao.NewSession(true)
	err := session.UseTransaction(ctx, func(ctx context.Context) error {
		return fun(ctx)
	}, opt)
	if err != nil {
		logs.Errorf(ctx, "func=StartNeo4jSession(); error:%v;", err.Error())
	}
	return err
}

// StartGraphSession
//
//	@Description:
//	@param ctx
//	@param fun
//	@param options
//	@return error
func StartGraphSession(ctx context.Context, fun Func, options ...*ddd_repository.SessionOptions) error {
	return StartNeo4jSession(ctx, fun, options...)
}
