package service

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	command2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/service/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/service/doris"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/event"
	doc "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/miniofs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
	"github.com/spf13/afero"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	task_pkg "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/pkg/readexcel"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/service/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types/times"
	gp2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
)

type ctxKey struct {
}

type TaskOptions struct {
	TenantId   string
	CaseId     string
	DocId      string
	FileId     string
	TaskId     string
	MasterId   string
	MasterType string
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

	runOpts, err := s.getRuntimeOptions(ctx)
	if err != nil {
		return nil, err
	}

	batchFun := func(ctx context.Context, list []*task_pkg.RecordIe, paging readexcel.Batching) error {
		return batchFunc(ctx, list, paging)
	}

	newCtx := NewContext(ctx, &TaskOptions{
		TenantId:   task.TenantId,
		CaseId:     task.CaseId,
		DocId:      task.DocId,
		FileId:     task.FileId,
		TaskId:     task.Id,
		MasterId:   task.MasterId,
		MasterType: task.MasterType,
	})
	table, err := readexcel.ReadByteToEntity[*task_pkg.RecordIe](newCtx, buffer, task.SheetName, tmp, isPreview, runOpts, newRecord, batchFun, &readexcel.Options{BatchSize: 1000})

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
func (s *RecordService) Create4Excel(ctx context.Context, cmd *command2.RecordCreate4ExcelCommand, batchBack func(batch readexcel.Batching) error) (res *Create4ExcelResult, err error) {
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
			cmd1 := command2.TaskUpdateStateCommand{}
			cmd1.CommandId = cmd.CommandId
			cmd1.Data = field.TaskUpdateStateFields{}
			cmd1.Data.Id = cmd.Data.TaskId
			cmd1.Data.State = task_pkg.TaskStateGenerated
			cmd1.Data.Message = fmt.Sprintf("成功%v条，错误%v条", res.RecordTotal, res.ErrorCount)
			err = s.taskService.UpdateState(ctx, &cmd1)
			if err != nil {
				return err
			}
			docModel := &doc.Document{}
			docModel.Id = idutils.NewId()
			docModel.FileId = task.FileId
			docModel.SourceId = task.DocId
			docModel.SourceType = "流水"
			docModel.CaseId = task.CaseId
			docModel.SourceApp = "document_service"
			err = s.docStatusProvider.UpdateStatus(ctx, docModel, task_pkg.TaskStateGenerated.Name(), "")
			if err != nil {
				return err
			}
		}
		return err
	}).Catch(func(e error) {
		err = e
	}).Finally(func() {

	})
	return res, err
}

func (s *RecordService) newDorisImporter() (*doris.DorisImporter, error) {
	dorisAny, ok := env.GetEnv().App.Meta["doris"]
	if !ok {
		return nil, errors.New("doris env not set")
	}

	value, ok := dorisAny.(map[string]any)
	if !ok {
		return nil, errors.New("doris env not set")
	}

	dorisConfig, err := doris.NewConfigWithMap(value)
	if err != nil {
		return nil, err
	}
	dorisConfig.Table = "master_record1"
	loader := doris.NewDorisImporter(*dorisConfig)
	return loader, nil
}

// Import2Master
// @Description: 将任务中的银行流水导入取主数据中。
// @receiver r
// @param ctx
// @param appcmd
// @return error
func (s *RecordService) Import2Master(ctx context.Context, appcmd *command2.RecordImport2MasterCommand) (err error) {

	tempFs, err := miniofs.NewFs(env.GetEnv(), "default", "importData")
	if err != nil {
		return err
	}

	loader, err := s.newDorisImporter()
	if err != nil {
		return err
	}
	err = tempFs.Mkdir("/record", os.ModeDir)
	if err != nil {
		return err
	}
	parquetFileName := fmt.Sprintf("/record/record_%s.parquet", appcmd.CommandId)
	parquetFile, err := doris.NewParquetFile[*doris.Record](tempFs, parquetFileName)
	if err != nil {
		return err
	}

	recordCount := int64(0)
	recordCount, err = s.writeParquetFile(ctx, parquetFile, appcmd)
	if err != nil {
		return err
	}

	return s.importData(ctx, tempFs, loader, appcmd, tempFs.GetBucketName(), parquetFileName, recordCount)

}

func (s *RecordService) importData(ctx context.Context, tempFs afero.Fs, loader *doris.DorisImporter, appcmd *command2.RecordImport2MasterCommand, s3BucketName, pFileName string, recordCount int64) error {
	taskId := appcmd.Data.TaskId
	// 导入数据到Doris数据库中
	_, err := loader.ImportFile(tempFs, pFileName, doris.LoadOptions{
		Format: doris.FormatParquet,
		Label:  "record_" + taskId,
	})
	if err != nil {
		return err
	}

	// 创建导入事件
	importEvent, err := newRecordCreateManyFromExcelCommand(ctx, appcmd, s3BucketName, pFileName)
	if err != nil {
		return err
	}

	startTime := times.PNow()
	// 更新任务状态
	cmd := command2.NewTaskUpdateProgressCommand(appcmd.CommandId, appcmd.Data.TaskId)
	cmd.Data.Complete = recordCount
	cmd.Data.StartTime = startTime
	cmd.Data.State = task_pkg.TaskStateImported
	cmd.Data.EndTime = times.PNow()

	docModel := &doc.Document{}
	docModel.Id = idutils.NewId()
	docModel.FileId = appcmd.Data.FileId
	docModel.SourceId = appcmd.Data.DocId
	docModel.SourceType = "流水"
	docModel.CaseId = appcmd.Data.CaseId
	docModel.SourceApp = "document_service"

	// 开启事务
	err = tx.StartTx(ctx, []string{config.DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
		// 发布领域事件
		if err := s.PublishImportRecordToMasterEvent(ctx, importEvent); err != nil {
			return err
		}
		if err := s.taskService.UpdateProgress(ctx, &cmd.Data); err != nil {
			return err
		}
		if err := s.docStatusProvider.UpdateStatus(ctx, docModel, task_pkg.TaskStateImported.Name(), ""); err != nil {
			return err
		}
		return nil
	})
	return err
}

// 生成parquet文件
func (s *RecordService) writeParquetFile(ctx context.Context, pFile *doris.ParquetFile[*doris.Record], appcmd *command2.RecordImport2MasterCommand) (recordCount int64, err error) {
	pageNum := int64(0)
	pageSize := gp2.IfElse[int64](appcmd.Data.PageSize == 0, 2000, appcmd.Data.PageSize)
	recordCount = int64(0)

	err = tx.StartTx(ctx, []string{config.DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
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

			for _, item := range res.GetData() {
				item.TaskId = appcmd.Data.TaskId
				if err1 := pFile.Write(doris.NewRecord(item)); err1 != nil {
					return err1
				}
			}
			recordCount += int64(len(res.GetData()))
			pageNum++
		}
	})

	if err != nil {
		return 0, err
	}
	if err = pFile.WriteStop(); err != nil {
		return 0, err
	}
	return recordCount, nil
}

func (s *RecordService) PublishImportRecordToMasterEvent(ctx context.Context, event *event.RecordImportMasterEvent) (err error) {
	return xbase.PublishEvent(ctx, event)
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
		BaseModel: xbase.BaseModel{
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
	record.Cash, _ = temp.ValueToString(task_pkg.FieldName_Cash.String(), row)

	record.Type, _ = temp.ValueToString(task_pkg.FieldName_Type.String(), row)
	record.Ccy, _ = temp.ValueToString(task_pkg.FieldName_Ccy.String(), row)
	if len(record.Ccy) == 0 {
		record.Ccy = "人民币"
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
func newRecordCreateManyFromExcelCommand(ctx context.Context, appcmd *command2.RecordImport2MasterCommand, s3Bucket, s3FilePath string) (*event.RecordImportMasterEvent, error) {
	if appcmd == nil {
		return nil, errors.New("appcmd 参数不能为空")
	}
	// 拼出 s3://user:pass@host/bucket/key 完整 URL 写入 S3FileUrl，
	// graph 域 RecordEventSubHandler 拿到后可直接传给 apoc.load.parquet
	minioCfg, ok := env.GetEnv().GetMinioByKey("default")
	if !ok {
		return nil, errors.New("env minio.default not configured")
	}
	s3FileUrl := fmt.Sprintf("s3://%s:%s@%s/%s/%s",
		minioCfg.AccessKey, minioCfg.SecretKey, minioCfg.Endpoint,
		s3Bucket, strings.TrimLeft(s3FilePath, "/"))

	eventData := &event.RecordImportMasterEventData{
		CaseId:     appcmd.Data.CaseId,
		DocId:      appcmd.Data.DocId,
		FileName:   appcmd.Data.FileName,
		FileId:     appcmd.Data.FileId,
		SheetId:    appcmd.Data.SheetId,
		SheetName:  appcmd.Data.SheetName,
		TaskId:     appcmd.Data.TaskId,
		S3FileUrl:  s3FileUrl,
		MasterType: appcmd.Data.MasterType,
		MasterId:   appcmd.Data.MasterId,
	}
	ev := event.NewRecordImportMasterEvent(ctx, config.ImportAppId, eventData)
	return ev, nil
}
