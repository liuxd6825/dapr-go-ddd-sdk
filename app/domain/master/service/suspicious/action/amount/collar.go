package amount

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/suspicious/action"
	"time"
)

type Collar struct {
	rule model.AmountCollarRule
}

func NewCollar(rule model.AmountCollarRule) *Collar {
	return &Collar{
		rule: rule,
	}
}

func (s *Collar) IsEnable() bool {
	return s.rule.IsEnable
}

// DoAction 对敲金额
func (s *Collar) DoAction(ctx context.Context, tx *model.Record, txIndex int, accTxs *model.AccountRecords, result *action.AnalyseResult) {
	if tx.Amount < s.rule.TxAmount {
		return
	}
	duration := time.Duration(s.rule.CollarDays) * 24 * time.Hour
	amount := tx.Amount * (s.rule.Percent / 100)
	txIO := tx.GetIO()
	for i := txIndex + 1; i < len(accTxs.Records); i++ {
		nextTx := accTxs.Records[i]
		if nextTx.Date.Sub(tx.Date) > duration {
			return
		}
		if txIO == nextTx.GetIO() { // 收付方式相同
			return
		}

		if amount < nextTx.Amount && tx.Name == nextTx.OppName && tx.OppName == nextTx.Name && tx.OppName != tx.Name {
			reason := fmt.Sprintf("%s与%s疑似生产“对敲交易”，金额：%.2f", tx.Date.Format(time.DateTime), nextTx.Date.Format(time.DateTime), tx.Amount)
			result.AddBatch(reason, model.SuType_AmountCollar, tx, nextTx)
		}
	}
}
