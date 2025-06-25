package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"

type FileCreateCommand struct {
	CommandId string     `json:"commandId"`
	Data      model.File `json:"data"`
}

type FileUpdateCommand struct {
	CommandId string     `json:"commandId"`
	Data      model.File `json:"data"`
}

type FileDeleteCommand struct {
	CommandId string     `json:"commandId"`
	Data      model.File `json:"data"`
}
