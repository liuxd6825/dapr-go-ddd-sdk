package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/model"

// ProductProductCreateCommand 创建命令
type ProductProductCreateCommand struct {
	CommandId string                  `json:"commandId"`
	Data      []*model.ProductProduct `json:"data"`
}

// ProductProductUpdateCommand 更新命令
type ProductProductUpdateCommand struct {
	CommandId string                  `json:"commandId"`
	Data      []*model.ProductProduct `json:"data"`
}

// ProductProductSubmitManyCommand 批量创建或更新命令
type ProductProductSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.ProductProduct `json:"insertData"`
		UpdateData []*model.ProductProduct `json:"updateData"`
	} `json:"data"`
}