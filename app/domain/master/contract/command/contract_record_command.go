package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/model"

// ContractRecordCreateCommand 创建命令
type ContractRecordCreateCommand struct {
	CommandId string                 `json:"commandId"`
	Data      []*model.ContractRecord `json:"data"`
}

// ContractRecordUpdateCommand 更新命令
type ContractRecordUpdateCommand struct {
	CommandId string                 `json:"commandId"`
	Data      []*model.ContractRecord `json:"data"`
}

// ContractRecordSubmitManyCommand 批量创建或更新命令
type ContractRecordSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.ContractRecord `json:"insertData"`
		UpdateData []*model.ContractRecord `json:"updateData"`
	} `json:"data"`
}