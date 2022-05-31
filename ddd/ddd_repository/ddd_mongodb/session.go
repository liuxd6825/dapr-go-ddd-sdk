package ddd_mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoSession struct {
	mongodb *MongoDB
}

func NewSession(db *MongoDB) ddd_repository.Session {
	return &MongoSession{mongodb: db}
}

func (r *MongoSession) UseTransaction(ctx context.Context, dbFunc ddd_repository.SessionFunc) error {
	return r.mongodb.client.UseSession(ctx, func(sCtx mongo.SessionContext) error {
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
		return err
	})
}
