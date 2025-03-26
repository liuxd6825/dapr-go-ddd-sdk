package tx

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/store_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/store_sql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"gorm.io/gorm"
)

// StartTx 开启事务
func StartTx(ctx context.Context, dbKeys []string, txFunc store.TxFunc, options ...*store.SessionOptions) (err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()
	newCtx := ctx
	newTxFunc := txFunc
	if dbKeys != nil && len(dbKeys) > 0 {
		for _, dbKey := range dbKeys {
			item := restapp.GetDB(dbKey)
			if item == nil {
				return fmt.Errorf("db %s not found", dbKey)
			}
			dbType := item.GetDBType()
			switch dbType {
			case restapp.DbType_Redis:
				break
			case restapp.DbType_Neo4j:
				break
			case restapp.DbType_MongoDB:
				newTxFunc = newMongoFunc(item.GetMongo(), dbKey, newTxFunc)
				break
			case restapp.DbType_Sqlite:
				newTxFunc = newGormFunc(item.GetGormDB(), dbKey, newTxFunc)
				break
			case restapp.DbType_MySQL:
				newTxFunc = newGormFunc(item.GetGormDB(), dbKey, newTxFunc)
				break
			case restapp.DbType_MsSQL:
				newTxFunc = newGormFunc(item.GetGormDB(), dbKey, newTxFunc)
				break
			case restapp.DbType_Oracle:
				newTxFunc = newGormFunc(item.GetGormDB(), dbKey, newTxFunc)
				break
			case restapp.DbType_Postgres:
				newTxFunc = newGormFunc(item.GetGormDB(), dbKey, newTxFunc)
				break
			}
		}
	}
	return newTxFunc(newCtx, options...)
}

func newGormFunc(db *gorm.DB, dbKey string, txFunc store.TxFunc, opts ...*store.SessionOptions) store.TxFunc {
	return func(ctx context.Context, opts ...*store.SessionOptions) error {
		return store_sql.StartTx(ctx, db, dbKey, txFunc, opts...)
	}
}

func newMongoFunc(mongodb *store_mongodb.MongoDB, dbKey string, txFunc store.TxFunc, opts ...*store.SessionOptions) store.TxFunc {
	return func(ctx context.Context, opts ...*store.SessionOptions) error {
		return store_mongodb.StartTx(ctx, mongodb, dbKey, txFunc, opts...)
	}
}
