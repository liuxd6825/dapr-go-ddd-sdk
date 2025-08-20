package times

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/suspicious/action"
	"time"
)

type FastInOutHours struct {
	rule model.TimeFastInOutRule
}

func NewFastInOutHours(rule model.TimeFastInOutRule) *FastInOutHours {
	return &FastInOutHours{
		rule: rule,
	}
}

func (s *FastInOutHours) IsEnable() bool {
	return s.rule.IsEnable
}

// DoAction 快进快出
func (s *FastInOutHours) DoAction(ctx context.Context, tx *model.Record, txIndex int, accTxs *model.AccountRecords, result *action.AnalyseResult) {

	// --- 步骤 1: 寻找“进” (Credit) ---
	inTx := tx
	if inTx.Name == inTx.OppName {
		return
	}
	if inTx.Income < s.rule.TxAmount {
		// 如果当前交易不是收入，则跳过，继续寻找下一个可能的“进”
		return
	}
	// 这里我们使用一个简单的匹配规则
	amountMatchThreshold := inTx.Income * (s.rule.Percent / 100)

	// --- 步骤 2: 在时间窗口内向后寻找匹配的“出” (Debit) ---
	// 定义时间窗口的截止时间
	windowEndTime := inTx.Date.Add(time.Duration(s.rule.FastInOutHours * float64(time.Hour)))
	// 使用滑动窗口的方式，遍历所有可能的“一进一出”交易对
	for i := txIndex + 1; i < len(accTxs.Records)-1; i++ {
		outTx := accTxs.Records[i]

		// 检查2a: 如果搜索的交易时间已经超出了窗口，那么后续的交易也必然超时
		// 因为数据是排序的，这是一个重要的性能优化，可以提前终止内层循环。
		if outTx.Date.After(windowEndTime) {
			break
		}

		// 检查2b: 确保这是一笔支出交易
		// 检查2c: (可选，但推荐) 检查金额是否匹配或相近
		if outTx.GetPayout() >= amountMatchThreshold {
			// --- 步骤 3: 发现模式，记录结果 ---
			duration := outTx.Date.Sub(inTx.Date)
			inDate := inTx.Date.Format(time.DateTime)
			outDate := outTx.Date.Format(time.DateTime)
			reason := fmt.Sprintf("快进快出(进): %s收入%.2f后%s支出%.2f，时隔:%.2f小时", inDate, inTx.Income, outDate, outTx.GetPayout(), duration.Hours())
			result.AddBatch(reason, model.SuType_TimeFastInOut, inTx, outTx)
			// 找到了一个匹配的“出”，通常就可以结束本次对“进”的搜索
			// 因为一笔钱通常只会被转走一次。
			break
		}

	}

}
