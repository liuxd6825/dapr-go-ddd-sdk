package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type FolderCreateCommand struct {
	xbase.Command[model.FolderView]
}

type FolderUpdateCommand struct {
	xbase.Command[model.Folder]
}

type FolderDeleteCommand struct {
	xbase.Command[model.Folder]
}

type FolderRenameCommand struct {
	xbase.Command[model.RenameFolder]
}

type FolderMoveCommand struct {
	xbase.Command[model.MoveFolder]
}
