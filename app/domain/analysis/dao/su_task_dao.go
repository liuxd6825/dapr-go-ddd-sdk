package dao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type SuTaskDao struct {
	idao.Dao[*model.SuTask]
}

type SuTaskUpdateInfo struct {
	RecordCount *int64              `json:"recordCount"  title:"流水总数"`
	TotalAmount *float64            `json:"totalAmount"  title:"可疑总金额"`
	SuCount     *int64              `json:"suCount"  title:"可疑交易数"`
	SuHighCount *int64              `json:"suHighCount"  title:"高可疑交易数"`
	Status      *model.SuTaskStatus `json:"status" title:"状态"`
}

func NewSuTaskDao(dbKey string) *SuTaskDao {
	tableName := "su_task"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.SuTask{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.SuTask](newCfg)
	daoVal := &SuTaskDao{Dao: baseDao}
	return daoVal
}

func (dao *SuTaskDao) UpdateStatus(ctx context.Context, id string, status model.SuTaskStatus) error {
	data := map[string]any{
		"status": status,
	}
	return dao.UpdateMap(ctx, id, data).GetError()
}

func (dao *SuTaskDao) UpdateWorkflowIdStatus(ctx context.Context, id string, workflowId string, status model.SuTaskStatus) error {
	data := map[string]any{
		"workflow_id": workflowId,
		"status":      status,
	}
	return dao.UpdateMap(ctx, id, data).GetError()
}

func (dao *SuTaskDao) UpdateInfo(ctx context.Context, id string, info SuTaskUpdateInfo) error {
	data := map[string]any{}
	if info.RecordCount != nil {
		data["record_count"] = *info.RecordCount
	}
	if info.TotalAmount != nil {
		data["total_amount"] = *info.TotalAmount
	}
	if info.SuCount != nil {
		data["su_count"] = *info.SuCount
	}
	if info.SuHighCount != nil {
		data["su_high_count"] = *info.SuHighCount
	}
	if info.Status != nil {
		data["status"] = *info.Status
	}
	// SuTask 可疑分析任务
	return dao.UpdateMap(ctx, id, data).GetError()
}
