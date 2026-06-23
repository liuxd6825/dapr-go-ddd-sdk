package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/model"

// ContractProductCreateCommand 创建命令
type ContractProductCreateCommand struct {
	CommandId string                  `json:"commandId"`
	Data      []*model.ContractProduct `json:"data"`
}

// ContractProductUpdateCommand 更新命令
type ContractProductUpdateCommand struct {
	CommandId string                  `json:"commandId"`
	Data      []*model.ContractProduct `json:"data"`
}

// ContractProductSubmitManyCommand 批量创建或更新命令
type ContractProductSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.ContractProduct `json:"insertData"`
		UpdateData []*model.ContractProduct `json:"updateData"`
	} `json:"data"`
}