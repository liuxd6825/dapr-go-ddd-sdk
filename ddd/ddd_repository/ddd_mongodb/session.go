package ddd_mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type MongoSession struct {
	mongodb *MongoDB
}

type sessionCtxKey struct {
}

func NewSession(isWrite bool, db *MongoDB) ddd_repository.Session {
	return &MongoSession{mongodb: db}
}

func (r *MongoSession) UseTransaction(ctx context.Context, dbFunc ddd_repository.SessionFunc, opts ...*ddd_repository.SessionOptions) (err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
		if err != nil {
			logs.Error(ctx, "", logs.Fields{"func": "MongoSession.UseTransaction()", "error": err.Error()})
		}
	}()

	writeTime := ddd_repository.NewSessionOptions(opts...).GetWriteTime()

	// 事务选项
	opt := &options.SessionOptions{
		DefaultMaxCommitTime: &writeTime,
		DefaultWriteConcern: &writeconcern.WriteConcern{
			WTimeout: writeTime,
		},
	}

	// 事务处理
	err = r.mongodb.client.UseSessionWithOptions(ctx, opt, func(ctx mongo.SessionContext) error {
		serverCount := r.mongodb.config.ServerCount()
		if serverCount == 1 {
			return dbFunc(ctx)
		} else {
			// 开启事务
			if err = ctx.StartTransaction(); err != nil {
				return err
			}
			// 执行业务
			if err = dbFunc(ctx); err != nil {
				err = ctx.AbortTransaction(ctx)
			} else {
				err = ctx.CommitTransaction(ctx)
			}
			return err
		}
	})
	return err
}
