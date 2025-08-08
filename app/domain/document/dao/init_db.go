package dao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

func CreateTables(dbKey string) {
	ctx := context.Background()
	tables := []idao.Table{
		NewDocumentDao(dbKey).Table(),
		NewFileDao(dbKey).Table(),
		NewFolderDao(dbKey).Table(),
	}
	dao.CreateTables(ctx, tables)
}
