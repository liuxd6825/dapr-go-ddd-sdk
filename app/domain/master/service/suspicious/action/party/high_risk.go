package party

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/suspicious/action"
)

// HighRisk
// 高风险对手方筛查：
// 将对手方公司名称通过工商信息查询工具（如天眼查、企查查）进行核查，关注其成立时间（是否为新成立即有大额交易）、经营范围（是否与交易内容匹配）、股东背景以及是否存在经营异常或司法风险。
// 对手方名称模糊或与主营业务严重不符（如科技公司向一家农产品合作社支付大额“咨询费”）。
type HighRisk struct {
	rule model.PartHighRiskRule
	repo action.DataRepo
}

func NewHighRisk(rule model.PartHighRiskRule, dataRepo action.DataRepo) *HighRisk {
	return &HighRisk{
		rule: rule,
		repo: dataRepo,
	}
}

func (s *HighRisk) IsEnable() bool {
	return s.rule.IsEnable
}

func (s *HighRisk) DoAction(ctx context.Context, tx *model.Record, txIndex int, accTxs *model.AccountRecords, result *action.AnalyseResult) {

}

// findHighRiskCounterpartyTransactions 实现了“高风险对手方”的核心筛查逻辑。
// 它是一个复合的行级规则，会检查多个风险维度。
func (s *HighRisk) findHighRiskCounterpartyTransactions(tx *model.Record, result *action.AnalyseResult) {

	// --- 步骤 1: 获取对手方丰富的维度信息 ---
	info, found := s.repo.GetCounterparty(tx.OppName)
	if !found || info.EntityType != "Company" {
		// This rule only applies to companies
		return
	}

	amount := tx.Amount
	// --- 步骤 2: 逐一检查各个风险维度 ---

	// **风险维度 2a: 新成立即大额交易**
	daysSinceEstablishment := tx.Date.Sub(info.EstablishmentDate).Hours() / 24
	if daysSinceEstablishment >= 0 && daysSinceEstablishment < float64(s.rule.NewCompanyDaysThreshold) {
		if amount >= s.rule.NewCompanyLargeAmount {
			reason := fmt.Sprintf("高风险-新成立: 对手方成立仅 %.0f 天即发生 %.2f 元大额交易", daysSinceEstablishment, amount)
			result.AddRecord(tx, model.SuType_PartyHighRisk, reason)
		}
	}

	// **风险维度 2b: 工商状态异常**
	if info.Status != "正常" {
		if amount >= s.rule.AbnormalStatusTxAmount {
			reason := fmt.Sprintf("高风险-状态异常: 与状态为'%s'的对手方发生 %.2f 元交易", info.Status, amount)
			result.AddRecord(tx, model.SuType_PartyHighRisk, reason)
		}
	}

	// **风险维度 2c: 司法风险高**
	if info.LegalCasesCount > s.rule.HighLegalCasesThreshold {
		if amount >= s.rule.LegalCasesTxAmount {
			reason := fmt.Sprintf("高风险-司法: 对手方涉案%d起, 仍发生 %.2f 元交易", info.LegalCasesCount, amount)
			result.AddRecord(tx, model.SuType_PartyHighRisk, reason)
		}
	}

	//  业务不匹配 (Business Mismatch)**
	if amount > s.rule.BusinessMismatchTxAmount {

	}

}
