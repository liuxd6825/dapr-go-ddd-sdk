package amount

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/suspicious/action"
)

type AmountLarge struct {
	rule model.AmountLargeRule
}

func NewAmountLarge(rule model.AmountLargeRule) *AmountLarge {
	a := &AmountLarge{
		rule: rule,
	}
	return a
}

func (s *AmountLarge) IsEnable() bool {
	return s.rule.IsEnable
}

func (s *AmountLarge) DoAction(ctx context.Context, tx *model.Record, txIndex int, accTxs *model.AccountRecords, result *action.AnalyseResult) {
	if tx.Amount > s.rule.LargeValue {
		result.AddRecord(tx, model.SuType_AmountLarge, fmt.Sprintf("大额交易: 金额 %.2f", tx.Amount))
	}
}
