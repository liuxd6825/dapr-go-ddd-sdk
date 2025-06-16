package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/model"

type MessageCreateCommand struct {
	CommandId string        `json:"commandId"`
	Data      model.Message `json:"data"`
}

type MessageUpdateCommand struct {
	CommandId string        `json:"commandId"`
	Data      model.Message `json:"data"`
}
