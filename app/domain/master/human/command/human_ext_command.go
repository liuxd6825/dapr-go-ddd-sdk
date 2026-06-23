package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"

// HumanExtCreateCommand 创建命令
type HumanExtCreateCommand struct {
	CommandId string           `json:"commandId"`
	Data      []*model.HumanExt `json:"data"`
}

// HumanExtUpdateCommand 更新命令
type HumanExtUpdateCommand struct {
	CommandId string           `json:"commandId"`
	Data      []*model.HumanExt `json:"data"`
}

// HumanExtSubmitManyCommand 批量创建或更新命令
type HumanExtSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.HumanExt `json:"insertData"`
		UpdateData []*model.HumanExt `json:"updateData"`
	} `json:"data"`
}