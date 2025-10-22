package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type UserDao struct {
	idao.Dao[*model.User]
}

func NewUserDao(dbKey string) *UserDao {
	tableName := "portal_user"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.User{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.User](newCfg)
	daoVal := &UserDao{Dao: baseDao}
	return daoVal
}
