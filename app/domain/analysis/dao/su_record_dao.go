package dao

import (
	"context"
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
)

type SuRecordDao struct {
	idao.Dao[*model.SuRecord]
}

func NewSuRecordDao(dbKey string) *SuRecordDao {
	tableName := "su_record"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.SuRecord{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.SuRecord](newCfg)
	daoVal := &SuRecordDao{Dao: baseDao}
	return daoVal
}

func (dao *SuRecordDao) CountByTaskId(ctx context.Context, taskId string) (int64, error) {
	sql := rsql.NewBuilder().Eq("task_id", taskId).Build()
	suRecordCount, err := dao.CountByRSQL(ctx, sql)
	return suRecordCount, err
}

func (dao *SuRecordDao) CountByTaskIdHighRisk(ctx context.Context, taskId string) (int64, error) {
	sql := rsql.NewBuilder().And(rsql.Eq("task_id", taskId), rsql.Eq("risk", 5)).Build()
	suRecordCount, err := dao.CountByRSQL(ctx, sql)
	return suRecordCount, err
}

func (dao *SuRecordDao) SumByTaskId(ctx context.Context, taskId string) (float64, error) {
	sql := rsql.NewBuilder().Eq("task_id", taskId).Build()
	fields := []*store.ValueCol{
		{AggFunc: "sum", Field: "amount"},
	}
	sumFields, err := dao.SumByRSQL(ctx, sql, fields)
	if err != nil {
		return 0, err
	}
	if val, ok := sumFields["amount"]; ok {
		amount := val.(float64)
		return amount, nil
	}
	return 0, err
}

func (dao *SuTaskDao) SumAmount(ctx context.Context, taskId string) (map[string]any, error) {
	data, err := dao.SumByRSQL(ctx, fmt.Sprintf("task_id==%s", taskId), []*store.ValueCol{
		{AggFunc: "sum", Field: "amount"},
	})
	return data, err
}
