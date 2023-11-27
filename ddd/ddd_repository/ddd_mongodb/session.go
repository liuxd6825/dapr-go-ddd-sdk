package ddd_mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoSession struct {
	mongodb *MongoDB
	isWrite bool
}

func NewSession(isWrite bool, db *MongoDB) ddd_repository.Session {
	return &MongoSession{mongodb: db, isWrite: isWrite}
}

func (r *MongoSession) UseTransaction(ctx context.Context, dbFunc ddd_repository.SessionFunc) error {
	commitTime := 20 * time.Second
	opt := &options.SessionOptions{
		DefaultMaxCommitTime: &commitTime,
	}
	err := r.mongodb.client.UseSessionWithOptions(ctx, opt, func(sCtx mongo.SessionContext) error {
		serverCount := r.mongodb.config.ServerCount()
		if serverCount == 1 {
			return dbFunc(sCtx)
		} else {
			if err := sCtx.StartTransaction(); err != nil {
				return err
			}
			err := dbFunc(sCtx)
			if err != nil {
				err = sCtx.AbortTransaction(ctx)
			} else {
				err = sCtx.CommitTransaction(ctx)
			}
			return err
		}
	})
	if err != nil {
		logs.Error(ctx, err.Error())
	}
	return err
}
