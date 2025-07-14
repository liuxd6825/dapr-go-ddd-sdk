package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/model"

type TagCreateCommand struct {
	CommandId string    `json:"commandId"`
	Data      model.Tag `json:"data"`
}

type TagUpdateCommand struct {
	CommandId string    `json:"commandId"`
	Data      model.Tag `json:"data"`
}

type TagDeleteCommand struct {
	CommandId string    `json:"commandId"`
	Data      model.Tag `json:"data"`
}
