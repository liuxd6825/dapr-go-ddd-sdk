package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/model"
)

type ChatCreateCommand struct {
	CommandId string     `json:"commandId"`
	Data      model.Chat `json:"data"`
}

type ChatUpdateCommand struct {
	CommandId string     `json:"commandId"`
	Data      model.Chat `json:"data"`
}

type ChatDeleteCommand struct {
	CommandId string     `json:"commandId"`
	Data      model.Chat `json:"data"`
}
