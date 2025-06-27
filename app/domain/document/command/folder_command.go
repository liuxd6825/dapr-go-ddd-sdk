package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"

type FolderCreateCommand struct {
	CommandId string       `json:"commandId"`
	Data      model.Folder `json:"data"`
}

type FolderUpdateCommand struct {
	CommandId string       `json:"commandId"`
	Data      model.Folder `json:"data"`
}

type FolderDeleteCommand struct {
	CommandId string       `json:"commandId"`
	Data      model.Folder `json:"data"`
}

type FolderRenameCommand struct {
	CommandId string             `json:"commandId"`
	Data      model.RenameFolder `json:"data"`
}
