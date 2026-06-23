package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"

// HumanHumanCreateCommand 创建命令
type HumanHumanCreateCommand struct {
	CommandId string              `json:"commandId"`
	Data      []*model.HumanHuman `json:"data"`
}

// HumanHumanUpdateCommand 更新命令
type HumanHumanUpdateCommand struct {
	CommandId string              `json:"commandId"`
	Data      []*model.HumanHuman `json:"data"`
}

// HumanHumanSubmitManyCommand 批量创建或更新命令
type HumanHumanSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.HumanHuman `json:"insertData"`
		UpdateData []*model.HumanHuman `json:"updateData"`
	} `json:"data"`
}