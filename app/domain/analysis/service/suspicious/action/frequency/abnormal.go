package frequency

import (
	"context"
	"fmt"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/service/suspicious/action"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"sort"
)

/*
### **1. 算法业务逻辑和审计目标**

*   **业务定义 (Business Definition)**:
    “异常规律性支付”指的是，针对某一个**非核心业务对手方**（特别是个人），在一段时间内，发生了一系列满足以下条件的支付：
    *   **周期性 (Periodicity)**: 支付发生的时间间隔大致相等（例如，每隔30天左右）。
    *   **金额稳定性 (Amount Stability)**: 每次支付的金额大致相等。
    *   **持续性 (Continuity)**: 这种规律性的支付连续发生了多次（例如，至少3次以上）。

*   **审计目标 (Audit Objective)**:
    正常的业务采购或费用支付，其时间和金额通常是根据业务需求波动的。高度规律性的支付行为非常罕见，它强烈暗示了这笔钱的性质可能不是常规的业务往来，而是：
    1.  **账外借款的利息**: 公司可能有一笔未入账的私人借款，正在定期支付利息。
    2.  **回扣或佣金**: 向某个关键人物或中介公司定期支付的好处费。
    3.  **虚构的顾问费/服务费**: 与某个人或“壳公司”签订虚假合同，定期套取公司资金。
    4.  **分期支付的非法款项**: 将一笔大的非法支出，拆分成多次定期的支付来掩人耳目。
*/

// Abnormal
// 频率：异常规律性支付 (固定周期、固定金额支付给个人或非供应商)
// 识别向同一个非员工、非供应商的个人或单位，进行固定周期（如每周、每十天）、固定金额的支付。这可能是未入账的借款利息或回扣
type Abnormal struct {
	rule    model2.FreqAbnormalRule
	txsByCp map[string][]*model.Record
}

func NewAbnormal(rule model2.FreqAbnormalRule) *Abnormal {
	return &Abnormal{
		rule:    rule,
		txsByCp: make(map[string][]*model.Record),
	}
}

func (s *Abnormal) IsEnable() bool {
	return s.rule.IsEnable
}

func (s *Abnormal) DoAction(ctx context.Context, tx *model.Record, txIndex int, accTxs *model2.AccountRecords, result *action.AnalyseResult) {
	// Exclude core partners
	if tx.Name == "BigCorp Inc." {
		return
	}
	//s.txsByCp[tx.Name] = append(s.txsByCp[tx.Name], tx)
}

// NEW HELPER: groupTxsByCounterparty
func (s *Abnormal) groupTxsByCounterparty(allTxs []*model.Record) map[string][]*model.Record {
	txsByCp := make(map[string][]*model.Record)
	for _, tx := range allTxs {
		// Exclude core partners
		if tx.Name == "BigCorp Inc." {
			continue
		}
		txsByCp[tx.Name] = append(txsByCp[tx.Name], tx)
	}
	// Sort each group by time
	for _, txs := range txsByCp {
		sort.Slice(txs, func(i, j int) bool {
			return txs[i].Date.Before(txs[j].Date)
		})
	}
	return txsByCp
}

// Done “异常规律性支付”的核心算法
// 输入的txs是属于同一个对手方的、按时间排序的所有支出交易
func (s *Abnormal) Done(accTxs *model2.AccountRecords, result *action.AnalyseResult) {
	txsByCp := s.groupTxsByCounterparty(accTxs.Records)
	for _, tx := range txsByCp {
		s.findAbnormalRegularity(tx, result)
	}

}
func (s *Abnormal) findAbnormalRegularity(counterpartyTxs []*model.Record, result *action.AnalyseResult) {
	// We only care about payments (debits)
	var payments []*model.Record
	for _, tx := range counterpartyTxs {
		if tx.Payout != 0 {
			payments = append(payments, tx)
		}
	}

	if len(payments) < s.rule.MinOccurrences {
		return
	}

	periodDays := float64(s.rule.PeriodDays)
	dayTolerance := float64(s.rule.DayTolerance)
	// 这是一个状态机，用于跟踪连续的规律性支付
	var currentSequence []*model.Record

	for i := 1; i < len(payments); i++ {
		prevTx := payments[i-1]
		currTx := payments[i]

		// --- 步骤1: 检查周期是否在容差范围内 ---
		daysDiff := currTx.Date.Sub(prevTx.Date).Hours() / 24
		isPeriodMatch := daysDiff >= (periodDays-dayTolerance) && daysDiff <= (periodDays+dayTolerance)

		// --- 步骤2: 检查金额是否在容差范围内 ---
		amountDiffRatio := (currTx.Amount - prevTx.Amount) / prevTx.Amount
		amountToleranceRatio := s.rule.AmountTolerance / 100.0
		isAmountMatch := amountDiffRatio <= amountToleranceRatio

		// --- 步骤3: 状态转移 ---
		if isPeriodMatch && isAmountMatch {
			// 如果匹配，将交易加入当前序列
			if len(currentSequence) == 0 {
				currentSequence = append(currentSequence, prevTx)
			}
			currentSequence = append(currentSequence, currTx)
		} else {
			// 如果不匹配，检查之前的序列是否已满足最小次数要求
			if len(currentSequence) >= s.rule.MinOccurrences {
				// 序列中断，且长度达标，记录结果
				reason := fmt.Sprintf("异常规律性支付: 发现%d笔连续支付, 周期约%d天", len(currentSequence), s.rule.PeriodDays)
				result.AddBatch(reason, model2.SuType_FreqAbnormal, currentSequence...)
			}
			// 重置状态机
			currentSequence = nil
		}
	}

	// --- 步骤4: 循环结束后，检查最后一个序列 ---
	if len(currentSequence) >= s.rule.MinOccurrences {
		reason := fmt.Sprintf("异常规律性支付: 发现%d笔连续支付, 周期约%d天", len(currentSequence), s.rule.PeriodDays)
		result.AddBatch(reason, model2.SuType_FreqAbnormal, currentSequence...)
	}
}
