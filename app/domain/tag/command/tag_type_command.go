package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/model"

type TagTypeCreateCommand struct {
	CommandId string        `json:"commandId"`
	Data      model.TagType `json:"data"`
}

type TagTypeUpdateCommand struct {
	CommandId string        `json:"commandId"`
	Data      model.TagType `json:"data"`
}

type TagTypeDeleteCommand struct {
	CommandId string        `json:"commandId"`
	Data      model.TagType `json:"data"`
}
