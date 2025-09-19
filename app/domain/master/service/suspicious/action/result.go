package action

import (
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"sync"
)

type AnalyseResult struct {
	Account string
	Records map[string]*model2.SuRecord
	Batches []*AnalyseBatch
}

type AnalyseBatch struct {
	Name    string        `json:"name" gorm:"name" bson:"name"`
	Type    model2.SuType `json:"type" bson:"type" bson:"type"`
	Records []*model.Record
}

func NewAnalyseResult(taskId string, account string) *AnalyseResult {
	return &AnalyseResult{
		Account: account,
		Records: make(map[string]*model2.SuRecord),
		Batches: []*AnalyseBatch{},
	}
}

func (s *AnalyseResult) AddBatch(name string, suType model2.SuType, records ...*model.Record) *AnalyseBatch {
	batch := &AnalyseBatch{
		Name:    name,
		Type:    suType,
		Records: records,
	}
	for _, record := range records {
		s.AddRecord(record, suType, name)
	}
	s.Batches = append(s.Batches, batch)
	return batch
}

func (s *AnalyseResult) AddRecord(tx *model.Record, suType model2.SuType, reason string) *model2.SuRecord {
	var item *model2.SuRecord
	if val, ok := s.Records[tx.Id]; ok {
		item = val
	} else {
		item = &model2.SuRecord{
			Record: *tx,
		}
		s.Records[tx.Id] = item
	}
	switch suType {
	case model2.SuType_AmountLarge:
		if !item.AmountLarge {
			item.AmountLarge = true
			item.Risk++
		}

	case model2.SuType_AmountNumber:
		if !item.AmountNumber {
			item.AmountNumber = true
			item.Risk++
		}

	case model2.SuType_AmountCollar:
		if !item.AmountCollar {
			item.AmountCollar = true
			item.Risk++
		}

	case model2.SuType_AmountNear:
		if !item.AmountNear {
			item.AmountNear = true
			item.Risk++
		}

	case model2.SuType_TimeNonWorking:
		if !item.TimeNonWorking {
			item.TimeNonWorking = true
			item.Risk++
		}

	case model2.SuType_TimeConcentratedPayment:
		if !item.TimeConcentratedPayment {
			item.TimeConcentratedPayment = true
			item.Risk++
		}

	case model2.SuType_TimeFastInOut:
		if !item.TimeFastInOut {
			item.TimeFastInOut = true
			item.Risk++
		}

	case model2.SuType_TimeSignificantDate:
		if !item.TimeSignificantDate {
			item.TimeSignificantDate = true
			item.Risk++
		}

	case model2.SuType_FreqAbnormal:
		if !item.FreqAbnormal {
			item.FreqAbnormal = true
			item.Risk++
		}

	case model2.SuType_FreqHigh:
		if !item.FreqHigh {
			item.FreqHigh = true
			item.Risk++
		}

	case model2.SuType_FreqSleep:
		if !item.FreqSleep {
			item.FreqSleep = true
			item.Risk++
		}

	case model2.SuType_PartyBusiness:
		if !item.PartyBusiness {
			item.PartyBusiness = true
			item.Risk++
		}

	case model2.SuType_PartyAggregate:
		if !item.PartyAggregate {
			item.PartyAggregate = true
			item.Risk++
		}

	case model2.SuType_PartyPrivate:
		if !item.PartyPrivate {
			item.PartyPrivate = true
			item.Risk++
		}

	case model2.SuType_PartyHighRisk:
		if !item.PartyHighRisk {
			item.PartyHighRisk = true
			item.Risk++
		}

	case model2.SuType_PartyRelated:
		if !item.PartyRelated {
			item.PartyRelated = true
			item.Risk++
		}

	}
	item.Reasons = append(item.Reasons, model2.SuRecordReason{
		Reason: reason,
		Type:   suType,
	})
	return item
}

// SafeRecordList
// A thread-safe slice to store transactions
type SafeRecordList struct {
	mu          sync.Mutex
	Records     []*model.Record
	TotalAmount float64
}

func NewSafeRecordList() *SafeRecordList {
	return &SafeRecordList{Records: make([]*model.Record, 0)}
}

func (s *SafeRecordList) Add(record *model.Record) {
	s.mu.Lock()
	s.Records = append(s.Records, record)
	s.mu.Unlock()
}
