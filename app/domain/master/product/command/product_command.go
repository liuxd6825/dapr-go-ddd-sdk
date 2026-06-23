package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/model"

// ProductCreateCommand 创建产品命令
type ProductCreateCommand struct {
	CommandId string       `json:"commandId"`
	Data      model.Product `json:"data"`
}

// ProductUpdateCommand 更新产品命令
type ProductUpdateCommand struct {
	CommandId string       `json:"commandId"`
	Data      model.Product `json:"data"`
}

// ProductSubmitCommand 创建或更新产品命令
type ProductSubmitCommand struct {
	CommandId string       `json:"commandId"`
	Data      model.Product `json:"data"`
}

// ProductDeleteByIdsCommand 批量删除产品命令
type ProductDeleteByIdsCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		Ids    []string `json:"ids"`
		CaseId string   `json:"caseId"`
	} `json:"data"`
}