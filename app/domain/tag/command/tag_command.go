package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type TagCreateCommand struct {
	xbase.Command[model.Tag]
}

type TagUpdateCommand struct {
	xbase.Command[model.Tag]
}

type TagDeleteCommand struct {
	xbase.DeleteByIdCommand
}
