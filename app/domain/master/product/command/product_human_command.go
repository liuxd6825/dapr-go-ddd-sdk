package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/model"

// ProductHumanCreateCommand 创建命令
type ProductHumanCreateCommand struct {
	CommandId string                `json:"commandId"`
	Data      []*model.ProductHuman `json:"data"`
}

// ProductHumanUpdateCommand 更新命令
type ProductHumanUpdateCommand struct {
	CommandId string                `json:"commandId"`
	Data      []*model.ProductHuman `json:"data"`
}

// ProductHumanSubmitManyCommand 批量创建或更新命令
type ProductHumanSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.ProductHuman `json:"insertData"`
		UpdateData []*model.ProductHuman `json:"updateData"`
	} `json:"data"`
}