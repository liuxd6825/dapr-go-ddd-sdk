package field

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
)

type MapHead model.MapHead

// TemplateField
// @titleription:
type TemplateField model.TemplateField

type TaskCreateFields struct {
	Id         string                 `json:"id" validate:"required"`
	CaseId     string                 `json:"caseId" validate:"required"`
	Name       string                 `json:"name" validate:"required"`
	SheetName  string                 `json:"sheetName" validate:"required"`
	SheetId    string                 `json:"sheetId" validate:"required"`
	DocId      string                 `json:"docId" validate:"required"`
	FileId     string                 `json:"fileId" validate:"required"`
	FileName   string                 `json:"fileName" validate:"required"`
	SchemaName string                 `json:"schemaName" validate:"required"`
	SchemaId   string                 `json:"schemaId" validate:"required"`
	MapHeads   []*model.MapHead       `json:"mapHeads" validate:"required"`
	Fields     []*model.TemplateField `json:"fields" validate:"required"`
	Remark     string                 `json:"remark" title:"备注"`
}

type TaskDeleteFields struct {
	Id string `json:"id" title:"Id"  validate:"required"`
}

type TaskRecordCreateFields struct {
	Id        string `json:"id"  validate:"required" title:"Id"`
	BatchSize int64  `json:"batchSize"  validate:"required"  title:"批数据"`
}

type TaskRecordDeleteFields struct {
	Id string `json:"id" validate:"required"`
}

type TaskRecordImportFields struct {
	Id       string `json:"id"  validate:"required" title:"Id" `
	TenantId string `json:"tenantId" validate:"required" title:"租户Id"`
	PageSize int64  `json:"pageSize" validate:"required" title:"分页大小"`
}

type TaskRecordNoImportFields struct {
	Id       string `json:"id" validate:"required" title:"Id"`
	SchemaId string `json:"schemaId" validate:"required"  title:"主数据类型"`
}

type TaskUpdateStateFields struct {
	Id      string          `json:"id" validate:"required" `
	State   model.TaskState `json:"state" validate:"required" `
	Message string          `json:"message" validate:"required" `
}

type TaskUpdateFields struct {
	Id         string                 `json:"id" validate:"required"   title:"Id"`
	CaseId     string                 `json:"caseId" validate:"required"  title:"案件Id" `
	DocId      string                 `json:"docId" validate:"required"  title:"文档Id"`
	FileId     string                 `json:"fileId" validate:"required"  title:"文件Id"`
	FileName   string                 `json:"fileName"  validate:"required" title:"文件名称"`
	SheetId    string                 `json:"sheetId" validate:"required"`
	SheetName  string                 `json:"sheetName"  validate:"required"  title:"Sheet页"`
	SchemaId   string                 `json:"schemaId"  validate:"required" title:"主数据类型"`
	SchemaName string                 `json:"schemaName" validate:"required"`
	MapHeads   []*model.MapHead       `json:"mapHeads"  validate:"required"  title:"表头"`
	Fields     []*model.TemplateField `json:"fields"  validate:"required"  title:"字段"`
	Remark     string                 `json:"remarks" title:"备注"` // 备注
}

type TaskLockField struct {
	Id string `json:"id"  title:"id"  validate:"required" `
}

type TaskStopFields struct {
	Id string `json:"id" title:"Id"  validate:"required" `
}

type TaskUpdateProgressFields struct {
	Id        string          `json:"id" title:"ID"`
	Complete  int64           `json:"complete" title:"完成数量"`
	State     model.TaskState `json:"state"  title:"状态"`
	Message   string          `json:"message"  title:"信息"`
	StartTime *times.Time     `json:"startTime"   title:"开始时间"`
	EndTime   *times.Time     `json:"endTime"  title:"最后日期"`
}

type TaskUpdateRecordTemplateFields struct {
	Id       string               `json:"id"`
	Template model.RecordTemplate `json:"template"`
}

type TaskValidateFields struct {
	Id string `json:"id" title:"Id"`
}
