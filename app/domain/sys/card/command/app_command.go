package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type AppCreateCommand struct {
	xbase.Command[model.App]
}
type AppUpdateCommand struct {
	xbase.Command[model.App]
}
type AppDeleteCommand struct {
	xbase.DeleteByIdCommand
}
