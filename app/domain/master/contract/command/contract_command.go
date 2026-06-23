package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/model"

// ContractCreateCommand 创建合同命令
type ContractCreateCommand struct {
	CommandId string         `json:"commandId"`
	Data      model.Contract `json:"data"`
}

// ContractUpdateCommand 更新合同命令
type ContractUpdateCommand struct {
	CommandId string         `json:"commandId"`
	Data      model.Contract `json:"data"`
}

// ContractSubmitCommand 创建或更新合同命令
type ContractSubmitCommand struct {
	CommandId string         `json:"commandId"`
	Data      model.Contract `json:"data"`
}

// ContractDeleteByIdsCommand 批量删除合同命令
type ContractDeleteByIdsCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		Ids    []string `json:"ids"`
		CaseId string   `json:"caseId"`
	} `json:"data"`
}