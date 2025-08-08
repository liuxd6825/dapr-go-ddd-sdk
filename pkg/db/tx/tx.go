package tx

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/tx"
)

type TxCfg struct {
	DBKeys []string
}

func (c TxCfg) SetDBKeys(dbKeys ...string) TxCfg {
	c.DBKeys = dbKeys
	return c
}

func NewTxCfg(dbKeys ...string) TxCfg {
	return TxCfg{
		DBKeys: dbKeys,
	}
}

func StartTx(ctx context.Context, cfg TxCfg, txFunc store.TxFunc, options ...*store.SessionOptions) (err error) {
	return tx.StartTx(ctx, cfg.DBKeys, txFunc, options...)
}
