package task

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/recordie/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/template"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
)

type MapHead template.MapHead

// TemplateField
// @Description:
type TemplateField template.TemplateField

type TaskCreateFields struct {
	Id        string                    `json:"id"`
	CaseId    string                    `json:"caseId"`
	TenantId  string                    `json:"tenantId"`
	Name      string                    `json:"name"`
	SheetName string                    `json:"sheetName"`
	DocId     string                    `json:"docId"`
	FileId    string                    `json:"fileId"`
	FileName  string                    `json:"fileName"`
	SchemaId  string                    `json:"schemaId" `
	MapHeads  []*template.MapHead       `json:"mapHeads"`
	Fields    []*template.TemplateField `json:"fields"`
	Remark    string                    `json:"remark" desc:"备注"` // 备注
}

type TaskDeleteFields struct {
	Id       string `json:"id" desc:"Id"`
	TenantId string `json:"tenantId" desc:"租户Id"`
}

type TaskRecordCreateFields struct {
	Id        string `json:"id" desc:"Id"`
	BatchSize int64  `json:"batchSize"  desc:"批数据"`
}

type TaskRecordDeleteFields struct {
	TenantId string `json:"tenantId"`
	Id       string `json:"id"`
}

type TaskRecordImportFields struct {
	Id       string `json:"id" desc:"Id"`
	TenantId string `json:"tenantId" desc:"租户Id"`
	PageSize int64  `json:"pageSize" desc:"分页大小"`
}

type TaskRecordNoImportFields struct {
	Id       string `json:"id" desc:"Id"`
	TenantId string `json:"tenantId" desc:"租户Id"`
	SchemaId string `json:"schemaId" desc:"主数据类型"`
}

type TaskUpdateStateFields struct {
	Id      string    `json:"id"`
	State   TaskState `json:"state"`
	Message string    `json:"message"`
}

type TaskUpdateFields struct {
	Id        string                    `json:"id" desc:"Id"`
	TenantId  string                    `json:"tenantId" desc:"租户Id"`
	CaseId    string                    `json:"caseId"  desc:"案件Id" `
	DocId     string                    `json:"docId" desc:"文档Id"`
	FileId    string                    `json:"fileId" desc:"文件Id"`
	FileName  string                    `json:"fileName" desc:"文件名称"`
	SheetName string                    `json:"sheetName"  desc:"Sheet页"`
	StartTime *times.Time               `json:"startTime"  desc:"开始时间"`
	SchemaId  string                    `json:"schemaId" desc:"主数据类型"`
	MapHeads  []*template.MapHead       `json:"mapHeads"  desc:"表头"`
	Fields    []*template.TemplateField `json:"fields"  desc:"字段"`
	Remarks   string                    `json:"remarks" desc:"备注"` // 备注
}

type TaskLockField struct {
	Id string `json:"id"  desc:"Id"`
}

type TaskStopFields struct {
	Id       string `json:"id" desc:"Id"`
	TenantId string `json:"tenantId" desc:"租户Id"`
}

type TaskUpdateProgressFields struct {
	Id           string      `json:"id" title:"开始时间"`
	CompleteRows int64       `json:"completeRows" title:"开始时间"`
	State        TaskState   `json:"state"  title:"状态"`
	Message      string      `json:"message"  title:"开始时间"`
	StartTime    *times.Time `json:"startTime"   title:"开始时间"`
	EndTime      *times.Time `json:"endTime"  title:"最后日期"`
}

type TaskUpdateRecordTemplateFields struct {
	Id       string               `json:"id"`
	TenantId string               `json:"tenantId"`
	Template model.RecordTemplate `json:"template"`
}

type TaskValidateFields struct {
	Id       string `json:"id" desc:"Id"`
	TenantId string `json:"tenantId" desc:"租户Id"`
}
