package cmdquery

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"

type FindRowBySheetQuery struct {
	TenantId string `json:"tenantId"`
	CaseId   string `json:"caseId"`
	DocId    string `json:"docId"`
	FileId   string `json:"fileId"`
	SheetId  string `json:"sheetId"`
}

type FindSheetByNameQuery struct {
	TenantId string `json:"tenantId"`
	CaseId   string `json:"caseId"`
	DocId    string `json:"docId"`
	FileId   string `json:"fileId"`
	Name     string `json:"name"`
}

func (q *FindSheetByNameQuery) Validate() error {
	ve := errors.NewVerifyError()
	if len(q.TenantId) == 0 {
		ve.AppendField("TenantId", "【租户】不能为空")
	}
	if len(q.FileId) == 0 {
		ve.AppendField("FileId", "【文件ID】不能为空")
	}
	if len(q.CaseId) == 0 {
		ve.AppendField("CaseId", "【案件ID】不能为空")
	}
	if len(q.Name) == 0 {
		ve.AppendField("Name", "【工作表】名称不能为空")
	}
	return ve.GetError()
}
