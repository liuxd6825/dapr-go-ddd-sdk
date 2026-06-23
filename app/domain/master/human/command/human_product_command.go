package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"

// HumanProductCreateCommand 创建命令
type HumanProductCreateCommand struct {
	CommandId string                 `json:"commandId"`
	Data      []*model.HumanProduct `json:"data"`
}

// HumanProductUpdateCommand 更新命令
type HumanProductUpdateCommand struct {
	CommandId string                 `json:"commandId"`
	Data      []*model.HumanProduct `json:"data"`
}

// HumanProductSubmitManyCommand 批量创建或更新命令
type HumanProductSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.HumanProduct `json:"insertData"`
		UpdateData []*model.HumanProduct `json:"updateData"`
	} `json:"data"`
}