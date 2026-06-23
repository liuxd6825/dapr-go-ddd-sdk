package times

import (
	"context"
	"fmt"

	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/service/suspicious/action"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/model"
)

// ConcentratedPayments
// 实现了“集中支付”的核心审计逻辑。
// 它分析单个账户的交易流，以识别在风险时段内的密集支付行为。
type ConcentratedPayments struct {
	rule  model2.TimeConcentratedPaymentsRule
	txMap map[string]*model.Tran
}

func NewConcentratedPayments(rule model2.TimeConcentratedPaymentsRule) *ConcentratedPayments {
	return &ConcentratedPayments{
		rule: rule,
	}
}

func (s *ConcentratedPayments) IsEnable() bool {
	return s.rule.IsEnable
}

func (s *ConcentratedPayments) DoAction(ctx context.Context, tx *model.Record, txIndex int, accTxs *model2.AccountRecords, result *action.AnalyseResult) {
}
func (s *ConcentratedPayments) Done(accTxs *model2.AccountRecords, result *action.AnalyseResult) {
	txMap := make(map[string]*model.Record)

	// --- 步骤 1: 创建一个聚合桶 (Aggregation Bucket) ---
	// Key: YYYY-MM-DD 格式的日期字符串
	// Value: 一个切片，存储所有在该日风险时段内发生的支付交易的ID
	paymentsInRiskWindowByDay := make(map[string][]*model.Record)

	// --- 步骤 2: 遍历所有交易，将符合条件的交易放入桶中 ---
	for _, tx := range accTxs.Records {
		// 只关心支出交易
		if tx.Income == 0 {
			continue
		}

		txHour := tx.Date.Hour()

		// 检查交易时间是否落在用户定义的风险窗口内
		// 需要处理跨午夜的情况，例如 23:00 - 02:00
		inRiskWindow := false
		if s.rule.StartTime <= s.rule.EndTime {
			// e.g., 01:00 - 05:00
			if txHour >= s.rule.StartTime && txHour < s.rule.EndTime {
				inRiskWindow = true
			}
		} else {
			// e.g., 23:00 - 02:00 (crosses midnight)
			if txHour >= s.rule.StartTime || txHour < s.rule.EndTime {
				inRiskWindow = true
			}
		}

		if inRiskWindow {
			dateStr := tx.Date.Format("2006-01-02")
			txMap[tx.Id] = tx
			paymentsInRiskWindowByDay[dateStr] = append(paymentsInRiskWindowByDay[dateStr], tx)
		}
	}

	// --- 步骤 3: 检查每个桶，看是否超过了密度阈值 ---
	for dateStr, records := range paymentsInRiskWindowByDay {
		if len(records) > s.rule.Count {
			payout := 0.0
			income := 0.0
			for _, record := range records {
				payout += record.Payout
				income += record.Income
			}
			// 如果超过阈值，则这个桶里的所有交易都是可疑的
			reason := fmt.Sprintf("集中支付: 在 %s 的风险时段内发生 %d 笔支付, 超过阈值 %d, 收入总额 %.2f 支出总额 %.2f",
				dateStr, len(records), s.rule.Count, income, payout)
			result.AddBatch(reason, model2.SuType_TimeConcentratedPayment, records...)
		}
	}

}
