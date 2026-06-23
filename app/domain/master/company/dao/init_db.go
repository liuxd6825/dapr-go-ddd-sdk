package dao

import (
	"context"

	pkgDao "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

// CreateTables 注册 Company 及其子实体的数据库表
func CreateTables(dbKey string) {
	ctx := context.Background()
	tables := []idao.Table{
		NewCompanyDao(dbKey).Table(),
		NewCompanyAccountDao(dbKey).Table(),
		NewCompanyCompanyDao(dbKey).Table(),
		NewCompanyContractDao(dbKey).Table(),
		NewCompanyHumanDao(dbKey).Table(),
		NewCompanyProductDao(dbKey).Table(),
		NewCompanyRecordDao(dbKey).Table(),
	}
	pkgDao.CreateTables(ctx, tables)
}