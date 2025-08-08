package dao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

func CreateTables(dbKey string) {
	ctx := context.Background()
	tables := []idao.Table{
		NewExcelFileDao(dbKey).Table(),
		NewExcelRowDao(dbKey).Table(),
		NewExcelSheetDao(dbKey).Table(),
		NewRecordDao(dbKey).Table(),
		NewSchemaDao(dbKey).Table(),
		NewTaskDao(dbKey).Table(),
		NewTemplateDao(dbKey).Table(),
	}
	dao.CreateTables(ctx, tables)
}
