package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type FunCreateCommand struct {
	xbase.Command[model.Fun]
}
type FunUpdateCommand struct {
	xbase.Command[model.Fun]
}
type FunDeleteCommand struct {
	xbase.DeleteByIdCommand
}
