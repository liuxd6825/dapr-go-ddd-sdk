package times

import (
	"context"
	"fmt"

	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/service/suspicious/action"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/model"

	"time"
)

// SignificantDate
// 实现了“重大日期交易”的核心审计逻辑。
// 这是一个行级规则，但其判断逻辑依赖于复杂的日期计算。
type SignificantDate struct {
	rule model2.TimeSignificantDateRule
}

func NewSignificantDate(rule model2.TimeSignificantDateRule) *SignificantDate {
	return &SignificantDate{
		rule: rule,
	}
}

func (s *SignificantDate) IsEnable() bool {
	return s.rule.IsEnable
}

func (s *SignificantDate) DoAction(ctx context.Context, tx *model.Record, txIndex int, accTxs *model2.AccountRecords, result *action.AnalyseResult) {
	if !s.rule.IsEnable {
		return
	}

	// --- 步骤 1: 检查金额是否达到“重大”标准 ---
	amount := tx.Payout
	if tx.Income > amount {
		amount = tx.Income
	}

	if amount < s.rule.Amount {
		return
	}

	txDate := tx.Date

	// --- 步骤 2: 获取当前交易日期所在月份的最后一天 ---
	// time.Date(year, month+1, 0, ...) 是一个获取月末日期的巧妙技巧
	lastDayOfMonth := time.Date(txDate.Year(), txDate.Month()+1, 0, 0, 0, 0, 0, txDate.Location())

	// --- 步骤 3: 检查这个月末日是否是我们关心的“重大日期”类型 ---
	isTargetDateType := false
	month := lastDayOfMonth.Month()
	timeRule := s.rule
	if timeRule.CheckYear && month == time.December {
		isTargetDateType = true
	} else if timeRule.CheckQuarter && (month == time.March || month == time.June || month == time.September || month == time.December) {
		isTargetDateType = true
	} else if timeRule.CheckMonth {
		isTargetDateType = true
	}

	if !isTargetDateType {
		return
	}

	// --- 步骤 4: 计算交易日期与重大日期的天数差，判断是否在窗口内 ---
	daysDiff := txDate.Sub(lastDayOfMonth).Hours() / 24

	// Case A: 交易在重大日期之前 (daysDiff will be negative)
	if daysDiff <= 0 && daysDiff >= float64(-timeRule.DaysBefore) {
		reason := fmt.Sprintf("重大日期交易: 临近 %s (%.0f天前)", lastDayOfMonth.Format("2006-01-02"), -daysDiff)
		result.AddRecord(tx, model2.SuType_TimeConcentratedPayment, reason)
	}

	// Case B: 交易在重大日期之后 (daysDiff will be positive)
	if daysDiff > 0 && daysDiff <= float64(timeRule.DaysAfter) {
		reason := fmt.Sprintf("重大日期交易: 临近 %s (%.0f天后)", lastDayOfMonth.Format("2006-01-02"), daysDiff)
		result.AddRecord(tx, model2.SuType_TimeConcentratedPayment, reason)
	}
}
