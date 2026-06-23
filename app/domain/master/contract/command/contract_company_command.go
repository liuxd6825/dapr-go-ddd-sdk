package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/model"

// ContractCompanyCreateCommand 创建命令
type ContractCompanyCreateCommand struct {
	CommandId string                    `json:"commandId"`
	Data      []*model.ContractCompany `json:"data"`
}

// ContractCompanyUpdateCommand 更新命令
type ContractCompanyUpdateCommand struct {
	CommandId string                    `json:"commandId"`
	Data      []*model.ContractCompany `json:"data"`
}

// ContractCompanySubmitManyCommand 批量创建或更新命令
type ContractCompanySubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.ContractCompany `json:"insertData"`
		UpdateData []*model.ContractCompany `json:"updateData"`
	} `json:"data"`
}