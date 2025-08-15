package action

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"runtime"
	"sync"
	"time"
)

// Counterparty represents dimension data about a transaction counterparty.
type Counterparty struct {
	Name              string
	EstablishmentDate time.Time
	BusinessScope     string
	EntityType        string // "Company" or "Individual"
	IsRelatedParty    bool
}

type SuTaskAnalyse struct {
	task     *model.SuTask
	rule     *model.SuTaskRule
	tranDao  *dao.TranDao
	accounts []*model.SuTaskAccount
}

// AccountTransactions is the unit of work for our workers.
type AccountTransactions struct {
	Account string
	Trans   []*model.Tran
}

const (
	TotalTransactions = 1_000_000 // Simulate processing 1 million transactions
	BatchSize         = 1_000     // Process 1000 transactions at a time
)

func NewSuTaskAnalyse(task *model.SuTask, accounts []*model.SuTaskAccount) *SuTaskAnalyse {
	return &SuTaskAnalyse{
		task:     task,
		rule:     &task.SuTaskRule,
		accounts: accounts,
		tranDao:  dao.NewTranDao(config.DBKey),
	}
}

func (s *SuTaskAnalyse) LoadCounterparties() (map[string]Counterparty, error) {
	fmt.Println("[DataRepo] Loading counterparties into memory map...")
	counterpartyMap := make(map[string]Counterparty)
	// In a real app, this would be a SELECT * FROM counterparties
	counterpartyMap["BigCorp Inc."] = Counterparty{Name: "BigCorp Inc.", EstablishmentDate: time.Now().AddDate(-5, 0, 0), EntityType: "Company", BusinessScope: "批发和零售业"}
	counterpartyMap["Consulting LLC"] = Counterparty{Name: "Consulting LLC", EstablishmentDate: time.Now().AddDate(-1, 0, 0), EntityType: "Company", BusinessScope: "技术服务,咨询服务"}
	counterpartyMap["John Smith"] = Counterparty{Name: "John Smith", EntityType: "Individual", IsRelatedParty: true}
	counterpartyMap["Fresh Start Co."] = Counterparty{Name: "Fresh Start Co.", EstablishmentDate: time.Now().AddDate(0, -2, 0), EntityType: "Company", BusinessScope: "广告业"}
	fmt.Printf("[DataRepo] Loaded %d counterparties.\n", len(counterpartyMap))
	return counterpartyMap, nil
}

// =============================================================================
// 4. ORCHESTRATOR (Grouping Producer Model)
// =============================================================================

func (s *SuTaskAnalyse) Analyze2(ctx context.Context) (int, time.Duration) {
	startTime := time.Now()

	// --- 1. Pre-load dimension data ---
	counterpartyMap, err := s.LoadCounterparties()
	if err != nil {
		logs.Errorfmt(ctx, "Failed to load counterparties: %v", err)
	}

	// --- 2. [CORE CHANGE] Group all transactions by Account first ---
	txsByAccount, err := s.getAccountTrans(ctx)
	if err != nil {
		logs.Errorfmt(ctx, "Failed to load transactions: %v", err)
	}

	fmt.Println("\n[Orchestrator] Starting data grouping phase...")

	// --- 3. Setup concurrent pipeline ---
	numWorkers := runtime.NumCPU()
	jobs := make(chan AccountTransactions, numWorkers)
	results := make(chan model.SuTaskResult, 1000)

	var workersWg sync.WaitGroup
	var aggregatorWg sync.WaitGroup

	// --- 4. Start the Aggregator ---
	aggregatorWg.Add(1)
	suspiciousCount := 0
	go func() {
		defer aggregatorWg.Done()
		for res := range results {
			suspiciousCount++
			if suspiciousCount <= 10 { // Print first few findings
				fmt.Printf("  -> Suspicious Finding: Account %s", res.Account)
			}
		}
	}()

	// --- 5. Start the Worker Pool ---
	fmt.Printf("[Orchestrator] Starting %d workers to process accounts...\n", numWorkers)
	for i := 0; i < numWorkers; i++ {
		workersWg.Add(1)
		go func() {
			defer workersWg.Done()
			for accountTxs := range jobs {
				taskResult := s.applyRulesToAccount(ctx, &accountTxs, counterpartyMap)
				results <- *taskResult
			}
		}()
	}

	// --- 6. [CORE CHANGE] Start the Grouping Producer ---
	// This goroutine now feeds whole account histories to the workers.
	go func() {
		defer close(jobs)
		for _, accTrans := range txsByAccount {
			jobs <- *accTrans
		}
		fmt.Println("[Producer] All accounts have been sent to workers.")
	}()

	// --- 7. Wait for shutdown ---
	workersWg.Wait()
	fmt.Println("[Orchestrator] All workers have finished.")
	close(results)
	aggregatorWg.Wait()
	fmt.Println("[Aggregator] All results have been processed.")

	return suspiciousCount, time.Since(startTime)
}

// applyRulesToAccount is a method of SuTaskAnalyse.
// It can access `sta.params` and `sta.counterpartyMap` directly.
func (s *SuTaskAnalyse) applyRulesToAccount(ctx context.Context, accTxs *AccountTransactions, counterpartyMap map[string]Counterparty) (result *model.SuTaskResult) {
	txs := accTxs.Trans
	result = model.NewSuTaskResult(s.task.Id, accTxs.Account)
	for i, tx := range txs {
		s.doAmountLarge(ctx, tx, result)
		s.doTimeIsNonWorkingHours(ctx, tx, result)
		s.doAmountCollar(ctx, tx, i, accTxs, result)
	}
	return result
}

func (s *SuTaskAnalyse) doAmountCollar(ctx context.Context, tx *model.Tran, txIndex int, accTxs *AccountTransactions, result *model.SuTaskResult) {
	if s.rule.IsAmountCollar {
		duration := time.Duration(s.rule.AmountRule.CollarDays) * 24 * time.Hour
		for i := txIndex + 1; i < len(accTxs.Trans); i++ {
			nextTx := accTxs.Trans[i]
			if nextTx.Date.Sub(tx.Date) > duration {
				return
			}
			if tx.Amount == nextTx.Amount && tx.Name == nextTx.OppName && tx.OppName == nextTx.Name && tx.OppName != tx.Name {
				item := result.AddItem(tx, model.SuType_AmountCollar, "金额折叠")
				item.AddAmountCollar(nextTx)
			}
		}
	}
	/*	if info, ok := counterpartyMap[tx.CounterpartyName]; ok && info.IsRelatedParty {
		suspiciousResults = append(suspiciousResults, SuspiciousResult{tx.ID, accTxs.AccountID, "关联方交易", tx.TransactionTime})
	}*/
}

// amountLarge Large Value Transaction
func (s *SuTaskAnalyse) doAmountLarge(ctx context.Context, tx *model.Tran, result *model.SuTaskResult) {
	if s.rule.IsAmountLarge {
		if tx.Amount > s.rule.AmountRule.LargeValue {
			result.AddItem(tx, model.SuType_AmountLarge, fmt.Sprintf("大额交易: 金额 %.2f", tx.Amount))
		}
	}
}

// timeIsNonWorkingHour Non-working hours
func (s *SuTaskAnalyse) doTimeIsNonWorkingHours(ctx context.Context, tx *model.Tran, result *model.SuTaskResult) {
	if s.rule.IsTimeNonWorkingHours {
		t := tx.Date
		if s.rule.TimeRule.IsNonWorkingSunday && t.Weekday() == time.Sunday {
			if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
				result.AddItem(tx, model.SuType_TimeNonWorkingHours, fmt.Sprintf("非工作日交易"))
				return
			}
		}
		if s.rule.TimeRule.IsNonWorkingHours {
			if t.Hour() < s.rule.TimeRule.NonWorkingHoursMin || t.Hour() > s.rule.TimeRule.NonWorkingHoursMax {
				result.AddItem(tx, model.SuType_TimeNonWorkingHours, fmt.Sprintf("非工作时间交易"))
				return
			}
		}
	}
}

func (s *SuTaskAnalyse) getAccountTrans(ctx context.Context) ([]*AccountTransactions, error) {
	res := []*AccountTransactions{}
	for _, account := range s.accounts {
		txs, err := s.tranDao.FindByAccountOppAccount(ctx, account.Account, s.task.StartTime, s.task.EndTime)
		if err != nil {
			return nil, err
		}
		res = append(res, &AccountTransactions{
			Account: account.Account,
			Trans:   txs,
		})
	}
	return res, nil
}
