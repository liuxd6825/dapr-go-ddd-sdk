package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/example/sys/pkg/xinfra/enums"

type FuncAddCommand struct {
	Base
	ModuleType     enums.ModuleType
	FunctionEntity model.FunctionEntity
}
