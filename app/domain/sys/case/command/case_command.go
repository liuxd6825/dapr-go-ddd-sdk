package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/case/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
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
	xbase.Command[DeleteByIds]
}

type DeleteByIds struct {
	Ids    []string `json:"ids"`
	CaseId string   `json:"caseId"`
}
