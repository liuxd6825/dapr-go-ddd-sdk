package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"

// HumanContractCreateCommand 创建命令
type HumanContractCreateCommand struct {
	CommandId string                 `json:"commandId"`
	Data      []*model.HumanContract `json:"data"`
}

// HumanContractUpdateCommand 更新命令
type HumanContractUpdateCommand struct {
	CommandId string                 `json:"commandId"`
	Data      []*model.HumanContract `json:"data"`
}

// HumanContractSubmitManyCommand 批量创建或更新命令
type HumanContractSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.HumanContract `json:"insertData"`
		UpdateData []*model.HumanContract `json:"updateData"`
	} `json:"data"`
}