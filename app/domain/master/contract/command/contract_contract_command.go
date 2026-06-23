package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/model"

// ContractContractCreateCommand 创建命令
type ContractContractCreateCommand struct {
	CommandId string                    `json:"commandId"`
	Data      []*model.ContractContract `json:"data"`
}

// ContractContractUpdateCommand 更新命令
type ContractContractUpdateCommand struct {
	CommandId string                    `json:"commandId"`
	Data      []*model.ContractContract `json:"data"`
}

// ContractContractSubmitManyCommand 批量创建或更新命令
type ContractContractSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.ContractContract `json:"insertData"`
		UpdateData []*model.ContractContract `json:"updateData"`
	} `json:"data"`
}