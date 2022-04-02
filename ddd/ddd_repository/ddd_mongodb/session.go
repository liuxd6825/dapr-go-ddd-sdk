package ddd_mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoSession struct {
	mongodb *MongoDB
}

func UseSession(ctx context.Context, db *MongoDB, fun ddd_repository.SessionFunc) error {
	sess := &MongoSession{mongodb: db}
	return sess.UseSession(ctx, fun)
}

func (r *MongoSession) UseSession(ctx context.Context, sessionFn func(ctx context.Context) error) error {
	return r.mongodb.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		if err := sessionContext.StartTransaction(); err != nil {
			return err
		}
		err := sessionFn(ctx)
		if err != nil {
			if e1 := sessionContext.AbortTransaction(ctx); e1 != nil {
				err = e1
			}
		} else {
			err = sessionContext.CommitTransaction(ctx)
		}
		return err
	})
}
