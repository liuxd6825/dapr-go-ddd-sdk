package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
)

type SchemaCreateCommand = xbase.Command[model.SchemaModel]
type SchemaUpdateCommand = xbase.Command[model.SchemaModel]
type SchemaDeleteCommand = xbase.DeleteByIdCommand
