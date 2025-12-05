package amount

import (
	"context"
	"fmt"

	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/service/suspicious/action"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
)

// AmountLarge
// @Description: 大额金额
type AmountLarge struct {
	rule model2.AmountLargeRule
}

func NewAmountLarge(rule model2.AmountLargeRule) *AmountLarge {
	a := &AmountLarge{
		rule: rule,
	}
	return a
}

func (s *AmountLarge) IsEnable() bool {
	return s.rule.IsEnable
}

func (s *AmountLarge) DoAction(ctx context.Context, tx *model.Record, txIndex int, accTxs *model2.AccountRecords, result *action.AnalyseResult) {
	if tx.Amount > s.rule.LargeValue {
		result.AddRecord(tx, model2.SuType_AmountLarge, fmt.Sprintf("大额交易: 金额 %.2f", tx.Amount))
	}
}
