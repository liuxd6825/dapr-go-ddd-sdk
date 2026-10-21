package tx

import (
	"context"
	"fmt"

	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	store_mongodb2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/store_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/store_sql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"gorm.io/gorm"
)

// StartTx 开启事务
func StartTx(ctx context.Context, dbKeys []string, txFunc store2.TxFunc, options ...*store2.SessionOptions) (err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()
	newCtx := ctx
	newTxFunc := txFunc
	if dbKeys != nil && len(dbKeys) > 0 {
		for _, dbKey := range dbKeys {
			item := envs.GetDB(dbKey)
			if item == nil {
				return fmt.Errorf("dbKey %s not found", dbKey)
			}
			DBType := item.GetDBType()
			switch DBType {
			case envs.DBType_Neo4j:
				break
			case envs.DBType_MongoDB:
				newTxFunc = newMongoFunc(item.GetMongo(), dbKey, newTxFunc)
				break
			case envs.DBType_Sqlite, envs.DBType_MySQL, envs.DBType_MsSQL, envs.DBType_Oracle, envs.DBType_Postgres:
				newTxFunc = newGormFunc(item.GetGormDB(), dbKey, newTxFunc)
				break
			}
		}
	}
	return newTxFunc(newCtx, options...)
}

type TxDB []TxDBItem

type TxDBItem struct {
	DB     any
	DBKey  string
	DBType envs.DBType
}

// Start 开启事务
func Start(ctx context.Context, txDb TxDB, txFunc store2.TxFunc, options ...*store2.SessionOptions) (err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()
	newCtx := ctx
	newTxFunc := txFunc
	if txDb != nil && len(txDb) > 0 {
		for _, item := range txDb {
			switch item.DBType {
			case envs.DBType_Neo4j:
				break
			case envs.DBType_MongoDB:
				db, ok := item.DB.(store_mongodb2.IMongoDB)
				if !ok {
					return fmt.Errorf("db %s type %s is not MongoDB", item.DBKey, item.DBType)
				}
				newTxFunc = newMongoFunc(db, item.DBKey, newTxFunc)
			case envs.DBType_Sqlite, envs.DBType_MySQL, envs.DBType_MsSQL, envs.DBType_Oracle, envs.DBType_Postgres:
				db := item.DB.(*gorm.DB)
				newTxFunc = newGormFunc(db, item.DBKey, newTxFunc)
				break
			}
		}
	}
	return newTxFunc(newCtx, options...)
}

func newGormFunc(db *gorm.DB, dbKey string, txFunc store2.TxFunc, opts ...*store2.SessionOptions) store2.TxFunc {
	return func(ctx context.Context, opts ...*store2.SessionOptions) error {
		return store_sql.StartTx(ctx, db, dbKey, txFunc, opts...)
	}
}

func newMongoFunc(mongodb store_mongodb2.IMongoDB, dbKey string, txFunc store2.TxFunc, opts ...*store2.SessionOptions) store2.TxFunc {
	return func(ctx context.Context, opts ...*store2.SessionOptions) error {
		return store_mongodb2.StartTx(ctx, mongodb, dbKey, txFunc, opts...)
	}
}
