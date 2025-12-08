package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type TenantUserDao struct {
	idao.Dao[*model.TenantUser]
}

func NewTenantUserDao(dbKey string) *TenantUserDao {
	tableName := "portal_tenant_user"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.TenantUser{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.TenantUser](newCfg)
	daoVal := &TenantUserDao{Dao: baseDao}
	return daoVal
}
