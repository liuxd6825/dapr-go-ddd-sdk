package party

import (
	"context"
	"fmt"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/suspicious/action"
	"sort"
	"sync"
)

// Aggregate
// 资金集中度分析： 分析资金支付是否高度集中于少数几个非核心供应商。
type Aggregate struct {
	rule             model2.PartAggregateRule
	coreSuppliersMap map[string]bool
	frequencyCounter *sync.Map
	// NEW: Concurrent map for concentration analysis
	paymentAggregator *sync.Map // Key: counterpartyName, Value: *atomic.Float64
}

func NewAggregate(rule model2.PartAggregateRule) *Aggregate {
	return &Aggregate{
		rule:              rule,
		paymentAggregator: &sync.Map{},
		frequencyCounter:  &sync.Map{},
	}
}

func (s *Aggregate) IsEnable() bool {
	return s.rule.IsEnable
}

func (s *Aggregate) DoAction(ctx context.Context, tx *model.Record, txIndex int, accTxs *model2.AccountRecords, result *action.AnalyseResult) {
	if tx.Payout > 0 {
		if _, isCore := s.coreSuppliersMap[tx.OppName]; isCore {
			return
		}
		val, _ := s.paymentAggregator.LoadOrStore(tx.OppName, action.NewSafeRecordList())
		xtList := val.(*action.SafeRecordList)
		xtList.Add(tx)
		xtList.TotalAmount += tx.Amount
	}
}

func (s *Aggregate) Done(accTxs *model2.AccountRecords, result *action.AnalyseResult) {
	// --- High Frequency Check (as before) ---
	// ...

	// --- [NEW] Fund Concentration Analysis ---
	fmt.Println("[Aggregator] Analyzing fund concentration...")

	// Step 1: Convert concurrent map to a slice for sorting
	type paymentStat struct {
		CounterpartyName string
		RecordList       *action.SafeRecordList
	}
	var paymentStats []paymentStat
	var grandTotalPaid float64

	s.paymentAggregator.Range(func(key, value interface{}) bool {
		name := key.(string)
		xtList := value.(*action.SafeRecordList)
		paymentStats = append(paymentStats, paymentStat{CounterpartyName: name, RecordList: xtList})
		grandTotalPaid += xtList.TotalAmount
		return true
	})

	// Step 2: Sort the slice by total paid in descending order
	sort.Slice(paymentStats, func(i, j int) bool {
		return paymentStats[i].RecordList.TotalAmount > paymentStats[j].RecordList.TotalAmount
	})

	// Step 3: Calculate the total paid to the Top N suppliers
	var topNTotalPaid float64
	topN := s.rule.TopN
	if topN > len(paymentStats) {
		topN = len(paymentStats)
	}

	for i := 0; i < topN; i++ {
		topNTotalPaid += paymentStats[i].RecordList.TotalAmount
	}

	// Step 4: Calculate the concentration percentage and check against the threshold
	if grandTotalPaid > 0 { // Avoid division by zero
		concentrationRatio := (topNTotalPaid / grandTotalPaid) * 100
		if concentrationRatio > s.rule.Percent {
			for i := 0; i < topN; i++ {
				item := paymentStats[i]
				reason := fmt.Sprintf("资金高度集中: Top %d 位非核心供应商 %s, 金额：%.0f", i+1, item.CounterpartyName, item.RecordList.TotalAmount)
				result.AddBatch(reason, model2.SuType_PartyAggregate, item.RecordList.Records...)
			}
		}
	}

}
