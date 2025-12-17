package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type TagTypeCreateCommand struct {
	xbase.Command[model.TagType]
}

type TagTypeUpdateCommand struct {
	xbase.Command[model.TagType]
}

type TagTypeDeleteCommand struct {
	xbase.DeleteByIdCommand
}
