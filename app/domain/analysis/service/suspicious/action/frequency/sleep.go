package frequency

import (
	"context"
	"fmt"

	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/service/suspicious/action"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/model"

	"sync"
)

// Sleep
// 休眠账户交易
type Sleep struct {
	rule             model2.FreqSleepRule
	frequencyCounter *sync.Map // Key: string (e.g., "CounterpartyName#YYYY-MM"), Value: *atomic.Int64
}

func NewSleep(rule model2.FreqSleepRule) *Sleep {
	return &Sleep{
		rule:             rule,
		frequencyCounter: &sync.Map{},
	}
}

func (s *Sleep) IsEnable() bool {
	return s.rule.IsEnable
}

func (s *Sleep) DoAction(ctx context.Context, tx *model.Record, txIndex int, accTxs *model2.AccountRecords, result *action.AnalyseResult) {

}

// Done is called after all workers are done.
func (s *Sleep) Done(accTxs *model2.AccountRecords, result *action.AnalyseResult) {
	txs := accTxs.Records
	// 如果总交易数都不足以触发激活阈值，则直接返回
	if len(accTxs.Records) < s.rule.ActivateThreshold {
		return
	}

	// --- 步骤 1: 遍历交易，寻找“休眠期”的结束点 ---
	// 我们从第二笔交易开始，检查它与前一笔交易的时间间隔
	for i := 1; i < len(txs); i++ {
		prevTx := txs[i-1]
		currTx := txs[i] // `currTx` 是休眠期后的第一笔交易，即“激活点”

		// 计算两笔相邻交易的时间间隔
		hibernationDays := currTx.Date.Sub(prevTx.Date).Hours() / 24

		// 如果间隔大于用户定义的“休眠期”，我们就找到了一个潜在的激活事件
		if hibernationDays > float64(s.rule.SleepDays) {

			// --- 步骤 2: 验证“激活期”内的交易频率 ---
			// 定义激活期的截止时间
			activationWindowEnd := currTx.Date.AddDate(0, 0, s.rule.ActivateDays)

			var activationRecords []*model.Record
			// 从“激活点”开始，向后统计在激活窗口内的所有交易
			for j := i; j < len(txs); j++ {
				activationTx := txs[j]
				// 如果交易时间在窗口内，则加入列表
				if !activationTx.Date.After(activationWindowEnd) {
					activationRecords = append(activationRecords, activationTx)
				} else {
					// 因为数据是排序的，一旦超出窗口，就可以提前终止搜索
					break
				}
			}

			// --- 步骤 3: 检查交易频率是否达到阈值并记录结果 ---
			if len(activationRecords) >= s.rule.ActivateThreshold {
				reason := fmt.Sprintf("休眠账户激活: 在休眠%.0f天后, 于%d天内发生%d笔交易",
					hibernationDays, s.rule.ActivateDays, len(activationRecords))

				// 将激活期内的所有交易都标记为可疑
				result.AddBatch(reason, model2.SuType_FreqAbnormal, activationRecords...)
				// [性能优化]：一旦发现并记录了一段激活期，我们可以跳过这段已分析的交易，
				// 从激活期的最后一笔交易之后开始寻找下一个休眠期。
				// i = i + len(activationTransactions) - 1
			}
		}
	}
}
