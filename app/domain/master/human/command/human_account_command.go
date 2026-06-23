package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"

// HumanAccountCreateCommand 创建命令
type HumanAccountCreateCommand struct {
	CommandId string                `json:"commandId"`
	Data      []*model.HumanAccount `json:"data"`
}

// HumanAccountUpdateCommand 更新命令
type HumanAccountUpdateCommand struct {
	CommandId string                `json:"commandId"`
	Data      []*model.HumanAccount `json:"data"`
}

// HumanAccountSubmitManyCommand 批量创建或更新命令
type HumanAccountSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.HumanAccount `json:"insertData"`
		UpdateData []*model.HumanAccount `json:"updateData"`
	} `json:"data"`
}