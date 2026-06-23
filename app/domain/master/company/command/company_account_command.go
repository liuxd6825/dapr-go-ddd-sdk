package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/model"
)


// CompanyAccountCreateCommand 创建命令
type CompanyAccountCreateCommand struct {
	CommandId string                  `json:"commandId"`
	Data      []*model.CompanyAccount `json:"data"`
}

// CompanyAccountUpdateCommand 更新命令
type CompanyAccountUpdateCommand struct {
	CommandId string                  `json:"commandId"`
	Data      []*model.CompanyAccount `json:"data"`
}

// CompanyAccountSubmitManyCommand 批量创建或更新命令
type CompanyAccountSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.CompanyAccount `json:"insertData"`
		UpdateData []*model.CompanyAccount `json:"updateData"`
	} `json:"data"`
}