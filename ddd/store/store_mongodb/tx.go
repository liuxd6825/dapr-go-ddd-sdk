package store_mongodb

import (
	ctx "context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"go.mongodb.org/mongo-driver/mongo"
	mongo_options "go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

type txMongoContextKey string

func NewContext(parentCtx ctx.Context, sessionCtx mongo.SessionContext, dbKey string) ctx.Context {
	return context.WithValue(parentCtx, txMongoContextKey(dbKey), sessionCtx)
}

func getSessionContext(ctx ctx.Context, dbKey string) mongo.SessionContext {
	db, ok := ctx.Value(txMongoContextKey(dbKey)).(mongo.SessionContext)
	if !ok {
		return nil
	}
	return db
}

func StartTx(ctx context.Context, mongodb *MongoDB, dbKey string, txFun store.TxFunc, opts ...*store.SessionOptions) error {
	sOpts := &mongo_options.SessionOptions{}
	client := mongodb.client
	serverCount := mongodb.config.ServerCount()

	sessionCtx := getSessionContext(ctx, dbKey)
	if sessionCtx != nil {
		return txFun(sessionCtx, opts...)
	}
	// 事务处理
	err := client.UseSessionWithOptions(ctx, sOpts, func(txCtx mongo.SessionContext) error {
		var tranErr error
		newCtx := NewContext(ctx, txCtx, dbKey)
		if serverCount == 1 {
			return txFun(sessionCtx, opts...)
		}

		err := txCtx.StartTransaction()
		// 开启事务
		if err != nil {
			return err
		}

		// 执行业务
		if err = txFun(newCtx, opts...); err != nil {
			tranErr = txCtx.AbortTransaction(ctx)
		} else {
			tranErr = txCtx.CommitTransaction(ctx)
		}
		if err != nil {
			return err
		}
		return tranErr
	})
	return err
}
