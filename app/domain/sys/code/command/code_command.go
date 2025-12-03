package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/code/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type CodeNewCommand struct {
	xbase.Command[CodeNewCommandData]
}

type CodeNewCommandData struct {
	Type   string          `json:"type" validate:"required"`
	Length int             `json:"length" validate:"required"`
	Style  model.CodeStyle `json:"style"`
}

type CodeNewCommandResult struct {
	Code string `json:"code" validate:"required"`
}

func NewCodeNewCommand() *CodeNewCommand {
	return &CodeNewCommand{}
}
