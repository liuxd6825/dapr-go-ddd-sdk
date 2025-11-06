package store_mongodb

import (
	ctx "context"

	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	mongo_options "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	writeconcern "go.mongodb.org/mongo-driver/mongo/writeconcern"
	"golang.org/x/net/context"
)

type txMongoContextKey string

func NewContext(parentCtx ctx.Context, sessionCtx mongo.SessionContext, dbKey string) ctx.Context {
	return context.WithValue(parentCtx, txMongoContextKey(dbKey), sessionCtx)
}

func getSessionContext(ctx ctx.Context, dbKey string) mongo.SessionContext {
	if sCtx, ok := ctx.(mongo.SessionContext); ok {
		return sCtx
	}
	if sCtx, ok := ctx.Value(txMongoContextKey(dbKey)).(mongo.SessionContext); ok {
		return sCtx
	}
	return nil
}

func StartTx(ctx context.Context, mongodb IMongoDB, dbKey string, txFun store2.TxFunc, opts ...*store2.SessionOptions) error {
	//serverCount := mongodb.GetServerCount()
	client := mongodb.GetClient()
	// 是否已经在事务中
	if sCtx := getSessionContext(ctx, dbKey); sCtx != nil {
		return txFun(sCtx, opts...)
	}

	// 定义事务选项
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot() // Snapshot 或更高隔离级别
	//sOpts := mongo_options.Session().SetDefaultReadConcern(rc).SetDefaultWriteConcern(wc)
	tOpts := mongo_options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, _ := client.StartSession()
	defer session.EndSession(ctx)

	if err := session.StartTransaction(tOpts); err != nil {
		return err
	}

	err := mongo.WithSession(ctx, session, func(txCtx mongo.SessionContext) (err error) {
		defer func() {
			err = errors.GetRecoverError(err, recover())
			if err != nil {
				_ = txCtx.AbortTransaction(ctx)
			} else {
				_ = txCtx.CommitTransaction(ctx)
			}
		}()

		err = txFun(txCtx, opts...)
		return err
	})

	return err

}
