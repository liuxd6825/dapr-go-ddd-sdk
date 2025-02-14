package tx

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_sql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"gorm.io/gorm"
)

func StartTx(ctx context.Context, dbNames []string, txFunc ddd_repository.TxFunc, options ...*ddd_repository.SessionOptions) (err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()
	for _, dbKey := range dbNames {
		item := restapp.GetDb(dbKey)
		if item == nil {
			return fmt.Errorf("db %s not found", dbKey)
		}
		dbType := item.GetDBType()
		switch dbType {
		case restapp.DbType_Redis:
			break
		case restapp.DbType_MongoDB:
			break
		case restapp.DbType_Neo4j:
			break
		case restapp.DbType_Sqlite:
			startGorm(ctx, dbKey, item.GetSqlite(), txFunc)
			break
		case restapp.DbType_MySQL:
			startGorm(ctx, dbKey, item.GetMySQL(), txFunc)
			break
		case restapp.DbType_MsSQL:
			startGorm(ctx, dbKey, item.GetMsSQL(), txFunc)
			break
		case restapp.DbType_Oracle:
			startGorm(ctx, dbKey, item.GetOracle(), txFunc)
			break
		case restapp.DbType_Postgres:
			startGorm(ctx, dbKey, item.GetPostgres(), txFunc)
			break
		}
	}
	return nil
}

func startGorm(ctx context.Context, dbKey string, db *gorm.DB, txFunc ddd_repository.TxFunc, options ...*ddd_repository.SessionOptions) {
	err := ddd_sql.StartTx(ctx, db, dbKey, txFunc)
	if err != nil {
		panic(err)
	}
}
