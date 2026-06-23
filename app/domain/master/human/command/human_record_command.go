package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"

// HumanRecordCreateCommand 创建命令
type HumanRecordCreateCommand struct {
	CommandId string                `json:"commandId"`
	Data      []*model.HumanRecord `json:"data"`
}

// HumanRecordUpdateCommand 更新命令
type HumanRecordUpdateCommand struct {
	CommandId string                `json:"commandId"`
	Data      []*model.HumanRecord `json:"data"`
}

// HumanRecordSubmitManyCommand 批量创建或更新命令
type HumanRecordSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.HumanRecord `json:"insertData"`
		UpdateData []*model.HumanRecord `json:"updateData"`
	} `json:"data"`
}