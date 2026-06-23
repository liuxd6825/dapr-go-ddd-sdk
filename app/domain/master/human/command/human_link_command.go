package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"

// HumanLinkCreateCommand 创建命令
type HumanLinkCreateCommand struct {
	CommandId string            `json:"commandId"`
	Data      []*model.HumanLink `json:"data"`
}

// HumanLinkUpdateCommand 更新命令
type HumanLinkUpdateCommand struct {
	CommandId string            `json:"commandId"`
	Data      []*model.HumanLink `json:"data"`
}

// HumanLinkSubmitManyCommand 批量创建或更新命令
type HumanLinkSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.HumanLink `json:"insertData"`
		UpdateData []*model.HumanLink `json:"updateData"`
	} `json:"data"`
}