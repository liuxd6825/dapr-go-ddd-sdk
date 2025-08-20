package suspicious

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/suspicious/action"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/suspicious/action/amount"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/suspicious/action/frequency"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/suspicious/action/party"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/suspicious/action/times"
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

type Analyse struct {
	task      *model.SuTask
	rule      *model.SuTaskRule
	recordDao *dao.RecordDao
	accounts  []*model.SuTaskAccount
	actions   []action.Action
	results   []*action.AnalyseResult
}

func NewAnalyse(task *model.SuTask, accounts []*model.SuTaskAccount, repo action.DataRepo) *Analyse {
	rules := &task.Rules
	actions := []action.Action{
		amount.NewCollar(rules.AmountCollar),
		amount.NewAmountLarge(rules.AmountLarge),
		amount.NewAmountNumber(rules.AmountNumber),
		amount.NewAmountNear(rules.AmountNear),

		times.NewFastInOutHours(rules.TimeFastInOut),
		times.NewNonWorkingHours(rules.TimeNonWorkingHours),
		times.NewConcentratedPayments(rules.TimeConcentratedPayments),
		times.NewSignificantDate(rules.TimeSignificantDate),

		frequency.NewHigh(rules.FreqHigh),
		frequency.NewAbnormal(rules.FreqAbnormal),
		frequency.NewSleep(rules.FreqSleep),

		party.NewAggregate(rules.PartAggregate),
		party.NewHighRisk(rules.PartHighRisk, repo),
		party.NewPrivate(rules.PartPrivate, repo),
		party.NewRelationParty(rules.PartRelation),
	}
	analyse := &Analyse{
		task:      task,
		rule:      rules,
		accounts:  accounts,
		recordDao: dao.NewRecordDao(config.DBKey),
		actions:   actions,
	}
	return analyse
}

// =============================================================================
// 4. ORCHESTRATOR (Grouping Producer Model)
// =============================================================================

func (s *Analyse) DoAction(ctx context.Context) (int, time.Duration) {
	startTime := time.Now()

	// --- 1. Pre-load dimension data ---

	// --- 2. [CORE CHANGE] Group all transactions by Account first ---
	txsByAccount, err := s.getAccountTrans(ctx)
	if err != nil {
		logs.Errorfmt(ctx, "Failed to load transactions: %v", err)
	}

	fmt.Println("\n[Orchestrator] Starting data grouping phase...")

	// --- 3. Setup concurrent pipeline ---
	numWorkers := runtime.NumCPU()
	jobs := make(chan model.AccountRecords, numWorkers)
	results := make(chan action.AnalyseResult, 100)

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
				taskResult := s.applyRulesToAccount(ctx, &accountTxs)
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

	for result := range results {
		s.results = append(s.results, &result)
	}
	return suspiciousCount, time.Since(startTime)
}

func (s *Analyse) GetResults() []*action.AnalyseResult {
	return s.results
}

// applyRulesToAccount 执行分析规则
// It can access `sta.params` and `sta.counterpartyMap` directly.
func (s *Analyse) applyRulesToAccount(ctx context.Context, accTxs *model.AccountRecords) (result *action.AnalyseResult) {
	txs := accTxs.Records
	result = action.NewAnalyseResult(s.task.Id, accTxs.Account)
	for i, tx := range txs {
		for _, a := range s.actions {
			if a.IsEnable() {
				a.DoAction(ctx, tx, i, accTxs, result)
			}
		}
	}
	for _, a := range s.actions {
		if a.IsEnable() {
			if done, ok := a.(action.ActionDone); ok {
				done.Done(accTxs, result)
			}
		}
	}
	return result
}

func (s *Analyse) getAccountTrans(ctx context.Context) ([]*model.AccountRecords, error) {
	res := []*model.AccountRecords{}
	for _, account := range s.accounts {
		txs, err := s.recordDao.FindByAccountOppAccount(ctx, account.Account, s.task.StartTime, s.task.EndTime)
		if err != nil {
			return nil, err
		}
		res = append(res, &model.AccountRecords{
			Account: account.Account,
			Records: txs,
		})
	}
	return res, nil
}
