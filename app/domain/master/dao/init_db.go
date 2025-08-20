package dao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

func CreateTables(dbKey string) {
	ctx := context.Background()
	tables := []idao.Table{
		NewRecordDao(dbKey).Table(),
		NewTranDao(dbKey).Table(),
		NewTagRelationDao(dbKey).Table(),

		NewSuBatchDao(dbKey).Table(),
		NewSuBatchItemDao(dbKey).Table(),
		NewSuRecordDao(dbKey).Table(),
		NewSuTaskAccountDao(dbKey).Table(),
		NewSuTaskDao(dbKey).Table(),
	}
	dao.CreateTables(ctx, tables)
}
