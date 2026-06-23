package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"

// HumanSuspectAmountCreateCommand 创建命令（单实体）
type HumanSuspectAmountCreateCommand struct {
	CommandId string                      `json:"commandId"`
	Data      model.HumanSuspectAmount   `json:"data"`
}

// HumanSuspectAmountUpdateCommand 更新命令
type HumanSuspectAmountUpdateCommand struct {
	CommandId string                      `json:"commandId"`
	Data      model.HumanSuspectAmount   `json:"data"`
}

// HumanSuspectAmountSubmitCommand 创建或更新命令
type HumanSuspectAmountSubmitCommand struct {
	CommandId string                      `json:"commandId"`
	Data      model.HumanSuspectAmount   `json:"data"`
}