package dao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

func CreateTables(dbKey string) {
	ctx := context.Background()
	tables := []idao.Table{
		NewSuBatchDao(dbKey).Table(),
		NewSuBatchItemDao(dbKey).Table(),
		NewSuRecordDao(dbKey).Table(),
		NewSuAccountDao(dbKey).Table(),
		NewSuTaskDao(dbKey).Table(),
		NewSuTaskLogDao(dbKey).Table(),
	}
	dao.CreateTables(ctx, tables)
}
