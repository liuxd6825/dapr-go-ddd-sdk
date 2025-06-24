package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/model"

type DocumentCreateCommand struct {
	CommandId string         `json:"commandId"`
	Data      model.Document `json:"data"`
}

type DocumentUpdateCommand struct {
	CommandId string         `json:"commandId"`
	Data      model.Document `json:"data"`
}
