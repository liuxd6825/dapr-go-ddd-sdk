package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type TenantDao struct {
	idao.Dao[*model.Tenant]
}

func NewTenantDao(dbKey string) *TenantDao {
	tableName := "portal_tenant"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Tenant{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.Tenant](newCfg)
	daoVal := &TenantDao{Dao: baseDao}
	return daoVal
}
