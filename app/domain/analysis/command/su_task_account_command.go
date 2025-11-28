package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type SuTaskAccountCreateCommand struct {
	xbase.Command[model.SuTaskAccount]
}
type SuTaskAccountUpdateCommand struct {
	xbase.Command[model.SuTaskAccount]
}
type SuTaskAccountBatchCreateCommand struct {
	xbase.Command[[]*model.SuTaskAccount]
}
type SuTaskAccountBatchDeleteCommand struct {
	xbase.Command[[]string]
}
