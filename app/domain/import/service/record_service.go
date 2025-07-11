package service

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/singleutils"
	"go.mongodb.org/mongo-driver/bson"
	"math"
	"strings"
)

type RecordService struct {
	repos *dao.RecordDao
}

func NewRecordService() *RecordService {
	return singleutils.CreateObj[*RecordService](func() *RecordService {
		return &RecordService{
			repos: dao.NewRecordDao(config.DBKey),
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

func (s *RecordService) Delete(ctx context.Context, appcmd *command.RecordIeDeleteCommand) error {
	return s.repos.DeleteById(ctx, appcmd.Data.Id).GetError()
}

func (s *RecordService) DeleteByTaskId(ctx context.Context, taskId string) error {
	return s.repos.DeleteById(ctx, taskId).GetError()
}

func (s *RecordService) Create(ctx context.Context, m *model.Record) error {
	return s.repos.Create(ctx, m).GetError()
}

func (s *RecordService) CreateMany(ctx context.Context, list []*model.Record) (int64, error) {
	for _, r := range list {
		data, err := s.Check(ctx, r)
		if err != nil {
			panic(err)
		}
		r.Errors = data
	}
	res := s.repos.CreateMany(ctx, list)
	return res.GetRowsAffected(), res.GetError()
}

func (s *RecordService) Update(ctx context.Context, cmd *command.RecordIeUpdateCommand) (*model.Record, error) {
	err := cmd.Validate()
	if err != nil {
		return nil, err
	}
	if cmd.IsValidOnly {
		return nil, err
	}

	data, err := maputils.NewMapWithOptions(cmd.Data, cmd.UpdateMask, false)
	if err != nil {
		return nil, err
	}
	setData := bson.M{
		"$set": data,
	}
	s.repos.UpdateMap(ctx, cmd.Data.Id, setData)
	record, err := s.repos.FindById(ctx, cmd.Data.Id)
	return record, err
}

func (s *RecordService) UpdateField(ctx context.Context, cmd *command.RecordIeUpdateFieldCommand) (*model.Record, bool, error) {
	var err error
	s.repos.UpdateMap(ctx, cmd.Data.Id, cmd.Data.Values)
	res, err := s.repos.FindById(ctx, cmd.Data.Id)
	res.Errors, err = s.Check(ctx, res)
	return res, true, err
}

func (s *RecordService) UpdateByFilter(ctx context.Context, cmd *command.RecordIeUpdateFilterCommand) error {
	if err := cmd.Validate(); err != nil {
		panic(err)
	}
	filter := fmt.Sprintf("%s and taskId=='%s'", cmd.Data.Filter, cmd.Data.TaskId)
	return s.repos.UpdateMapByRSQL(ctx, filter, cmd.Data.Values).GetError()
}

func (s *RecordService) FindById(ctx context.Context, tenantId, id string) (*model.Record, error) {
	return s.repos.FindById(ctx, id)
}

func (s *RecordService) FindPaging(ctx context.Context, qry idao.FindPagingQuery) store.FindPagingResult[*model.Record] {
	return s.repos.FindPaging(ctx, qry)
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
func (s *RecordService) FindPagingByTaskId(ctx context.Context, taskId string, isFindError bool) store.FindPagingResult[*model.Record] {
	var mustFilter string
	tenantId := appctx.GetTenantId2(ctx)
	if isFindError {
		mustFilter = s.getErrorSql(ctx, tenantId, taskId)
	} else {
		mustFilter = fmt.Sprintf("taskId=='%s'", taskId)
	}
	qry := store.NewFindPagingQueryRequest()
	qry.SetMustFilter(mustFilter)

	res := s.repos.FindPaging(ctx, qry)
	if res != nil {
		for _, v := range res.GetData() {
			v.Errors, _ = s.Check(ctx, v)
		}
	}

	return res
}

func (s *RecordService) CountErrorByTaskId(ctx context.Context, tenantId, taskId string) (int64, error) {
	rsql := s.getErrorSql(ctx, tenantId, taskId)
	return s.repos.CountByRSQL(ctx, rsql)
}

func (s *RecordService) getErrorSql(ctx context.Context, tenantId, taskId string) string {
	return fmt.Sprintf(`
			(   ((payout=null=0 or payout=='' or payout==0) and (income=null=0 or income=='' or income==0)) 
				or ((payout!=0) and (income!=0)) 
				or ((payout==0) and (income==0)) 
				or date=null=0 or date=='' 
				or name=null=0 or name=='' 
				or acct=null=0 or acct=='' 
				or bankName=null=0 or bankName=='' 
				or oppName=null=0 or oppName=='' 
				or oppAcct=null=0 or oppAcct==''
				or oppBankName=null=0 or oppBankName==''
				or amount=null=0 or amount==0 or amount==''
			) and taskId=='%s'
		`, taskId)
}

func (s *RecordService) Check(ctx context.Context, r *model.Record) (map[string][]string, error) {
	if r == nil {
		return nil, errors.New("检查的【流水记录】不能为空。")
	}
	data := map[string][]string{}
	if r.Date == nil {
		s.addFieldError(data, "date", "不能为空")
	}

	if *r.Payout != 0 && *r.Income != 0 {
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
