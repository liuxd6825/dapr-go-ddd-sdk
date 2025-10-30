package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/case/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type CaseCreateCommand struct {
	xbase.Command[model.Case]
}
type CaseUpdateCommand struct {
	xbase.Command[model.Case]
}
type CaseDeleteCommand struct {
	xbase.DeleteByIdCommand
}
type CaseDeleteBatchCommand struct {
	xbase.Command[[]string]
}
