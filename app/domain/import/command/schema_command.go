package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
)

type SchemaCreateCommand struct {
	xbase.Command[model.SchemaModel]
}
type SchemaUpdateCommand struct {
	xbase.Command[model.SchemaModel]
}
type SchemaDeleteCommand struct {
	xbase.DeleteByIdCommand
}
