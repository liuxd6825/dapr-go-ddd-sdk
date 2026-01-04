package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type RecordPreviewCommand struct {
	xbase.Command[field.RecordPreviewCommandFields]
}

type RecordCreate4ExcelCommand struct {
	xbase.Command[field.RecordCreate4ExcelCommandFields]
}

type RecordImport2MasterCommand struct {
	xbase.Command[field.RecordImport2MasterFields]
}

// RecordCreateCommand
// @Description:
type RecordCreateCommand struct {
	xbase.Command[field.RecordCreateFields]
}

// RecordDeleteCommand
// @Description:
type RecordDeleteCommand struct {
	xbase.Command[field.RecordIeDeleteFields]
}

// RecordUpdateFilterCommand
// @Description:
type RecordUpdateFilterCommand struct {
	xbase.Command[field.RecordIeUpdateFilterFields]
}

// RecordUpdateCommand
// @Description:
type RecordUpdateCommand struct {
	xbase.Command[field.RecordIeUpdateFields]
}

// RecordUpdateFieldCommand
// @Description:
type RecordUpdateFieldCommand struct {
	xbase.Command[field.RecordIeUpdateFieldFields]
}

// RecordRecordCommand
// @Description: 撤销导入任务中的所有流水
type RecordRecordCommand struct {
	xbase.Command[RecordRecordCommandData]
}

type RecordRecordCommandData struct {
	TaskId string `json:"taskId"`
}
