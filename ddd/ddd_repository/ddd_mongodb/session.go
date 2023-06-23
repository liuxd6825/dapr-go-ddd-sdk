package ddd_mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoSession struct {
	mongodb *MongoDB
	isWrite bool
}

func NewSession(isWrite bool, db *MongoDB) ddd_repository.Session {
	return &MongoSession{mongodb: db, isWrite: isWrite}
}

func (r *MongoSession) UseTransaction(ctx context.Context, dbFunc ddd_repository.SessionFunc) error {
	err := r.mongodb.client.UseSession(ctx, func(sCtx mongo.SessionContext) error {
		serverCount := r.mongodb.config.ServerCount()
		if serverCount == 1 {
			return dbFunc(sCtx)
		} else {
			if err := sCtx.StartTransaction(); err != nil {
				return err
			}
			err := dbFunc(sCtx)
			if err != nil {
				if e1 := sCtx.AbortTransaction(sCtx); e1 != nil {
					err = e1
				}
			} else {
				err = sCtx.CommitTransaction(sCtx)
			}
			if err != nil {
				println(err)
			}
			return err
		}
	})
	if err != nil {
		logs.Error(ctx, err)
	}
	return err
}
