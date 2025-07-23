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
	Id        string `json:"id"   title:"主键" validate:"required"`
	CaseId    string `json:"caseId" title:"案件Id"  validate:"required" `
	FileName  string `json:"fileName" title:"文件名" validate:"required" `
	SheetName string `json:"sheetName" title:"工作页" validate:"required" `
	DocId     string `json:"docId" title:"文档Id"  validate:"required" `
	DocFileId string `json:"docFileId" title:"文档文件Id" validate:"required" `
}

type ExcelCreateRowsCommandData struct {
	CaseId    string `json:"caseId" title:"案件Id" validate:"required" `
	DocId     string `json:"docId" title:"文档Id" validate:"required" `
	FileId    string `json:"fileId" title:"文件Id" validate:"required" `
	FileName  string `json:"fileName" title:"文件名" validate:"required" `
	SheetName string `json:"sheetName" title:"工作页" validate:"required" `
}

type ExcelFindRecordsQuery struct {
	CaseId    string                `json:"caseId"  validate:"required" `
	DocId     string                `json:"docId" validate:"required" `
	FileId    string                `json:"fileId" validate:"required" `
	FileName  string                `json:"fileName" validate:"required" `
	SheetName string                `json:"sheetName" validate:"required" `
	MaxRows   *int64                `json:"maxRows" validate:"required" `
	Template  *model.RecordTemplate `json:"template"`
}
