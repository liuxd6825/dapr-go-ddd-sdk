package dao

import (
	"context"
	dao2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

func CreateTables(dbKey string) {
	ctx := context.Background()
	tables := []idao.Table{
		NewRecordDao(dbKey).Table(),
		NewTranDao(dbKey).Table(),
		NewTagRelationDao(dbKey).Table(),

		dao2.NewSuBatchDao(dbKey).Table(),
		dao2.NewSuBatchItemDao(dbKey).Table(),
		dao2.NewSuRecordDao(dbKey).Table(),
		dao2.NewSuAccountDao(dbKey).Table(),
		dao2.NewSuTaskDao(dbKey).Table(),
	}
	dao.CreateTables(ctx, tables)
}
