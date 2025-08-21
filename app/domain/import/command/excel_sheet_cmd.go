package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type ExcelSheetCreateCommand struct {
	xbase.Command[*model.ExcelSheet]
}

type ExcelSheetCreateManyCommand struct {
	xbase.Command[[]*model.ExcelSheet]
}

type ExcelSheetUpdateCommand struct {
	xbase.Command[*model.ExcelSheet]
}
type ExcelSheetDeleteCommand struct {
	xbase.DeleteByIdCommand
}
