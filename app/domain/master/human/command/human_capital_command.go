package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"

// HumanCapitalCreateCommand 创建命令
type HumanCapitalCreateCommand struct {
	CommandId string                `json:"commandId"`
	Data      []*model.HumanCapital `json:"data"`
}

// HumanCapitalUpdateCommand 更新命令
type HumanCapitalUpdateCommand struct {
	CommandId string                `json:"commandId"`
	Data      []*model.HumanCapital `json:"data"`
}

// HumanCapitalSubmitManyCommand 批量创建或更新命令
type HumanCapitalSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.HumanCapital `json:"insertData"`
		UpdateData []*model.HumanCapital `json:"updateData"`
	} `json:"data"`
}