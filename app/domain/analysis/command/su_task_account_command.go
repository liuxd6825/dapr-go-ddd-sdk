package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type SuTaskAccountCreateCommand struct {
	xbase.Command[model.SuTaskAccount]
}
type SuTaskAccountUpdateCommand struct {
	xbase.Command[model.SuTaskAccount]
}
