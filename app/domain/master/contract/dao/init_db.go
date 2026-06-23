package dao

import (
	"context"

	pkgDao "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

// CreateTables 注册 Contract 及其子实体的数据库表
func CreateTables(dbKey string) {
	ctx := context.Background()
	tables := []idao.Table{
		NewContractDao(dbKey).Table(),
		NewContractCompanyDao(dbKey).Table(),
		NewContractContractDao(dbKey).Table(),
		NewContractHumanDao(dbKey).Table(),
		NewContractProductDao(dbKey).Table(),
		NewContractRecordDao(dbKey).Table(),
	}
	pkgDao.CreateTables(ctx, tables)
}