package service

import (
	"context"
	"fmt"
	docfile "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	dao2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/readexcel"

	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/singleutils"
	"go.mongodb.org/mongo-driver/bson"
	"math"
	"strings"
)

type RecordService struct {
	dao            *dao2.RecordDao
	docFileService *docfile.FileService
	taskService    *TaskService
}

func NewRecordService() *RecordService {
	return singleutils.CreateObj[*RecordService](func() *RecordService {
		return &RecordService{
			dao:            dao2.NewRecordDao(config.DBKey),
			docFileService: docfile.NewFileService(),
			taskService:    NewTaskService(),
		}
	})
}

func (s *RecordService) addFieldError(data map[string][]string, field string, message string) {
	if list, ok := data[field]; ok {
		list = append(list, message)
		data[field] = list
	} else {
		list = []string{message}
		data[field] = list
	}
}

func (s *RecordService) Create(ctx context.Context, cmd *command.RecordCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		dr := cmd.Data.Record
		record := &model.RecordIe{
			BaseModel: xbase.BaseModel{
				TenantId: cmd.Data.TenantId,
			},
			DocId:       cmd.Data.DocId,
			FileId:      cmd.Data.FileId,
			TaskId:      dr.TaskId,
			RowNum:      dr.RowNum,
			Name:        dr.Name,
			Acct:        dr.Acct,
			AcctType:    dr.AcctType,
			Category:    dr.Category,
			BankName:    dr.BankName,
			OppIden:     dr.OppIden,
			OppName:     dr.OppName,
			OppAcct:     dr.OppAcct,
			OppAcctType: dr.OppAcctType,
			OppBankName: dr.OppBankName,
			Amount:      dr.Amount,
			OppCategory: dr.OppCategory,
			Payout:      dr.Payout,
			Income:      dr.Income,
			Summary:     dr.Summary,
			Date:        dr.Date,
			Type:        dr.Type,
			Balance:     dr.Balance,
			Serial:      dr.Serial,
		}
		return s.dao.Create(ctx, record).GetError()
	})
}

func (s *RecordService) Delete(ctx context.Context, appcmd *command.RecordDeleteCommand) error {
	return s.dao.DeleteById(ctx, appcmd.Data.Id).GetError()
}

func (s *RecordService) DeleteByTaskId(ctx context.Context, taskId string) error {
	return s.dao.DeleteById(ctx, taskId).GetError()
}

func (s *RecordService) CreateMany(ctx context.Context, list []*model.RecordIe) (int64, error) {
	for _, r := range list {
		data, err := s.Check(ctx, r)
		if err != nil {
			panic(err)
		}
		r.Errors = data
	}
	res := s.dao.CreateMany(ctx, list)
	return res.GetRowsAffected(), res.GetError()
}

func (s *RecordService) Update(ctx context.Context, cmd *command.RecordUpdateCommand) (recordIe *model.RecordIe, err error) {
	err = xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		data, err := maputils.NewMapWithOptions(cmd.Data, cmd.UpdateMask, false)
		if err != nil {
			return err
		}
		setData := bson.M{
			"$set": data,
		}
		s.dao.UpdateMap(ctx, cmd.Data.Id, setData)
		recordIe, err = s.dao.FindById(ctx, cmd.Data.Id)
		return err
	})
	return recordIe, err
}

func (s *RecordService) UpdateField(ctx context.Context, cmd *command.RecordUpdateFieldCommand) (*model.RecordIe, bool, error) {
	var err error
	s.dao.UpdateMap(ctx, cmd.Data.Id, cmd.Data.Values)
	res, err := s.dao.FindById(ctx, cmd.Data.Id)
	if err != nil {
		return nil, false, err
	}
	if res == nil {
		return nil, false, nil
	}
	res.Errors, err = s.Check(ctx, res)
	return res, true, err
}

func (s *RecordService) UpdateByFilter(ctx context.Context, cmd *command.RecordUpdateFilterCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		filter := fmt.Sprintf("%s and taskId=='%s'", cmd.Data.Filter, cmd.Data.TaskId)
		return s.dao.UpdateMapByRSQL(ctx, filter, cmd.Data.Values).GetError()
	})
}

func (s *RecordService) FindById(ctx context.Context, id string) (*model.RecordIe, error) {
	return s.dao.FindById(ctx, id)
}

func (s *RecordService) FindPaging(ctx context.Context, qry idao.FindPagingQuery) store.FindPagingResult[*model.RecordIe] {
	return s.dao.FindPaging(ctx, qry)
}

// FindPagingByTaskId
// @Description: 按任务ID查找错误数据
// @receiver s
// @param ctx
// @param qry
// @param opts
// @return *ddd_repository.FindPagingResult[*model.Record]
// @return bool
// @return error
func (s *RecordService) FindPagingByTaskId(ctx context.Context, qry *query.RecordIeFindPagingByTaskIdQuery) store.FindPagingResult[*model.RecordIe] {
	var mustFilter string
	tenantId := appctx.GetTenantId2(ctx)
	taskId := qry.TaskId
	if qry.IsFindError {
		mustFilter = s.getErrorSql(ctx, tenantId, qry.TaskId)
	} else {
		mustFilter = fmt.Sprintf("task_id=='%s'", taskId)
	}
	qry.SetMustFilter(mustFilter)

	res := s.dao.FindPaging(ctx, qry)
	if res != nil {
		for _, v := range res.GetData() {
			v.Errors, _ = s.Check(ctx, v)
		}
	}

	return res
}

func (s *RecordService) CountErrorByTaskId(ctx context.Context, tenantId, taskId string) (int64, error) {
	rsql := s.getErrorSql(ctx, tenantId, taskId)
	return s.dao.CountByRSQL(ctx, rsql)
}

func (s *RecordService) getErrorSql(ctx context.Context, tenantId, taskId string) string {
	return fmt.Sprintf(`
			( 	(payout=!null=0 and income=!null=0) 
				or (payout=null=0 and income=null=0) 
				or date=null=0 or date=='' 
				or name=null=0 or name=='' 
				or acct=null=0 or acct=='' 
				or bank_name=null=0 or bank_name=='' 
				or opp_name=null=0 or opp_name=='' 
				or opp_acct=null=0 or opp_acct==''
				or opp_bank_name=null=0 or opp_bank_name==''
				or amount=null=0 or amount==0 or amount==''
			) and taskId=='%s'
		`, taskId)
}

func (s *RecordService) Check(ctx context.Context, r *model.RecordIe) (map[string][]string, error) {
	if s == nil {
		return nil, errors.New("检查的【流水记录】不能为空。")
	}
	data := map[string][]string{}
	if r.Date.IsNil() {
		s.addFieldError(data, "date", "不能为空")
	}

	if r.Payout != nil && r.Income != nil {
		s.addFieldError(data, "payout", "【支出金额】和【收入金额】不能同时有值")
		s.addFieldError(data, "income", "【支出金额】和【收入金额】不能同时有值")
	} else if (r.Payout == nil || *r.Payout >= 0) && (r.Income == nil || *r.Income <= 0) {
		s.addFieldError(data, "payout", "不能为空或大于0")
		s.addFieldError(data, "income", "不能为空或小于0")
	} else if (r.Payout != nil && *r.Payout != 0) && r.Amount != nil && math.Abs(*r.Payout) != *r.Amount {
		s.addFieldError(data, "amount", "必须等于【支出金额】")
	} else if (r.Income != nil && *r.Income != 0) && r.Amount != nil && *r.Income != *r.Amount {
		s.addFieldError(data, "amount", "必须等于【收入金额】")
	}

	if r.Amount == nil || *r.Amount <= 0 {
		s.addFieldError(data, "amount", "不能为空或小于0")
	}

	if len(strings.ReplaceAll(r.Acct, " ", "")) == 0 {
		s.addFieldError(data, "acct", "不能为空")
	}
	if len(strings.ReplaceAll(r.Name, " ", "")) == 0 {
		s.addFieldError(data, "name", "不能为空")
	}
	if len(strings.ReplaceAll(r.BankName, " ", "")) == 0 {
		s.addFieldError(data, "bankName", "不能为空")
	}

	if len(strings.ReplaceAll(r.OppAcct, " ", "")) == 0 {
		s.addFieldError(data, "oppAcct", "不能为空")
	}
	if len(strings.ReplaceAll(r.OppName, " ", "")) == 0 {
		s.addFieldError(data, "oppName", "不能为空")
	}
	if len(strings.ReplaceAll(r.OppBankName, " ", "")) == 0 {
		s.addFieldError(data, "oppBankName", "不能为空")
	}
	if len(data) > 0 {
		return data, nil
	}
	return nil, nil
}

func (s *RecordService) readMap(ctx context.Context, task *TaskOptions, temp *model.RecordTemplate, mapList []map[string]any) ([]*model.RecordIe, error) {
	if temp == nil {
		return nil, errors.New("parameter 'task.Template' is null")
	}

	tmp, err := temp.NewTemplate()
	if err != nil {
		return nil, err
	}

	newCtx := NewContext(ctx, task)
	_, list, err := readexcel.ReadMapToEntity[*model.RecordIe](newCtx, mapList, tmp, newRecord)
	if err != nil {
		return nil, err
	}

	return list, err
}
