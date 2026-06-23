package dao

import (
	"context"

	pkgDao "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

// CreateTables 注册 Product 及其子实体的数据库表
func CreateTables(dbKey string) {
	ctx := context.Background()
	tables := []idao.Table{
		NewProductDao(dbKey).Table(),
		NewProductCompanyDao(dbKey).Table(),
		NewProductContractDao(dbKey).Table(),
		NewProductHumanDao(dbKey).Table(),
		NewProductProductDao(dbKey).Table(),
		NewProductRecordDao(dbKey).Table(),
	}
	pkgDao.CreateTables(ctx, tables)
}