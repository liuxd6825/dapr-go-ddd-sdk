package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"

// HumanAddressCreateCommand 创建命令
type HumanAddressCreateCommand struct {
	CommandId string                `json:"commandId"`
	Data      []*model.HumanAddress `json:"data"`
}

// HumanAddressUpdateCommand 更新命令
type HumanAddressUpdateCommand struct {
	CommandId string                `json:"commandId"`
	Data      []*model.HumanAddress `json:"data"`
}

// HumanAddressSubmitManyCommand 批量创建或更新命令
type HumanAddressSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.HumanAddress `json:"insertData"`
		UpdateData []*model.HumanAddress `json:"updateData"`
	} `json:"data"`
}