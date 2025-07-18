package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
)

type RecordCreate4ExcelCommand = xbase.Command[field.RecordCreate4ExcelCommandFields]

type RecordImport2MasterAppCmd = xbase.Command[field.RecordImport2MasterFields]

// RecordCreateCommand
// @Description:
type RecordCreateCommand = xbase.Command[field.RecordCreateFields]

// RecordDeleteCommand
// @Description:
type RecordDeleteCommand = xbase.Command[field.RecordIeDeleteFields]

// RecordUpdateFilterCommand
// @Description:
type RecordUpdateFilterCommand = xbase.Command[field.RecordIeUpdateFilterFields]

// RecordUpdateCommand
// @Description:
type RecordUpdateCommand = xbase.Command[field.RecordIeUpdateFields]

// RecordUpdateFieldCommand
// @Description:
type RecordUpdateFieldCommand = xbase.Command[field.RecordIeUpdateFieldFields]

// RecordCreateManyFromExcelCommand
// @Description:
type RecordCreateManyFromExcelCommand = xbase.Command[field.RecordCreateManyFromExcelFields]
