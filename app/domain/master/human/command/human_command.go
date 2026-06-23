package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"

// HumanCreateCommand 创建人员命令
type HumanCreateCommand struct {
	CommandId string      `json:"commandId"`
	Data      model.Human `json:"data"`
}

// HumanUpdateCommand 更新人员命令
type HumanUpdateCommand struct {
	CommandId string      `json:"commandId"`
	Data      model.Human `json:"data"`
}

// HumanSubmitCommand 创建或更新人员命令
type HumanSubmitCommand struct {
	CommandId string      `json:"commandId"`
	Data      model.Human `json:"data"`
}

// HumanDeleteByIdsCommand 批量删除人员命令
type HumanDeleteByIdsCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		Ids    []string `json:"ids"`
		CaseId string   `json:"caseId"`
	} `json:"data"`
}