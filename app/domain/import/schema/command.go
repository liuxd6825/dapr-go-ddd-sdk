package schema

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"

type SchemaCreateCommand = xbase.Command[SchemaModel]
type SchemaUpdateCommand = xbase.Command[SchemaModel]
type SchemaDeleteCommand = xbase.DeleteByIdCommand
