package field

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/template/model"
	"time"
)

type TaskUpdateFields struct {
	Id        string                 `json:"id" desc:"Id"`
	TenantId  string                 `json:"tenantId" desc:"租户Id"`
	CaseId    string                 `json:"caseId"  desc:"案件Id" `
	DocId     string                 `json:"docId" desc:"文档Id"`
	FileId    string                 `json:"fileId" desc:"文件Id"`
	FileName  string                 `json:"fileName" desc:"文件名称"`
	SheetName string                 `json:"sheetName"  desc:"Sheet页"`
	StartTime *time.Time             `json:"startTime"  desc:"开始时间"`
	SchemaId  string                 `json:"schemaId" desc:"主数据类型"`
	MapHeads  []*model.MapHead       `json:"mapHeads"  desc:"表头"`
	Fields    []*model.TemplateField `json:"fields"  desc:"字段"`
	Remarks   string                 `json:"remarks" desc:"备注"` // 备注
}
