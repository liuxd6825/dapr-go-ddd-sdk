package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/currency/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type CurrencyCreateCommand struct {
	xbase.Command[model.Currency]
}
type CurrencyCreateManyCommand struct {
	xbase.Command[[]*model.Currency]
}

type CurrencyUpdateCommand struct {
	xbase.Command[model.Currency]
}
type CurrencyDeleteCommand struct {
	xbase.DeleteByIdCommand
}
type CurrencyDeleteBatchCommand struct {
	xbase.Command[[]string]
}

type CurrencyDeleteAllCommand struct {
	xbase.Command[map[string]string]
}

type CurrencyUpdateAllRateCommand struct {
	xbase.Command[CurrencyUpdateAllRateData]
}

type CurrencyUpdateAllRateData struct {
	CnyRate float64 `json:"cnyRate"`
}
