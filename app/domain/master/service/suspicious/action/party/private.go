package party

import (
	"context"
	"fmt"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/suspicious/action"
)

// Private
// 个人账户筛选： 筛选出所有与个人账户之间的大额交易（公转私或私转公）。
type Private struct {
	rule model2.PartPrivateRule
	repo action.DataRepo
}

func NewPrivate(rule model2.PartPrivateRule, repo action.DataRepo) *Private {
	return &Private{
		rule: rule,
		repo: repo,
	}
}

func (s *Private) IsEnable() bool {
	return s.rule.IsEnable
}

func (s *Private) DoAction(ctx context.Context, tx *model.Record, txIndex int, accTxs *model2.AccountRecords, result *action.AnalyseResult) {
	isSuspicious, reason := s.isLargePrivateTransaction(tx)
	if isSuspicious {
		result.AddRecord(tx, model2.SuType_PartyPrivate, reason)
	}
}

// isLargePrivateTransaction 实现了“个人账户大额交易筛选”的核心逻辑。
// 它依赖于预加载的 counterpartyMap 来识别对手方的实体类型。
func (s *Private) isLargePrivateTransaction(tx *model.Record) (bool, string) {
	// --- 步骤 1: 从预加载的map中获取对手方信息 ---
	counterpartyInfo, found := s.repo.GetCounterparty(tx.Name)

	// 如果在我们的主数据中找不到该对手方，或者它不是个人，则直接返回
	if !found || counterpartyInfo.EntityType != "Individual" {
		return false, ""
	}

	// --- 步骤 2: 检查交易金额是否超过对私交易的阈值 ---
	// COALESCE a.k.a get the non-zero value between debit and credit
	amount := tx.Amount

	if amount >= s.rule.TxThreshold {
		// --- 步骤 3: 确定交易方向并生成报警信息 ---
		direction := "公转私"
		if tx.Income > 0 {
			direction = "私转公"
		}

		reason := fmt.Sprintf("大额对私交易(%s): 金额 %.2f, 对手方: %s", direction, amount, tx.Name)

		return true, reason
	}

	return false, ""
}
