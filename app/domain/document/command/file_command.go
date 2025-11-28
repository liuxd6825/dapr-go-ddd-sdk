package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type FileCreateCommand struct {
	xbase.Command[model.File]
}

type FileUpdateCommand struct {
	xbase.Command[model.File]
}

type FileDeleteCommand struct {
	xbase.Command[model.File]
}
