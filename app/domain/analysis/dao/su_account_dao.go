package dao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
)

type SuAccountDao struct {
	idao.Dao[*model.SuTaskAccount]
}

func NewSuAccountDao(dbKey string) *SuAccountDao {
	tableName := "su_account"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.SuTaskAccount{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.SuTaskAccount](newCfg)
	daoVal := &SuAccountDao{Dao: baseDao}
	return daoVal
}

func (s *SuAccountDao) FindByTaskId(ctx context.Context, taskId string) ([]*model.SuTaskAccount, error) {
	builder := rsql.NewBuilder().And(
		rsql.Eq("su_task_id", taskId),
	)
	return s.FindByRSQL(ctx, builder.Build())
}
