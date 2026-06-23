package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/model"

// ContractHumanCreateCommand 创建命令
type ContractHumanCreateCommand struct {
	CommandId string                 `json:"commandId"`
	Data      []*model.ContractHuman `json:"data"`
}

// ContractHumanUpdateCommand 更新命令
type ContractHumanUpdateCommand struct {
	CommandId string                 `json:"commandId"`
	Data      []*model.ContractHuman `json:"data"`
}

// ContractHumanSubmitManyCommand 批量创建或更新命令
type ContractHumanSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.ContractHuman `json:"insertData"`
		UpdateData []*model.ContractHuman `json:"updateData"`
	} `json:"data"`
}