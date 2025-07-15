package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/recordie/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/recordie/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

type FindExcelViewAppCmd struct {
	TenantId  string `json:"tenantId" desc:"租户Id"`
	CaseId    string `json:"caseId" desc:"案件Id"`
	FileId    string `json:"fileId" desc:"文件Id"`
	FileName  string `json:"fileName" desc:"文件名"`
	SheetName string `json:"sheetName" desc:"Excel页名"`
	MaxRows   *int64 `json:"maxRows" desc:"最大记录数"`
}

type CreateFileAppCmd struct {
	TenantId string `json:"tenantId" desc:"租户Id"`
	CaseId   string `json:"caseId" desc:"案件Id"`
	DocId    string `json:"docId" desc:"文档Id"`
	FileId   string `json:"fileId" desc:"文件Id"`
	FileName string `json:"fileName" desc:"文件名"`
}

type CreateRowsAppCmd struct {
	TenantId  string `json:"tenantId" desc:"租户Id"`
	CaseId    string `json:"caseId" desc:"案件Id"`
	DocId     string `json:"docId" desc:"文档Id"`
	FileId    string `json:"fileId" desc:"文件Id"`
	FileName  string `json:"fileName" desc:"文件名"`
	SheetName string `json:"sheetName" desc:"工作页"`
}

type FindFileAppQuery struct {
	TenantId string `json:"tenantId" desc:"租户Id"`
	CaseId   string `json:"caseId" desc:"案件Id"`
	DocId    string `json:"docId" desc:"文档Id"`
	FileId   string `json:"fileId" desc:"文件Id"`
}

type FindRowsAppQuery struct {
	TenantId  string `json:"tenantId" desc:"租户Id"`
	CaseId    string `json:"caseId" desc:"案件Id"`
	DocId     string `json:"docId" desc:"文档Id"`
	FileId    string `json:"fileId" desc:"文件Id"`
	SheetName string `json:"sheetName" desc:"工作页"`
}

type FindRowsAppResult struct {
	Columns   []string         `json:"columns" desc:"列头"`
	SheetName string           `json:"sheetName" desc:"工作表"`
	MaxRow    int64            `json:"maxRow" desc:"最大行"`
	MaxCol    int64            `json:"maxCol" desc:"最大列"`
	Rows      []map[string]any `json:"rows" desc:"数据行"`
}

type FindRecordsAppQuery struct {
	TenantId  string                `json:"tenantId"`
	CaseId    string                `json:"caseId"`
	DocId     string                `json:"docId"`
	FileId    string                `json:"fileId"`
	FileName  string                `json:"fileName"`
	SheetName string                `json:"sheetName"`
	MaxRows   *int64                `json:"maxRows"`
	Template  *model.RecordTemplate `json:"template"`
}

type FindSheetBySheetNameAppQuery = query.FindSheetByNameQuery

func (q *CreateFileAppCmd) Validate() error {
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
	return ve.GetError()
}

func (q *FindFileAppQuery) Validate() error {
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
	return ve.GetError()
}

func (c *FindRecordsAppQuery) Validate() error {
	ve := errors.NewVerifyError()
	if len(c.TenantId) == 0 {
		ve.AppendField("TenantId", "【租户】不能为空")
	}
	if len(c.FileId) == 0 {
		ve.AppendField("FileId", "【文件ID】不能为空")
	}
	if len(c.FileName) == 0 {
		ve.AppendField("FileName", "【文件名称】不能为空")
	}
	if c.Template == nil {
		ve.AppendField("Template", "【模板】不能为空")
	}
	return ve.GetError()
}

func (c *FindRowsAppQuery) Validate() error {
	ve := errors.NewVerifyError()
	if len(c.TenantId) == 0 {
		ve.AppendField("TenantId", "【租户】不能为空")
	}
	if len(c.FileId) == 0 {
		ve.AppendField("FileId", "【文件ID】不能为空")
	}
	if len(c.SheetName) == 0 {
		ve.AppendField("Name", "【Name】不能为空")
	}
	return ve.GetError()
}

func (c *CreateRowsAppCmd) Validate() error {
	ve := errors.NewVerifyError()
	if len(c.TenantId) == 0 {
		ve.AppendField("TenantId", "【租户】不能为空")
	}
	if len(c.DocId) == 0 {
		ve.AppendField("DocId", "【文档ID】不能为空")
	}
	if len(c.CaseId) == 0 {
		ve.AppendField("CaseId", "【案件ID】不能为空")
	}
	if len(c.FileId) == 0 {
		ve.AppendField("FileId", "【文件ID】不能为空")
	}
	if len(c.FileName) == 0 {
		ve.AppendField("FileName", "【文件名称】不能为空")
	}
	return ve.GetError()
}
