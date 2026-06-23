package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"

// HumanReportedAmountCreateCommand 创建命令（单实体）
type HumanReportedAmountCreateCommand struct {
	CommandId string                       `json:"commandId"`
	Data      model.HumanReportedAmount   `json:"data"`
}

// HumanReportedAmountUpdateCommand 更新命令
type HumanReportedAmountUpdateCommand struct {
	CommandId string                       `json:"commandId"`
	Data      model.HumanReportedAmount   `json:"data"`
}

// HumanReportedAmountSubmitCommand 创建或更新命令
type HumanReportedAmountSubmitCommand struct {
	CommandId string                       `json:"commandId"`
	Data      model.HumanReportedAmount   `json:"data"`
}