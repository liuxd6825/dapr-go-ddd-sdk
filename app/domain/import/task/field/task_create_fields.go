package field

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/template/model"

type TaskCreateFields struct {
	Id        string                 `json:"id"`
	CaseId    string                 `json:"caseId"`
	TenantId  string                 `json:"tenantId"`
	Name      string                 `json:"name"`
	SheetName string                 `json:"sheetName"`
	DocId     string                 `json:"docId"`
	FileId    string                 `json:"fileId"`
	FileName  string                 `json:"fileName"`
	SchemaId  string                 `json:"schemaId" `
	MapHeads  []*model.MapHead       `json:"mapHeads"`
	Fields    []*model.TemplateField `json:"fields"`
	Remark    string                 `json:"remark" desc:"备注"` // 备注
}
