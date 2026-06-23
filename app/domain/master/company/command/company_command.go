package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/model"
)

// CompanyCreateCommand 创建公司命令
type CompanyCreateCommand struct {
	CommandId string        `json:"commandId"`
	Data      model.Company `json:"data"`
}

// CompanyUpdateCommand 更新公司命令
type CompanyUpdateCommand struct {
	CommandId string        `json:"commandId"`
	Data      model.Company `json:"data"`
}

// CompanySubmitCommand 创建或更新公司命令
type CompanySubmitCommand struct {
	CommandId string        `json:"commandId"`
	Data      model.Company `json:"data"`
}

// CompanyDeleteByIdsCommand 批量删除公司命令
type CompanyDeleteByIdsCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		Ids    []string `json:"ids"`
		CaseId string   `json:"caseId"`
	} `json:"data"`
}