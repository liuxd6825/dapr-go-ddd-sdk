package party

import (
	"context"
	"fmt"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/service/suspicious/action"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
)

// Relation
// 关联方筛选：
// 将交易对手方与公司提供的股东、高管、关键员工及其近亲属名单，以及子公司、联营企业名单进行比对，筛选出所有关联方交易。
type Relation struct {
	rule            model2.PartRelationPartyRule
	repo            action.DataRepo
	nameMap         map[string]string // Key: Name, Value: RelationType
	accountMap      map[string]string // Key: AccountNumber, Value: RelationType
	counterpartyMap map[string]action.Counterparty
}

func NewRelationParty(rule model2.PartRelationPartyRule) *Relation {
	return &Relation{
		rule:       rule,
		nameMap:    map[string]string{},
		accountMap: map[string]string{},
	}
}

func (s *Relation) IsEnable() bool {
	return s.rule.IsEnable
}

func (s *Relation) DoAction(ctx context.Context, tx *model.Record, txIndex int, accTxs *model2.AccountRecords, result *action.AnalyseResult) {

}

// Update the prepare method
func (s *Relation) prepare() error {
	fmt.Println("\n[Orchestrator] Starting preparation phase...")

	// --- Step 1: Load Counterparties (as before) ---
	var err error
	s.counterpartyMap, err = s.repo.LoadCounterparties()
	if err != nil {
		return fmt.Errorf("failed to load counterparties: %w", err)
	}

	// --- Step 2: [NEW] Load and build lookup maps for Related Parties ---
	relatedParties, err := s.repo.LoadRelatedParties()
	if err != nil {
		return fmt.Errorf("failed to load related parties: %w", err)
	}

	s.nameMap = make(map[string]string)
	s.accountMap = make(map[string]string)
	for _, rp := range relatedParties {
		if rp.Name != "" {
			s.nameMap[rp.Name] = rp.RelationType
		}
		if rp.AccountNumber != "" {
			s.accountMap[rp.AccountNumber] = rp.RelationType
		}
	}
	fmt.Printf("[Orchestrator] Built lookup maps for %d related parties.\n", len(relatedParties))

	return nil
}

// isRelatedPartyTransaction implements the core matching logic.
// It's a high-performance row-level check using the pre-built maps.
func (s *Relation) isRelatedPartyTransaction(tx model.Tran) (bool, string) {
	// --- 步骤 1: 按对手方名称匹配 ---
	if relationType, found := s.nameMap[tx.Name]; found {
		reason := fmt.Sprintf("关联方交易(按名称): 对手方为'%s', 关系:%s", tx.Name, relationType)
		return true, reason
	}

	// --- 步骤 2: 按对手方账号匹配 ---
	// (Only check if name didn't match, and if account number is available)
	if tx.OppAcct != "" {
		if relationType, found := s.accountMap[tx.Acct]; found {
			reason := fmt.Sprintf("关联方交易(按账号): 对手方为'%s', 关系:%s", tx.Name, relationType)
			return true, reason
		}
	}

	// --- 步骤 3: 如果都没有匹配上 ---
	return false, ""
}

// isRelatedPartyTransaction implements the core matching logic.
// It's a high-performance row-level check using the pre-built maps.
func (s *Relation) isRelatedPartyTransaction2(tx model.Tran) (bool, string) {
	// --- 步骤 1: 按对手方名称匹配 ---
	if relationType, found := s.nameMap[tx.Name]; found {
		reason := fmt.Sprintf("关联方交易(按名称): 对手方为'%s', 关系:%s", tx.Name, relationType)
		return true, reason
	}

	// --- 步骤 2: 按对手方账号匹配 ---
	// (Only check if name didn't match, and if account number is available)
	if tx.Acct != "" {
		if relationType, found := s.accountMap[tx.Acct]; found {
			reason := fmt.Sprintf("关联方交易(按账号): 对手方为'%s', 关系:%s", tx.Name, relationType)
			return true, reason
		}
	}

	// --- 步骤 3: 如果都没有匹配上 ---
	return false, ""
}
