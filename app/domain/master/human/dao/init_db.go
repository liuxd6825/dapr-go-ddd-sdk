package dao

import (
	"context"

	pkgDao "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

// CreateTables 注册 Human 及其子实体的数据库表
func CreateTables(dbKey string) {
	ctx := context.Background()
	tables := []idao.Table{
		NewHumanDao(dbKey).Table(),
		NewHumanAccountDao(dbKey).Table(),
		NewHumanAddressDao(dbKey).Table(),
		NewHumanCapitalDao(dbKey).Table(),
		NewHumanCompanyDao(dbKey).Table(),
		NewHumanContractDao(dbKey).Table(),
		NewHumanCredentialDao(dbKey).Table(),
		NewHumanExtDao(dbKey).Table(),
		NewHumanHumanDao(dbKey).Table(),
		NewHumanLinkDao(dbKey).Table(),
		NewHumanProductDao(dbKey).Table(),
		NewHumanRecordDao(dbKey).Table(),
		NewHumanReportedAmountDao(dbKey).Table(),
		NewHumanSuspectAmountDao(dbKey).Table(),
	}
	pkgDao.CreateTables(ctx, tables)
}