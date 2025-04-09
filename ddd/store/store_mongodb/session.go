package store_mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type MongoSession struct {
	mongodb IMongoDB
}

type sessionCtxKey struct {
}

func NewSession(isWrite bool, db IMongoDB) store.Session {
	return &MongoSession{mongodb: db}
}

func (r *MongoSession) UseTransaction(ctx context.Context, dbFunc store.SessionFunc, opts ...*store.SessionOptions) (err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
		if err != nil {
			logs.Error(ctx, "", logs.Fields{"func": "MongoSession.UseTransaction()", "error": err.Error()})
		}
	}()

	writeTime := store.NewSessionOptions(opts...).GetWriteTime()

	// 事务选项
	opt := &options.SessionOptions{
		DefaultMaxCommitTime: &writeTime,
		DefaultWriteConcern: &writeconcern.WriteConcern{
			WTimeout: writeTime,
		},
	}

	// 事务处理
	err = r.mongodb.GetClient().UseSessionWithOptions(ctx, opt, func(ctx mongo.SessionContext) error {
		serverCount := r.mongodb.GetServerCount()
		if serverCount == 1 {
			return dbFunc(ctx)
		} else {
			// 开启事务
			if err = ctx.StartTransaction(); err != nil {
				return err
			}
			var tranErr error
			// 执行业务
			if err = dbFunc(ctx); err != nil {
				tranErr = ctx.AbortTransaction(ctx)
			} else {
				tranErr = ctx.CommitTransaction(ctx)
			}
			if err != nil {
				return err
			}
			return tranErr
		}
	})
	return err
}
