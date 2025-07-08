package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type TemplateDao struct {
	idao.Dao[*model.Template]
}

func NewTemplateDao(dbKey string) *TemplateDao {
	tableName := "import_template"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Template{}, tableName)
	newCfg := &dao.NewConfig{
		DBKey:      dbKey,
		IsPubEvent: dao.IsFalse(),
		TableName:  tableName,
		DBSchema:   dbSch,
	}
	baseDao := dao.NewDao[*model.Template](newCfg)
	daoVal := &TemplateDao{Dao: baseDao}
	return daoVal
}
