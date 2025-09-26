package amount

import (
	"context"
	"fmt"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/service/suspicious/action"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
)

type AmountNumber struct {
	rule model2.AmountNumberRule
}

func NewAmountNumber(rule model2.AmountNumberRule) *AmountNumber {
	return &AmountNumber{
		rule: rule,
	}
}

func (s *AmountNumber) IsEnable() bool {
	return s.rule.IsEnable
}

// DoAction 整数交易
func (s *AmountNumber) DoAction(ctx context.Context, tx *model.Record, txIndex int, accTxs *model2.AccountRecords, result *action.AnalyseResult) {
	if int(tx.Amount)%s.rule.Number == 0 {
		result.AddRecord(tx, model2.SuType_AmountNumber, fmt.Sprintf("整数交易: 金额 %.2f", tx.Amount))
	}
}
