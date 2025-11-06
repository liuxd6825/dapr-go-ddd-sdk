package tx

import (
	"context"

	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/tx"
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

func StartTx(ctx context.Context, cfg TxCfg, txFunc store2.TxFunc, options ...*store2.SessionOptions) (err error) {
	return tx.StartTx(ctx, cfg.DBKeys, txFunc, options...)
}
