package ddd_mongodb

import (
	ctx "context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type txMongoContextKey string

func NewContext(parentCtx ctx.Context, sessionCtx mongo.SessionContext, dbName string) ctx.Context {
	return context.WithValue(parentCtx, txMongoContextKey(dbName), sessionCtx)
}

func getSessionContext(ctx ctx.Context, dbName string) mongo.SessionContext {
	db, ok := ctx.Value(txMongoContextKey(dbName)).(mongo.SessionContext)
	if !ok {
		return nil
	}
	return db
}

func StartTx(ctx context.Context, txFunc ddd_repository.TxFunc, options ...*ddd_repository.SessionOptions) error {
	opt := ddd_repository.NewSessionOptions(options...)

	return nil
}
