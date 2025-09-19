package service

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/query"
	dao2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/dao"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/suspicious"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/suspicious/action"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

type SuTaskService struct {
	tranDao        *dao2.TranDao
	taskDao        *dao.SuTaskDao
	accountDao     *dao.SuAccountDao
	suRecordDao    *dao.SuRecordDao
	suBatchDao     *dao.SuBatchDao
	suBatchItemDao *dao.SuBatchItemDao
}

func NewSuTaskService() *SuTaskService {
	return &SuTaskService{
		tranDao:        dao2.NewTranDao(config.DBKey),
		taskDao:        dao.NewSuTaskDao(config.DBKey),
		suRecordDao:    dao.NewSuRecordDao(config.DBKey),
		accountDao:     dao.NewSuAccountDao(config.DBKey),
		suBatchDao:     dao.NewSuBatchDao(config.DBKey),
		suBatchItemDao: dao.NewSuBatchItemDao(config.DBKey),
	}
}

// Create
//
//	@Description:
//	@receiver s
//	@param ctx
//	@param v
//	@param opts
//	@return error
func (s *SuTaskService) Create(ctx context.Context, v *model.SuTask, opts ...idao.CallOptions) error {
	return s.taskDao.Create(ctx, v, opts...).GetError()
}

func (s *SuTaskService) CreateMany(ctx context.Context, v []*model.SuTask, opts ...idao.CallOptions) error {
	return s.taskDao.CreateMany(ctx, v, opts...).GetError()
}

func (s *SuTaskService) Update(ctx context.Context, v *model.SuTask, opts ...idao.CallOptions) error {
	return s.taskDao.Update(ctx, v, opts...).GetError()
}

func (s *SuTaskService) UpdateMany(ctx context.Context, v []*model.SuTask, opts ...idao.CallOptions) error {
	return s.taskDao.UpdateMany(ctx, v, opts...).GetError()
}

func (s *SuTaskService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.taskDao.DeleteById(ctx, id, opts...).GetError()
}

func (s *SuTaskService) FindById(ctx context.Context, qry *query.SuTaskFindByIdQuery, opts ...idao.CallOptions) (*model.SuTask, error) {
	return s.taskDao.FindById(ctx, qry.Id, opts...)
}

func (s *SuTaskService) FindBillById(ctx context.Context, qry *query.SuTaskFindByIdQuery, opts ...idao.CallOptions) (*model.SuTaskBillView, error) {
	task, err := s.FindById(ctx, qry, opts...)
	if err != nil {
		return nil, err
	}
	accounts, err := s.accountDao.FindByTaskId(ctx, qry.Id)
	if err != nil {
		return nil, err
	}
	taskView := &model.SuTaskBillView{
		Task:     task,
		Accounts: accounts,
	}
	return taskView, nil
}

func (s *SuTaskService) FindPaging(ctx context.Context, qry *ddd_query.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.SuTask] {
	return s.taskDao.FindPaging(ctx, qry, opts...)
}

func (s *SuTaskService) Start(ctx context.Context, taskId string) error {
	task, err := s.taskDao.FindById(ctx, taskId)
	if err != nil {
		return err
	}
	if task.Status != model.SuTaskStatus_New {
		err := errors.NewVerifyError()
		err.AppendField("status", "当前任务状态不是“新建”")
		return err.GetError()
	}
	taskAccounts, err := s.accountDao.FindByTaskId(ctx, task.Id)
	if err != nil {
		return err
	}
	if len(taskAccounts) == 0 {
		return errors.New("taskAccounts is length 0 ")
	}

	if err := s.UpdateStatus(ctx, taskId, model.SuTaskStatus_PendingInfo); err != nil {
		return err
	}
	if err = s.analyse(ctx, task, taskAccounts); err != nil {
		if err := s.UpdateStatus(ctx, taskId, model.SuTaskStatus_New); err != nil {
			return err
		}
		return err
	}
	if err = s.UpdateStatus(ctx, taskId, model.SuTaskStatus_InProgress); err != nil {
		return err
	}
	return nil
}

func (s *SuTaskService) UpdateStatus(ctx context.Context, taskId string, status model.SuTaskStatus) error {
	return s.taskDao.UpdateStatus(ctx, taskId, status)
}

func (s *SuTaskService) UpdateInfo(ctx context.Context, taskId string, info dao.SuTaskUpdateInfo) error {
	return s.taskDao.UpdateInfo(ctx, taskId, info)
}

func (s *SuTaskService) analyse(ctx context.Context, task *model.SuTask, taskAccounts []*model.SuTaskAccount) error {
	if task == nil {
		return errors.New("task is null ")
	}
	if len(taskAccounts) == 0 {
		return errors.New("taskAccounts is length 0 ")
	}

	analyse := suspicious.NewAnalyse(task, taskAccounts, newDataRepo(s))
	results, useTime, err := analyse.DoAction(ctx)

	logs.Infofmt(ctx, "用时：%s", useTime.String())

	if err != nil {
		return err
	}

	for _, result := range results {
		records := make([]*model.SuRecord, 0)
		batches := make([]*model.SuBatch, 0)
		batchItems := make([]*model.SuBatchItem, 0)

		for _, record := range result.Records {
			records = append(records, record)
			record.RecordId = record.Id
			record.TaskId = task.Id
			record.Id = getSuRecordId(task.Id, record.Id)
		}
		for _, item := range result.Batches {
			batch := &model.SuBatch{
				BaseModel: xbase.BaseModel{},
				TaskId:    task.Id,
				Name:      item.Name,
				Type:      item.Type,
				Count:     len(item.Records),
				Payout:    0,
				Income:    0,
			}
			batches = append(batches, batch)
			for _, record := range item.Records {
				batchItems = append(batchItems, &model.SuBatchItem{
					TaskId:     task.Id,
					BatchId:    batch.Id,
					SuRecordId: getSuRecordId(task.Id, record.Id),
					RecordId:   record.Id,
				})
			}
		}
		s.suBatchDao.CreateMany(ctx, batches)
		s.suBatchItemDao.CreateMany(ctx, batchItems)
		s.suRecordDao.CreateMany(ctx, records)
	}
	return nil
}

func (s *SuTaskService) ClearAnalyseResults(ctx context.Context, taskId string) error {
	builder := rsql.NewBuilder().And(rsql.Eq("task_id", taskId))
	where := builder.Build()
	if err := s.suRecordDao.DeleteByRSQL(ctx, where).GetError(); err != nil {
		return err
	}
	if err := s.suBatchDao.DeleteByRSQL(ctx, where).GetError(); err != nil {
		return err
	}
	if err := s.suBatchItemDao.DeleteByRSQL(ctx, where).GetError(); err != nil {
		return err
	}
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

func (d *dataRepo) GetAllTransactionsGroupedByAccount() (map[string][]*model2.Tran, error) {
	return map[string][]*model2.Tran{}, nil
}

func (d *dataRepo) LoadRelatedParties() ([]action.RelatedParty, error) {
	return []action.RelatedParty{}, nil
}

func (d *dataRepo) GetCounterparty(name string) (part *action.Counterparty, found bool) {
	return nil, false
}
