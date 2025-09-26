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

type SuTaskFields struct {
	RecordCount *int64              `json:"recordCount"  title:"流水总数"`
	TotalAmount *float64            `json:"totalAmount"  title:"可疑总金额"`
	SuCount     *int64              `json:"suCount"  title:"可疑交易数"`
	SuHighCount *int64              `json:"suHighCount"  title:"高可疑交易数"`
	Status      *model.SuTaskStatus `json:"status" title:"状态"`
	WorkflowId  *string             `json:"workflowId" title:"工作流Id"`
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
	values := NewSuTaskFields().SetWorkflowId(workflowId).SetStatus(status)
	return dao.UpdateFields(ctx, id, *values)
}

func (dao *SuTaskDao) UpdateFields(ctx context.Context, id string, fields SuTaskFields) error {
	data := map[string]any{}
	if fields.RecordCount != nil {
		data["record_count"] = *fields.RecordCount
	}
	if fields.TotalAmount != nil {
		data["total_amount"] = *fields.TotalAmount
	}
	if fields.SuCount != nil {
		data["su_count"] = *fields.SuCount
	}
	if fields.SuHighCount != nil {
		data["su_high_count"] = *fields.SuHighCount
	}
	if fields.Status != nil {
		data["status"] = *fields.Status
		data["status_name"] = fields.Status.String()
	}
	if fields.WorkflowId != nil {
		data["workflow_id"] = *fields.WorkflowId
	}
	// SuTask 可疑分析任务
	return dao.UpdateMap(ctx, id, data).GetError()
}

func NewSuTaskFields() *SuTaskFields {
	return &SuTaskFields{}
}

func (f *SuTaskFields) SetRecordCount(recordCount int64) *SuTaskFields {
	f.RecordCount = &recordCount
	return f
}

func (f *SuTaskFields) SetTotalAmount(totalAmount float64) *SuTaskFields {
	f.TotalAmount = &totalAmount
	return f
}

func (f *SuTaskFields) SetSuCount(suCount int64) *SuTaskFields {
	f.SuCount = &suCount
	return f
}

func (f *SuTaskFields) SetSuHighCount(suHighCount int64) *SuTaskFields {
	f.SuHighCount = &suHighCount
	return f
}

func (f *SuTaskFields) SetStatus(status model.SuTaskStatus) *SuTaskFields {
	f.Status = &status
	return f
}

func (f *SuTaskFields) SetWorkflowId(workflowId string) *SuTaskFields {
	f.WorkflowId = &workflowId
	return f
}
