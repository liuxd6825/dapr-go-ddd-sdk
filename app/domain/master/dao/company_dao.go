package dao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
)

type CompanyDao struct {
	idao.Dao[*model.Company]
}

func NewCompanyDao(dbKey string) *CompanyDao {
	tableName := "master_record"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Company{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.Company](newCfg)
	daoVal := &CompanyDao{Dao: baseDao}
	return daoVal
}

func (d *CompanyDao) FindByCaseId(ctx context.Context, caseId string) ([]*model.Company, error) {
	return d.FindByRSQL(ctx, rsql.NewBuilder().And(rsql.Eq("case", caseId)).Build())
}
