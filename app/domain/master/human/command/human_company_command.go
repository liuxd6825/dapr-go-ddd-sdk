package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"

// HumanCompanyCreateCommand 创建命令
type HumanCompanyCreateCommand struct {
	CommandId string                `json:"commandId"`
	Data      []*model.HumanCompany `json:"data"`
}

// HumanCompanyUpdateCommand 更新命令
type HumanCompanyUpdateCommand struct {
	CommandId string                `json:"commandId"`
	Data      []*model.HumanCompany `json:"data"`
}

// HumanCompanySubmitManyCommand 批量创建或更新命令
type HumanCompanySubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.HumanCompany `json:"insertData"`
		UpdateData []*model.HumanCompany `json:"updateData"`
	} `json:"data"`
}