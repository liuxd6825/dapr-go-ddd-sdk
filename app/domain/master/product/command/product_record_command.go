package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/model"

// ProductRecordCreateCommand 创建命令
type ProductRecordCreateCommand struct {
	CommandId string                `json:"commandId"`
	Data      []*model.ProductRecord `json:"data"`
}

// ProductRecordUpdateCommand 更新命令
type ProductRecordUpdateCommand struct {
	CommandId string                `json:"commandId"`
	Data      []*model.ProductRecord `json:"data"`
}

// ProductRecordSubmitManyCommand 批量创建或更新命令
type ProductRecordSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.ProductRecord `json:"insertData"`
		UpdateData []*model.ProductRecord `json:"updateData"`
	} `json:"data"`
}