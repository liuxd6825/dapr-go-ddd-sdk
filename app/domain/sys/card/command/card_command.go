package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type CardCreateCommand struct {
	xbase.Command[model.Card]
}
type CardUpdateCommand struct {
	xbase.Command[model.Card]
}
type CardDeleteCommand struct {
	xbase.DeleteByIdCommand
}
