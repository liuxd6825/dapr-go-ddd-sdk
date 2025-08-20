package frequency

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/suspicious/action"
	"sync"
)

// High
// 高频交易对手方
// 统计与各对手方的交易次数。对于非主要供应商或客户，如果出现与其业务性质不符的高频次交易，应重点关注
type High struct {
	rule             model.FreqHighRule
	frequencyCounter *sync.Map // Key: string (e.g., "CounterpartyName#YYYY-MM"), Value: *atomic.Int64
}

func NewHigh(rule model.FreqHighRule) *High {
	return &High{
		rule:             rule,
		frequencyCounter: &sync.Map{},
	}
}

func (s *High) IsEnable() bool {
	return s.rule.IsEnable
}

func (s *High) DoAction(ctx context.Context, tx *model.Record, txIndex int, accTxs *model.AccountRecords, result *action.AnalyseResult) {
	// Exclude known core partners
	// In a real system, this list would come from a pre-loaded map.
	if tx.Name == accTxs.OwnerName { // Assuming BigCorp is a core partner
		return
	}

	// Determine the period key based on config (e.g., "2024-03" for Month)
	var periodKey string
	switch s.rule.HighPeriod {
	case model.SuFreqHighPeriod_Quarter:
		quarter := (tx.Date.Month()-1)/3 + 1
		periodKey = fmt.Sprintf("%d-Q%d", tx.Date.Year(), quarter)
	case model.SuFreqHighPeriod_Year:
		periodKey = fmt.Sprintf("%d", tx.Date.Year())
	case model.SuFreqHighPeriod_Month:
		periodKey = tx.Date.Format("2006-01")
	}

	mapKey := fmt.Sprintf("%s#%s", tx.Name, periodKey)

	// sync.Map is complex. `LoadOrStore` gets existing value or stores a new one.
	val, _ := s.frequencyCounter.LoadOrStore(mapKey, action.NewSafeRecordList())
	txList := val.(*action.SafeRecordList)
	// Lock the specific list, append the transaction, then unlock
	txList.Add(tx)
}

// Done is called after all workers are done.
func (s *High) Done(accTxs *model.AccountRecords, result *action.AnalyseResult) {
	highThreshold := s.rule.HighThreshold
	// --- High Frequency Algorithm ---
	s.frequencyCounter.Range(func(key, value interface{}) bool {
		mapKey := key.(string)
		txList := value.(*action.SafeRecordList)
		count := len(txList.Records)
		if count > s.rule.HighThreshold {
			// This key (e.g., "John Smith#2024-03") is suspicious
			// In a real system, we'd now go back to the DB to fetch the actual txIDs
			// For this demo, we'll create a summary result.
			reason := fmt.Sprintf("高频交易: 对手方 %s 在周期内交易 %d 次, 超过阈值 %d", mapKey, count, highThreshold)
			result.AddBatch(reason, model.SuType_FreqHigh, txList.Records...)
		}
		return true // continue iteration
	})
}
