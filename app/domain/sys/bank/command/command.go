package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/bank/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type BankCreateCommand struct {
	xbase.Command[model.Bank]
}
type BankCreateManyCommand struct {
	xbase.Command[[]*model.Bank]
}

type BankUpdateCommand struct {
	xbase.Command[model.Bank]
}
type BankDeleteCommand struct {
	xbase.DeleteByIdCommand
}
type BankDeleteBatchCommand struct {
	xbase.Command[[]string]
}

type BankDeleteAllCommand struct {
	xbase.Command[map[string]string]
}

type BankUpdateAllRateCommand struct {
	xbase.Command[BankUpdateAllRateData]
}

type BankUpdateAllRateData struct {
	CnyRate float64 `json:"cnyRate"`
}
