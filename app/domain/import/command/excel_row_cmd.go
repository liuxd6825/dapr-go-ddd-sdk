package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type ExcelRowCreateCommand struct {
	xbase.Command[*model.ExcelRow]
}

type ExcelRowCreateManyCommand struct {
	xbase.Command[[]*model.ExcelRow]
}

type ExcelRowUpdateCommand struct {
	xbase.Command[*model.ExcelRow]
}
type ExcelRowDeleteCommand struct {
	xbase.DeleteByIdCommand
}
