package service

import (
	"context"
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/service/suspicious"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/service/suspicious/action"
	dao2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/dao"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/code/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ddd/ddd_query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/tasks"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
	"go.temporal.io/sdk/client"
)

type SuTaskService struct {
	tranDao        *dao2.TranDao
	taskDao        *dao.SuTaskDao
	accountDao     *dao.SuAccountDao
	suRecordDao    *dao.SuRecordDao
	suBatchDao     *dao.SuBatchDao
	suBatchItemDao *dao.SuBatchItemDao
	taskLogDao     *dao.SuTaskLogDao
	codeService    *service.CodeService
}

func NewSuTaskService() *SuTaskService {
	return &SuTaskService{
		tranDao:        dao2.NewTranDao(config.DBKey),
		taskDao:        dao.NewSuTaskDao(config.DBKey),
		suRecordDao:    dao.NewSuRecordDao(config.DBKey),
		accountDao:     dao.NewSuAccountDao(config.DBKey),
		suBatchDao:     dao.NewSuBatchDao(config.DBKey),
		suBatchItemDao: dao.NewSuBatchItemDao(config.DBKey),
		taskLogDao:     dao.NewSuTaskLogDao(config.DBKey),
		codeService:    service.NewCodeService(),
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
func (s *SuTaskService) Create(ctx context.Context, cmd *command.SuTaskCreateCommand) (*model.SuTask, error) {
	task := cmd.NewTask()
	if task.Code == "" {
		task.Code = s.codeService.NewSuCode(ctx, cmd.Data.CaseId)
	}
	task.StatusName = task.Status.String()
	err := s.taskDao.Create(ctx, task).GetError()
	return task, err
}

func (s *SuTaskService) CreateMany(ctx context.Context, v []*model.SuTask, opts ...idao.CallOptions) error {
	for _, task := range v {
		task.StatusName = task.Status.String()
	}
	return s.taskDao.CreateMany(ctx, v, opts...).GetError()
}

func (s *SuTaskService) Update(ctx context.Context, v *model.SuTask, opts ...idao.CallOptions) error {
	v.StatusName = v.Status.String()
	return s.taskDao.Update(ctx, v, opts...).GetError()
}

func (s *SuTaskService) UpdateMany(ctx context.Context, v []*model.SuTask, opts ...idao.CallOptions) error {
	for _, task := range v {
		task.StatusName = task.Status.String()
	}
	return s.taskDao.UpdateMany(ctx, v, opts...).GetError()
}

func (s *SuTaskService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.taskDao.DeleteById(ctx, id, opts...).GetError()
}

func (s *SuTaskService) QueryById(ctx context.Context, qry *query.SuTaskFindByIdQuery, opts ...idao.CallOptions) (*model.SuTask, error) {
	return s.taskDao.FindById(ctx, qry.Id, opts...)
}

func (s *SuTaskService) FindById(ctx context.Context, id string, opts ...idao.CallOptions) (*model.SuTask, error) {
	return s.taskDao.FindById(ctx, id, opts...)
}

func (s *SuTaskService) FindByCaseId(ctx context.Context, caseId string, opts ...idao.CallOptions) ([]*model.SuTask, error) {
	return s.taskDao.FindByRSQL(ctx, fmt.Sprintf("case_id=='%s' and status>=3", caseId), opts...)
}

func (s *SuTaskService) QueryBillById(ctx context.Context, qry *query.SuTaskFindByIdQuery, opts ...idao.CallOptions) (*model.SuTaskBillView, error) {
	task, err := s.QueryById(ctx, qry, opts...)
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

func (s *SuTaskService) QueryPaging(ctx context.Context, qry *ddd_query.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.SuTask] {
	return s.taskDao.FindPaging(ctx, qry, opts...)
}

func (s *SuTaskService) Analysis(ctx context.Context, taskId string) error {
	verifyError := errors.NewVerifyError()
	task, err := s.taskDao.FindById(ctx, taskId)
	if err != nil {
		return err
	}

	if task.Status != model.SuTaskStatus_New {
		verifyError.AppendField("status", fmt.Sprintf("当前任务状态不是“%s”", model.SuTaskStatus_New))
	}
	/*
		if task.WorkflowId != "" {
			verifyError.AppendField("workflowId", "当前任务已经在执行中")
		}
	*/

	accountsCount, err := s.accountDao.CountByTaskId(ctx, task.Id)
	if err != nil {
		return err
	}
	if accountsCount == 0 {
		verifyError.AppendField("account", "分析“账号”信息不能为空")
	}

	if verifyError.HasError() {
		return verifyError
	}

	if err := s.UpdateStatus(ctx, taskId, model.SuTaskStatus_TaskQueuing); err != nil {
		return err
	}

	// 配置工作流选项
	workflowId := idutils.NewId()
	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowId,           // 工作流的唯一ID
		TaskQueue: tasks.GetTaskQueue(), // 必须与 Worker 监听的任务队列名称一致
	}
	ctxMap := appctx.NewMapWithContext(ctx)
	_, err = tasks.ExecuteWorkflow(ctx, workflowOptions, TaskAnalysisWorkflow, taskId, workflowId, ctxMap)
	if err != nil {
		return err
	}
	return err
}

// analysisTask
// @Description: 执行数据分析，在Temporal工作流中执行。
// @receiver s
// @param ctx
// @param taskId
// @param workflowId
// @return err
func (s *SuTaskService) analysisTask(ctx context.Context, taskId string, workflowId string) (err error) {
	verr := errors.NewVerifyError()
	task, err := s.taskDao.FindById(ctx, taskId)
	if err != nil {
		return err
	}
	if task.Status != model.SuTaskStatus_TaskQueuing {
		verr.AppendField("status", fmt.Sprintf("当前任务状态为“%s”，应为“新建”", task.Status))
	}
	taskAccounts, err := s.accountDao.FindByTaskId(ctx, task.Id)
	if err != nil {
		return err
	}
	if len(taskAccounts) == 0 {
		verr.AppendField("account", "分析“账号”信息不能为空")
	}

	// 当执行错误是，不可以返回error， 否则任务会重复执行。
	gp.Try(func() error {
		if verr.HasError() {
			return verr.GetError()
		}
		if err := s.UpdateWorkflowIdStatus(ctx, taskId, workflowId, model.SuTaskStatus_InProgress); err != nil {
			return err
		}
		recordCount, err := s.analyse(ctx, task, taskAccounts)
		if err != nil {
			return err
		}
		err = s.updateInspect(ctx, task, recordCount)
		return err
	}).Catch(func(e error) {
		index := int64(1)
		msg := fmt.Sprintf("分析可疑任务是出错, error:%s", e.Error())
		s.taskLogDao.Create(ctx, model.NewTaskLog(taskId, index, msg))
		if e1 := s.UpdateWorkflowIdStatus(ctx, taskId, "", model.SuTaskStatus_New); e1 != nil {
			msg = fmt.Sprintf("分析可疑任务是出错后, 恢复状态时出错:%s", e.Error())
			s.taskLogDao.Create(ctx, model.NewTaskLog(taskId, index+1, msg))
		}
	})
	return nil
}

// updateInspect
// @Description: 更新统计后的任务信息
// @receiver s
// @param ctx
// @param task
// @param recordCount
// @return err
func (s *SuTaskService) updateInspect(ctx context.Context, task *model.SuTask, recordCount int64) (err error) {
	suRecordCount, err := s.suRecordDao.CountByTaskId(ctx, task.Id)
	if err != nil {
		return err
	}
	suHighCount, err := s.suRecordDao.CountByTaskIdHighRisk(ctx, task.Id)
	if err != nil {
		return err
	}
	suRecordSum, err := s.suRecordDao.SumByTaskId(ctx, task.Id)
	if err != nil {
		return err
	}
	status := model.SuTaskStatus_Inspect
	fields := dao.SuTaskFields{
		RecordCount: &recordCount,
		SuCount:     &suRecordCount,
		TotalAmount: &suRecordSum,
		SuHighCount: &suHighCount,
		Status:      &status,
	}

	err = s.UpdateFields(ctx, task.Id, fields)
	return err
}

func (s *SuTaskService) UpdateWorkflowIdStatus(ctx context.Context, taskId string, workflowId string, status model.SuTaskStatus) error {
	return s.taskDao.UpdateWorkflowIdStatus(ctx, taskId, workflowId, status)
}

func (s *SuTaskService) UpdateStatus(ctx context.Context, taskId string, status model.SuTaskStatus) error {
	return s.taskDao.UpdateStatus(ctx, taskId, status)
}

func (s *SuTaskService) UpdateFields(ctx context.Context, taskId string, info dao.SuTaskFields) error {
	return s.taskDao.UpdateFields(ctx, taskId, info)
}

func (s *SuTaskService) analyse(ctx context.Context, task *model.SuTask, taskAccounts []*model.SuTaskAccount) (recordCount int64, err error) {
	if task == nil {
		return 0, errors.New("task is null ")
	}
	if len(taskAccounts) == 0 {
		return 0, errors.New("taskAccounts is length 0 ")
	}

	analyse := suspicious.NewAnalyse(task, taskAccounts, newDataRepo(s))
	results, useTime, recordCount, err := analyse.DoAction(ctx)

	logs.Infofmt(ctx, "用时：%s", useTime.String())

	if err != nil {
		return 0, err
	}
	for _, result := range results {
		records := make([]*model.SuRecord, 0)
		batches := make([]*model.SuBatch, 0)
		batchItems := make([]*model.SuBatchItem, 0)
		for _, record := range result.SuRecords {
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
	return recordCount, nil
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
