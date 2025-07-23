package event

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
)

// RecordImportMasterEvent
// @Description:
type RecordImportMasterEvent struct {
	xbase.Command[RecordImportMasterEventData]
}

type RecordImportMasterEventData struct {
	CaseId     string                `json:"caseId" validate:"required"  title:"案件ID"`   // 案件Id
	DocId      string                `json:"docId" validate:"required"  title:"文档ID"`    // 文档Id
	FileId     string                `json:"fileId" validate:"required"  title:"文件ID"`   // 文档Id
	FileName   string                `json:"fileName" validate:"required"  title:"文件名称"` // 文件名称
	SheetName  string                `json:"sheetName" validate:"required"  title:"工作表"` // 工作表
	TaskId     string                `json:"taskId"  validate:"required" title:"任务ID" `  // 任务ID
	Items      []*field.RecordFields `json:"records" validate:"required"  title:"流水明细"`  // 流水明细
	IsAddItems bool                  `json:"isAddItems"`
}

func (f RecordImportMasterEventData) Id() string {
	return fmt.Sprintf("%s:%s:%s:%s", f.TaskId, f.CaseId, f.FileId, f.SheetName)
}
