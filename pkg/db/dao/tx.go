package dao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/tx"
)

func StartTx(ctx context.Context, dbKeys []string, txFunc store.TxFunc, options ...*store.SessionOptions) (err error) {
	return tx.StartTx(ctx, dbKeys, txFunc, options...)
}
