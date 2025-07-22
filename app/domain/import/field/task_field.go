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
	Id        string                 `json:"id" required:"true"`
	CaseId    string                 `json:"caseId" required:"true"`
	Name      string                 `json:"name" required:"true"`
	SheetName string                 `json:"sheetName" required:"true"`
	DocId     string                 `json:"docId" required:"true"`
	FileId    string                 `json:"fileId" required:"true"`
	FileName  string                 `json:"fileName" required:"true"`
	SchemaId  string                 `json:"schemaId" required:"true"`
	MapHeads  []*model.MapHead       `json:"mapHeads" required:"true"`
	Fields    []*model.TemplateField `json:"fields" required:"true"`
	Remark    string                 `json:"remark" title:"备注"`
}

type TaskDeleteFields struct {
	Id string `json:"id" title:"Id"  required:"true"`
}

type TaskRecordCreateFields struct {
	Id        string `json:"id"  required:"true" title:"Id"`
	BatchSize int64  `json:"batchSize"  required:"true"  title:"批数据"`
}

type TaskRecordDeleteFields struct {
	Id string `json:"id" required:"true"`
}

type TaskRecordImportFields struct {
	Id       string `json:"id"  required:"true" title:"Id" `
	TenantId string `json:"tenantId" required:"true" title:"租户Id"`
	PageSize int64  `json:"pageSize" required:"true" title:"分页大小"`
}

type TaskRecordNoImportFields struct {
	Id       string `json:"id" required:"true" title:"Id"`
	SchemaId string `json:"schemaId" required:"true"  title:"主数据类型"`
}

type TaskUpdateStateFields struct {
	Id      string          `json:"id" required:"true" `
	State   model.TaskState `json:"state" required:"true" `
	Message string          `json:"message" required:"true" `
}

type TaskUpdateFields struct {
	Id        string                 `json:"id" required:"true"   title:"Id"`
	CaseId    string                 `json:"caseId" required:"true"  title:"案件Id" `
	DocId     string                 `json:"docId" required:"true"  title:"文档Id"`
	FileId    string                 `json:"fileId" required:"true"  title:"文件Id"`
	FileName  string                 `json:"fileName"  required:"true" title:"文件名称"`
	SheetName string                 `json:"sheetName"  required:"true"  title:"Sheet页"`
	StartTime *times.Time            `json:"startTime"  required:"true"  title:"开始时间"`
	SchemaId  string                 `json:"schemaId"  required:"true" title:"主数据类型"`
	MapHeads  []*model.MapHead       `json:"mapHeads"  required:"true"  title:"表头"`
	Fields    []*model.TemplateField `json:"fields"  required:"true"  title:"字段"`
	Remarks   string                 `json:"remarks"  required:"true" title:"备注"` // 备注
}

type TaskLockField struct {
	Id string `json:"id"  title:"id"  required:"true" `
}

type TaskStopFields struct {
	Id string `json:"id" title:"Id"  required:"true" `
}

type TaskUpdateProgressFields struct {
	Id           string          `json:"id" title:"ID"`
	CompleteRows int64           `json:"completeRows" title:"完成数量"`
	State        model.TaskState `json:"state"  title:"状态"`
	Message      string          `json:"message"  title:"信息"`
	StartTime    *times.Time     `json:"startTime"   title:"开始时间"`
	EndTime      *times.Time     `json:"endTime"  title:"最后日期"`
}

type TaskUpdateRecordTemplateFields struct {
	Id       string               `json:"id"`
	Template model.RecordTemplate `json:"template"`
}

type TaskValidateFields struct {
	Id string `json:"id" title:"Id"`
}
