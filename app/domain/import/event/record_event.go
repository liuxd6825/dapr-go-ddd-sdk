package event

import (
	"fmt"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/events"
)

// RecordImportMasterEvent
// @Description:
type RecordImportMasterEvent struct {
	events.Event[RecordImportMasterEventData]
}

const RecordImportMasterEventType = "RecordImportMasterEvent"

type RecordImportMasterEventData struct {
	Date       *time.Time            `json:"date" gorm:"date" bson:"date" validate:"-" title:"时间"`
	CaseId     string                `json:"caseId" gorm:"case_id" bson:"case_id" validate:"required"  title:"案件ID"`         // 案件Id
	DocId      string                `json:"docId" gorm:"doc_Id" bson:"doc_id" validate:"required"  title:"文档ID"`            // 文档Id
	FileId     string                `json:"fileId" gorm:"file_id" bson:"file_id" validate:"required"  title:"文件ID"`         // 文档Id
	FileName   string                `json:"fileName" gorm:"file_name" bson:"file_name" validate:"required"  title:"文件名称"`   // 文件名称
	SheetId    string                `json:"sheetId" gorm:"sheet_id" bson:"sheet_id" validate:"required"  title:"工作表ID"`     // 工作表
	SheetName  string                `json:"sheetName" gorm:"sheet_name" bson:"sheet_name" validate:"required"  title:"工作表"` // 工作表
	TaskId     string                `json:"taskId" gorm:"task_id" bson:"task_id"  validate:"required" title:"任务ID" `        // 任务ID
	Items      []*field.RecordFields `json:"items" gorm:"item" bson:"items" validate:"required"  title:"流水明细"`               // 流水明细
	IsAddItems bool                  `json:"isAddItems"`
	MasterType string                `json:"masterType" gorm:"master_type" title:"主数据类型"`
	MasterId   string                `json:"masterId" gorm:"master_id" bson:"master_id" title:"主数据Id"`
}

func (f RecordImportMasterEventData) Id() string {
	return fmt.Sprintf("%s:%s:%s:%s", f.TaskId, f.CaseId, f.FileId, f.SheetName)
}
