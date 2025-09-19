package action

import (
	"context"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"time"
)

type Action interface {
	IsEnable() bool
	DoAction(ctx context.Context, tx *model.Record, txIndex int, accTxs *model2.AccountRecords, result *AnalyseResult)
}

type ActionDone interface {
	Done(accTxs *model2.AccountRecords, result *AnalyseResult)
}

type EntityType string

const (
	EntityType_Individual EntityType = "individual"
	EntityType_Company    EntityType = "company"
)

// Counterparty
// 交易对方
type Counterparty struct {
	Name              string     // 对手方名称 (主键)
	EstablishmentDate time.Time  // 成立日期 (用于“新成立即大额交易”规则)
	BusinessScope     string     // 经营范围 (用于“业务不匹配”规则)
	EntityType        EntityType // 实体类型, 例如 "Company" 或 "Individual" (用于“公转私”规则)
	IsRelatedParty    bool       // 是否为关联方 (一个简化的关联方标记)
	// NEW FIELDS for high-risk check
	Status          string // "正常", "经营异常", "吊销", "注销"
	LegalCasesCount int    // 涉及的司法案件数量
}

type RelatedParty struct {
	Name          string
	AccountNumber string
	RelationType  string // e.g., "高管", "股东", "子公司"
}

// DataRepo interface
type DataRepo interface {
	LoadCounterparties() (map[string]Counterparty, error)
	GetAllTransactionsGroupedByAccount() (map[string][]*model.Tran, error)
	// LoadRelatedParties 查找关联方
	LoadRelatedParties() ([]RelatedParty, error)
	// GetCounterparty 查找对手方
	GetCounterparty(name string) (part *Counterparty, found bool)
	// GetCoreSuppliers 取得核心供应商名单
	GetCoreSuppliers() (map[string]Counterparty, error)
}
