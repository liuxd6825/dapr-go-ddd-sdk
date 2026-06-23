package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"

// HumanCredentialCreateCommand 创建命令
type HumanCredentialCreateCommand struct {
	CommandId string                   `json:"commandId"`
	Data      []*model.HumanCredential `json:"data"`
}

// HumanCredentialUpdateCommand 更新命令
type HumanCredentialUpdateCommand struct {
	CommandId string                   `json:"commandId"`
	Data      []*model.HumanCredential `json:"data"`
}

// HumanCredentialSubmitManyCommand 批量创建或更新命令
type HumanCredentialSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.HumanCredential `json:"insertData"`
		UpdateData []*model.HumanCredential `json:"updateData"`
	} `json:"data"`
}