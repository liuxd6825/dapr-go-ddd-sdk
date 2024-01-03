package session

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/dao/mongo_dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
)

// StartMongoSession
//
//	@Description:
//	@param ctx
//	@param fun
//	@param options
//	@return error
func StartMongoSession(ctx context.Context, fun Func, options ...*ddd_repository.SessionOptions) error {
	opt := ddd_repository.NewSessionOptions(options...)
	session := mongo_dao.NewSession(true)
	err := session.UseTransaction(ctx, func(ctx context.Context) error {
		return fun(ctx)
	}, opt)
	if err != nil {
		logs.Errorf(ctx, "func=StartMongoSession(); error:%v;", err.Error())
	}
	return err
}
