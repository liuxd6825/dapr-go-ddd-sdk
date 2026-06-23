package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/model"

// ProductCompanyCreateCommand 创建命令
type ProductCompanyCreateCommand struct {
	CommandId string                  `json:"commandId"`
	Data      []*model.ProductCompany `json:"data"`
}

// ProductCompanyUpdateCommand 更新命令
type ProductCompanyUpdateCommand struct {
	CommandId string                  `json:"commandId"`
	Data      []*model.ProductCompany `json:"data"`
}

// ProductCompanySubmitManyCommand 批量创建或更新命令
type ProductCompanySubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.ProductCompany `json:"insertData"`
		UpdateData []*model.ProductCompany `json:"updateData"`
	} `json:"data"`
}