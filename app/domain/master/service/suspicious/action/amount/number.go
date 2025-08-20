package amount

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/suspicious/action"
)

type AmountNumber struct {
	rule model.AmountNumberRule
}

func NewAmountNumber(rule model.AmountNumberRule) *AmountNumber {
	return &AmountNumber{
		rule: rule,
	}
}

func (s *AmountNumber) IsEnable() bool {
	return s.rule.IsEnable
}

// DoAction 整数交易
func (s *AmountNumber) DoAction(ctx context.Context, tx *model.Record, txIndex int, accTxs *model.AccountRecords, result *action.AnalyseResult) {
	if int(tx.Amount)%s.rule.NumberValue == 0 {
		result.AddRecord(tx, model.SuType_AmountNumber, fmt.Sprintf("整数交易: 金额 %.2f", tx.Amount))
	}
}
