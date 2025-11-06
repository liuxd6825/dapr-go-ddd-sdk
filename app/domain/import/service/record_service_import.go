package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/enum"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/event"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	task_pkg "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/pkg/readexcel"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/config"
	xbase2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types/times"
	gp2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
)

type ctxKey struct {
}

type TaskOptions struct {
	TenantId string
	CaseId   string
	DocId    string
	FileId   string
	TaskId   string
}

type Create4ExcelResult struct {
	RecordTotal int64 `json:"recordTotal"`
	ErrorCount  int64 `json:"errorCount"`
}

// Preview Excel文件内容预览
func (s *RecordService) Preview(ctx context.Context, task *task_pkg.Task, temp *task_pkg.RecordTemplate) (records []*task_pkg.RecordIe, res *Create4ExcelResult, err error) {
	fileByte, err := s.docFileService.ReadByteByFileId(ctx, task.FileId)
	if err != nil {
		return nil, nil, err
	}

	buffer := bytes.NewBuffer(fileByte)
	gp2.Try(func() error {
		// 读取缓存数据
		res, err = s.readExcel(ctx, task, temp, buffer, true, func(ctx context.Context, list []*task_pkg.RecordIe, batch readexcel.Batching) error {
			records = list
			return nil
		})
		return err
	}).Catch(func(e error) {
		err = e
	}).Finally(func() {

	})
	return records, res, err
}

// readExcel
// @Description: 读取excel文件
// @receiver r
// @param ctx
// @param fileName
// @param sheetName
// @param recordTemp
// @param batchFunc
// @param opts
// @return error
func (s *RecordService) readExcel(ctx context.Context, task *task_pkg.Task, temp *task_pkg.RecordTemplate,
	buffer *bytes.Buffer, isPreview bool, batchFunc func(ctx context.Context, list []*task_pkg.RecordIe, paging readexcel.Batching) error,
) (*Create4ExcelResult, error) {
	if task == nil {
		return nil, errors.New("parameter 'task' is null")
	}
	if temp == nil {
		return nil, errors.New("parameter 'task.Template' is null")
	}

	tmp, err := temp.NewTemplate()
	if err != nil {
		return nil, err
	}

	batchFun := func(ctx context.Context, list []*task_pkg.RecordIe, paging readexcel.Batching) error {
		return batchFunc(ctx, list, paging)
	}

	newCtx := NewContext(ctx, &TaskOptions{
		TenantId: task.TenantId,
		CaseId:   task.CaseId,
		DocId:    task.DocId,
		FileId:   task.FileId,
		TaskId:   task.Id,
	})
	table, err := readexcel.ReadByteToEntity[*task_pkg.RecordIe](newCtx, buffer, task.SheetName, tmp, isPreview, newRecord, batchFun, &readexcel.Options{BatchSize: 1000})

	if err != nil {
		return nil, err
	}

	res := &Create4ExcelResult{
		RecordTotal: int64(len(table.Rows)),
		ErrorCount:  int64(len(table.Errors)),
	}
	return res, err
}

// Create4Excel
// @Description: 从Excel文件批量创建流水记录
// @receiver r
// @param ctx
// @param cmd
// @return error
func (s *RecordService) Create4Excel(ctx context.Context, cmd *command.RecordCreate4ExcelCommand, batchBack func(batch readexcel.Batching) error) (res *Create4ExcelResult, err error) {
	now := time.Now()
	task := &task_pkg.Task{
		DocId:     cmd.Data.DocId,
		FileId:    cmd.Data.FileId,
		FileName:  cmd.Data.FileName,
		SheetName: cmd.Data.SheetName,
		Total:     0,
		Complete:  0,
		StartTime: &now,
		State:     task_pkg.TaskStateEditing,
	}
	taskId := cmd.Data.TaskId
	task.Id = taskId
	task.CaseId = cmd.Data.CaseId

	fileByte, err := s.docFileService.ReadByteByFileId(ctx, cmd.Data.FileId)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(fileByte)

	gp2.Try(func() error {
		// 读取缓存数据
		res, err = s.readExcel(ctx, task, cmd.Data.Template, buffer, false, func(ctx context.Context, list []*task_pkg.RecordIe, batch readexcel.Batching) error {
			fields := logs.Fields{
				"call":       "readExcel()",
				"taskId":     task.Id,
				"batchIndex": batch.BatchIndex,
				"total":      batch.RowTotal,
				"batchSize":  batch.BatchSize,
				"pageCount":  len(list),
			}
			logs.Debug(ctx, fields)
			createRes := s.dao.CreateMany(ctx, list)
			if createRes != nil {
				return err
			}
			if batchBack != nil {
				return batchBack(batch)
			}
			//  time.Sleep(20 * time.Second) // 用于Actor超时测试使用
			return nil
		})
		if err == nil {
			cmd1 := command.TaskUpdateStateCommand{}
			cmd1.CommandId = cmd.CommandId
			cmd1.Data = field.TaskUpdateStateFields{}
			cmd1.Data.Id = cmd.Data.TaskId
			cmd1.Data.State = task_pkg.TaskStateGenerated
			cmd1.Data.Message = fmt.Sprintf("成功%v条，错误%v条", res.RecordTotal, res.ErrorCount)
			err = s.taskService.UpdateState(ctx, &cmd1)
		}
		return err
	}).Catch(func(e error) {
		err = e
	}).Finally(func() {

	})
	return res, err
}

// Import2Master
// @Description: 将任务中的银行流水导入取主数据中。
// @receiver r
// @param ctx
// @param appcmd
// @return error
func (s *RecordService) Import2Master(ctx context.Context, appcmd *command.RecordImport2MasterCommand) (err error) {
	taskId := appcmd.Data.TaskId
	/*
		tenantId := appctx.GetTenantId2(ctx)
		if count, err := s.CountErrorByTaskId(ctx, tenantId, taskId); err != nil {
			return err
		} else if count > 0 {
			return errors.New("发现错误数据%v条，请更正后再提交。", count)
		}
	*/

	pageNum := int64(0)
	recordCount := int64(0)
	pageSize := gp2.IfElse[int64](appcmd.Data.PageSize == 0, 2000, appcmd.Data.PageSize)

	// 批量导入数据
	createMany := func(ctx context.Context, appcmd *command.RecordImport2MasterCommand) error {
		logs.Debugf(ctx, nil, "createMany() context=%v", func() any {
			return appctx.GetMessage(ctx)
		})
		return tx.StartTx(ctx, []string{config.DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
			for {
				// 取得批量数据
				qry := query.NewRecordIeFindPagingByTaskIdQuery(appcmd.Data.TaskId, false)
				qry.SetPageNum(pageNum)
				qry.SetPageSize(pageSize)

				res := s.FindPagingByTaskId(ctx, qry)
				if err != nil {
					return err
				}
				if res.GetDataLength() == 0 {
					return nil
				}
				// 如果没有数据
				if res.GetPageNum() > pageNum {
					return nil
				}

				// 通过领域事件，导入流水到主数据
				// 生成领域事件
				create, err := newRecordCreateManyFromExcelCommand(appcmd, res.GetData())
				if err != nil {
					return err
				}
				if pageNum > 0 {
					create.Data.IsAddItems = true
				}
				// 发布领域事件
				if err := s.PublishImportRecordToMasterEvent(ctx, create); err != nil {
					return err
				}
				recordCount += int64(len(res.GetData()))
				pageNum++
			}

		})
	}

	// 执行导入
	gp2.Try(func() error {
		startTime := times.PNow()
		// 批量导入数据
		err = createMany(ctx, appcmd)
		if err == nil {
			// 更新任务状态
			cmd := command.NewTaskUpdateProgressCommand(appcmd.CommandId, taskId)
			cmd.Data.Complete = recordCount
			cmd.Data.StartTime = startTime
			cmd.Data.State = task_pkg.TaskStateImported
			cmd.Data.EndTime = times.PNow()
			err = s.taskService.UpdateProgress(ctx, &cmd.Data)
		}
		return err
	}).Catch(func(e error) {
		err = e
	}).Finally(func() {

	})
	return err
}

func (s *RecordService) PublishImportRecordToMasterEvent(ctx context.Context, event *event.RecordImportMasterEvent) (err error) {
	return xbase2.PublishEvent(ctx, config.ImportAppId, event, nil)
}

func NewContext(ctx context.Context, task *TaskOptions) context.Context {
	return context.WithValue(ctx, ctxKey{}, task)
}

func getTaskByContext(ctx context.Context) (*TaskOptions, bool) {
	val := ctx.Value(ctxKey{})
	if task, ok := val.(*TaskOptions); ok {
		return task, true
	}
	return nil, false
}

func newRecord(ctx context.Context, row *readexcel.DataRow, temp *readexcel.Template) (*task_pkg.RecordIe, error) {
	task, ok := getTaskByContext(ctx)
	if !ok {
		return nil, errors.New("task is nil")
	}
	record := &task_pkg.RecordIe{
		BaseModel: xbase2.BaseModel{
			Id:       uuid.NewString(),
			TenantId: task.TenantId,
			CaseId:   task.CaseId,
		},
		DocId:  task.DocId,
		FileId: task.FileId,
		TaskId: task.TaskId,
		RowNum: row.RowNum,
	}
	record.Name, _ = temp.ValueToString(task_pkg.FieldName_Name.String(), row)
	record.Acct, _ = temp.ValueToString(task_pkg.FieldName_Account.String(), row)
	record.AcctType, _ = temp.ValueToString(task_pkg.FieldName_AccountType.String(), row)
	record.Category, _ = temp.ValueToString(task_pkg.FieldName_Category.String(), row)
	record.BankName, _ = temp.ValueToString(task_pkg.FieldName_BankName.String(), row)

	record.OppIden, _ = temp.ValueToString(task_pkg.FieldName_OppIden.String(), row)
	record.OppName, _ = temp.ValueToString(task_pkg.FieldName_OppName.String(), row)
	record.OppAcct, _ = temp.ValueToString(task_pkg.FieldName_OppAccount.String(), row)
	record.OppAcctType, _ = temp.ValueToString(task_pkg.FieldName_OppAccountType.String(), row)
	record.OppBankName, _ = temp.ValueToString(task_pkg.FieldName_OppBankName.String(), row)
	record.OppCategory, _ = temp.ValueToString(task_pkg.FieldName_OppCategory.String(), row)

	record.Type, _ = temp.ValueToString(task_pkg.FieldName_Type.String(), row)
	if ccy, _ := temp.ValueToString(task_pkg.FieldName_Ccy.String(), row); len(ccy) == 0 {
		record.Ccy = "CNY"
	}
	record.Summary, _ = temp.ValueToString(task_pkg.FieldName_Summary.String(), row)
	record.Notes, _ = temp.ValueToString(task_pkg.FieldName_Notes.String(), row)
	record.Serial, _ = temp.ValueToString(task_pkg.FieldName_Serial.String(), row)

	row.GetString(temp, task_pkg.FieldName_Place.String(), func(v string) {
		record.Place = v
	}, func(err ...error) {
		record.AddErrors(err...)
	})

	row.GetString(temp, task_pkg.FieldName_Iden.String(), func(v string) {
		record.Iden = v
	}, func(err ...error) {
		record.AddErrors(err...)
	})

	row.GetFloat(temp, task_pkg.FieldName_Payout.String(), func(v *float64) {
		record.Payout = v
	}, func(err ...error) {
		record.AddErrors(err...)
	})

	row.GetFloat(temp, task_pkg.FieldName_Income.String(), func(v *float64) {
		record.Income = v
	}, func(err ...error) {
		record.AddErrors(err...)
	})

	row.GetFloat(temp, task_pkg.FieldName_Amount.String(), func(v *float64) {
		record.Amount = v
	}, func(err ...error) {
		record.AddErrors(err...)
	})

	row.GetFloat(temp, task_pkg.FieldName_Balance.String(), func(v *float64) {
		record.Balance = v
	}, func(err ...error) {
		record.AddErrors(err...)
	})

	row.GetDate(temp, task_pkg.FieldName_Date.String(), func(v *time.Time) {
		record.Date = v
	}, func(err ...error) {
		record.AddErrors(err...)
	})

	record.Errors = row.Errors
	record.Cells = newCells(row.Cells)

	return record, nil

}

func newCells(mapDataCells map[string]readexcel.DataCells) map[string]task_pkg.RecordIeCells {
	if mapDataCells == nil {
		return nil
	}
	res := map[string]task_pkg.RecordIeCells{}
	for fk, cells := range mapDataCells {
		cs := task_pkg.RecordIeCells{}
		for _, c := range cells {
			cs = append(cs, task_pkg.RecordIeCell{
				Key:    c.Key,
				Row:    c.Row,
				Col:    c.Col,
				Value:  c.Value,
				Errors: c.Errors,
			})
		}
		res[fk] = cs
	}
	return res
}

// newCreateManyCommand
// @Description:
// @param appcmd
// @param list
// @return *command.RecordCreateManyFromExcelCommand
// @return error
func newRecordCreateManyFromExcelCommand(appcmd *command.RecordImport2MasterCommand, list []*task_pkg.RecordIe) (*event.RecordImportMasterEvent, error) {
	items := make([]*field.RecordFields, 0)
	for _, e := range list {
		record := &field.RecordFields{
			RowNum:   e.RowNum,
			Id:       e.Id,
			Iden:     e.Iden,
			Name:     e.Name,
			Acct:     e.Acct,
			AcctType: e.AcctType,
			Category: e.Category,
			Balance:  e.Balance,
			BankName: e.BankName,

			OppIden:     e.OppIden,
			OppName:     e.OppName,
			OppAcct:     e.OppAcct,
			OppAcctType: e.OppAcctType,
			OppCategory: e.OppCategory,
			OppBankName: e.OppBankName,

			Serial:  e.Serial,
			Payout:  e.Payout,
			Income:  e.Income,
			Amount:  e.Amount,
			Date:    e.Date,
			Type:    e.Type,
			Ccy:     e.Ccy,
			Place:   e.Place,
			Summary: e.Summary,
			Notes:   e.Notes,
		}
		if len(record.Ccy) == 0 {
			record.Ccy = enum.CurrencyCNY.String()
		}
		/*
			if _, err := json.Marshal(record); err != nil {
				fmt.Println(fmt.Sprintf("id:%s", record.Id))
				return nil, err
			}
		*/
		items = append(items, record)
	}

	cmd := &event.RecordImportMasterEvent{}
	cmd.EventId = idutils.NewId()
	cmd.OccurredOn = time.Now()
	cmd.Data = event.RecordImportMasterEventData{
		CaseId:    appcmd.Data.CaseId,
		DocId:     appcmd.Data.DocId,
		FileName:  appcmd.Data.FileName,
		FileId:    appcmd.Data.FileId,
		SheetId:   appcmd.Data.SheetId,
		SheetName: appcmd.Data.SheetName,
		TaskId:    appcmd.Data.TaskId,
		Items:     items,
	}

	return cmd, nil
}
