package service

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/suspicious"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/suspicious/action"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

type SuTaskService struct {
	tranDao        *dao.TranDao
	taskDao        *dao.SuTaskDao
	accountDao     *dao.SuTaskAccountDao
	suRecordDao    *dao.SuRecordDao
	suBatchDao     *dao.SuBatchDao
	suBatchItemDao *dao.SuBatchItemDao
}

func NewSuTaskService() *SuTaskService {
	return &SuTaskService{
		tranDao:        dao.NewTranDao(config.DBKey),
		taskDao:        dao.NewSuTaskDao(config.DBKey),
		suRecordDao:    dao.NewSuRecordDao(config.DBKey),
		accountDao:     dao.NewSuTaskAccountDao(config.DBKey),
		suBatchDao:     dao.NewSuBatchDao(config.DBKey),
		suBatchItemDao: dao.NewSuBatchItemDao(config.DBKey),
	}
}

func (s *SuTaskService) Analyse(ctx context.Context, taskId string) error {
	task, err := s.taskDao.FindById(ctx, taskId)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not found " + taskId)
	}
	accounts, err := s.accountDao.FindByTaskId(ctx, taskId)
	if err != nil {
		return err
	}
	analyse := suspicious.NewAnalyse(task, accounts, newDataRepo(s))
	analyse.DoAction(ctx)
	results := analyse.GetResults()

	records := make([]*model.SuRecord, 0)
	batches := make([]*model.SuBatch, 0)
	batchItems := make([]*model.SuBatchItem, 0)

	for _, result := range results {
		for _, record := range result.Records {
			records = append(records, record)
			record.RecordId = record.Id
			record.Id = getSuRecordId(task.Id, record.Id)
		}
		for _, item := range result.Batches {
			batch := &model.SuBatch{
				BaseModel: xbase.BaseModel{},
				Name:      item.Name,
				Type:      item.Type,
				Count:     len(item.Records),
				Payout:    0,
				Income:    0,
			}
			batches = append(batches, batch)
			for _, record := range item.Records {
				batchItems = append(batchItems, &model.SuBatchItem{
					BatchId:    batch.Id,
					SuRecordId: getSuRecordId(task.Id, record.Id),
					RecordId:   record.Id,
				})
			}
		}
	}

	s.suBatchDao.CreateMany(ctx, batches)
	s.suBatchItemDao.CreateMany(ctx, batchItems)
	s.suRecordDao.CreateMany(ctx, records)

	return nil
}

func getSuRecordId(taskId string, recordId string) string {
	return fmt.Sprintf("%s_%s", taskId, recordId)
}

type dataRepo struct {
	service *SuTaskService
}

func newDataRepo(service *SuTaskService) *dataRepo {
	return &dataRepo{
		service: service,
	}
}

func (d *dataRepo) GetCoreSuppliers() (map[string]action.Counterparty, error) {
	return map[string]action.Counterparty{}, nil
}

func (d *dataRepo) LoadCounterparties() (map[string]action.Counterparty, error) {
	return map[string]action.Counterparty{}, nil
}

func (d *dataRepo) GetAllTransactionsGroupedByAccount() (map[string][]*model.Tran, error) {
	return map[string][]*model.Tran{}, nil
}

func (d *dataRepo) LoadRelatedParties() ([]action.RelatedParty, error) {
	return []action.RelatedParty{}, nil
}

func (d *dataRepo) GetCounterparty(name string) (part *action.Counterparty, found bool) {
	return nil, false
}
