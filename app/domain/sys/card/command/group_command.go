package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type GroupCreateCommand struct {
	xbase.Command[model.Group]
}
type GroupUpdateCommand struct {
	xbase.Command[model.Group]
}
type GroupDeleteCommand struct {
	xbase.DeleteByIdCommand
}
