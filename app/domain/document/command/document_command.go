package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type DocumentCreateCommand struct {
	xbase.Command[model.Document]
}

type DocumentUpdateCommand struct {
	xbase.Command[model.Document]
}

type DocumentDeleteCommand struct {
	xbase.Command[model.DeleteDocument]
}

type DocumentRenameCommand struct {
	xbase.Command[model.RenameDocument]
}

type DocumentMoveCommand struct {
	xbase.Command[model.MoveDocument]
}
