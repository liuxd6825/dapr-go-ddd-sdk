package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type CodeNewCommand struct {
	xbase.Command[CodeNewCommandData]
}

type CodeNewCommandData struct {
	Type string `json:"type" validate:"required"`
}

type CodeNewCommandResult struct {
	Code string `json:"code" validate:"required"`
}
