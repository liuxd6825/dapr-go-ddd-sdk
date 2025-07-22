package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
)

type ExcelFileCreateCommand struct {
	xbase.Command[ExcelCreateByFileCommandData]
}
type ExcelFileUpdateCommand struct {
	xbase.Command[*model.ExcelFile]
}
type ExcelFileDeleteCommand struct {
	xbase.Command[*model.ExcelFile]
}

type ExcelFileDeleteByIdCommand struct {
	xbase.DeleteByIdCommand
}

type ExcelCreateRowsCommand struct {
	xbase.Command[ExcelCreateRowsCommandData]
}

type ExcelCreateByFileCommandData struct {
	Id        string `json:"id" validate:"required"  title:"主键"`
	CaseId    string `json:"caseId" title:"案件Id"  required:"true"`
	FileName  string `json:"fileName" title:"文件名" required:"true"`
	SheetName string `json:"sheetName" title:"工作页" required:"true"`
	DocId     string `json:"docId" title:"文档Id"  required:"true"`
	DocFileId string `json:"docFileId" title:"文档文件Id" required:"true"`
}

type ExcelCreateRowsCommandData struct {
	CaseId    string `json:"caseId" title:"案件Id" required:"true"`
	DocId     string `json:"docId" title:"文档Id" required:"true"`
	FileId    string `json:"fileId" title:"文件Id" required:"true"`
	FileName  string `json:"fileName" title:"文件名" required:"true"`
	SheetName string `json:"sheetName" title:"工作页"`
}

type ExcelFindRecordsQuery struct {
	CaseId    string                `json:"caseId"  required:"true"`
	DocId     string                `json:"docId" required:"true"`
	FileId    string                `json:"fileId" required:"true"`
	FileName  string                `json:"fileName" required:"true"`
	SheetName string                `json:"sheetName" required:"true"`
	MaxRows   *int64                `json:"maxRows" required:"true"`
	Template  *model.RecordTemplate `json:"template"`
}
