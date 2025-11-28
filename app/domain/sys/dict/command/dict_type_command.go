package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/dict/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type DictTypeCreateCommand struct {
	xbase.Command[model.DictType]
}
type DictTypeUpdateCommand struct {
	xbase.Command[model.DictType]
}
type DictTypeDeleteCommand struct {
	xbase.DeleteByIdCommand
}
type DictTypeDeleteBatchCommand struct {
	xbase.Command[[]string]
}
