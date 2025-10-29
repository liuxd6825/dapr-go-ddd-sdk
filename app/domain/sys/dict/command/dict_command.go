package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/dict/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type DictCreateCommand struct {
	xbase.Command[model.Dict]
}
type DictUpdateCommand struct {
	xbase.Command[model.Dict]
}
type DictDeleteCommand struct {
	xbase.DeleteByIdCommand
}
type DictDeleteBatchCommand struct {
	xbase.Command[[]string]
}
