package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type CardFileCreateCommand struct {
	xbase.Command[model.CardFile]
}
type CardFileUpdateCommand struct {
	xbase.Command[model.CardFile]
}
type CardFileDeleteCommand struct {
	xbase.DeleteByIdCommand
}
