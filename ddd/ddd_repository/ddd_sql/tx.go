package ddd_sql

import (
	ctx "context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type txSqlContextKey string

func NewContext(parentCtx ctx.Context, tx *gorm.DB, dbKey string) ctx.Context {
	return context.WithValue(parentCtx, txSqlContextKey(dbKey), tx)
}

func GetTx(ctx ctx.Context, dbKey string) *gorm.DB {
	db, ok := ctx.Value(txSqlContextKey(dbKey)).(*gorm.DB)
	if !ok {
		return nil
	}
	return db
}

func StartTx(ctx context.Context, db *gorm.DB, dbKey string, fun ddd_repository.TxFunc, opts ...*ddd_repository.SessionOptions) (err error) {
	tx := GetTx(ctx, dbKey)
	if tx != nil {
		return fun(ctx, opts...)
	}

	err = db.Transaction(func(txDb *gorm.DB) error {
		txCtx := NewContext(ctx, txDb, dbKey)
		return fun(txCtx, opts...)
	})
	return err
}
