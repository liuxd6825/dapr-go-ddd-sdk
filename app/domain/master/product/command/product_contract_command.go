package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/model"

// ProductContractCreateCommand 创建命令
type ProductContractCreateCommand struct {
	CommandId string                   `json:"commandId"`
	Data      []*model.ProductContract `json:"data"`
}

// ProductContractUpdateCommand 更新命令
type ProductContractUpdateCommand struct {
	CommandId string                   `json:"commandId"`
	Data      []*model.ProductContract `json:"data"`
}

// ProductContractSubmitManyCommand 批量创建或更新命令
type ProductContractSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.ProductContract `json:"insertData"`
		UpdateData []*model.ProductContract `json:"updateData"`
	} `json:"data"`
}