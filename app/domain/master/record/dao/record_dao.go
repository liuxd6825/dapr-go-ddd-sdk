package dao

import (
	"context"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
)

type RecordDao struct {
	idao.Dao[*model.Record]
}

func NewRecordDao(dbKey string) *RecordDao {
	tableName := "master_record"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Record{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.Record](newCfg)
	daoVal := &RecordDao{Dao: baseDao}
	return daoVal
}

func (s *RecordDao) FindByAccountOppAccount(ctx context.Context, account string, startDate, endDate time.Time) ([]*model.Record, error) {
	builder := rsql.NewBuilder().And(
		rsql.Or(
			rsql.Eq("acct", account),
			rsql.Eq("oppAcct", account),
		),
		rsql.Gte("date", startDate),
		rsql.Lte("date", endDate),
	)

	qry := store.NewFindPagingQuery()
	qry.SetMustFilter(builder.Build())
	qry.SetSort("date:asc")
	qry.SetIsTotalRows(false)
	qry.SetPageSize(100000000000)

	res := s.FindPaging(ctx, qry)
	return res.GetData(), res.GetError()
}

// DistinctHuman
// @Description: 通过流水数据统计人员
// @receiver s
// @param ctx
// @param caseId
// @param name
// @return []*query.RecordAccountFindByNameResult
// @return error
func (s *RecordDao) DistinctHuman(ctx context.Context, caseId, masterType, masterId, myFilter, oppFilter string) ([]*query.DistinctNameResult, error) {
	var res []*query.DistinctNameResult
	keys := map[string]bool{}
	err := s.findDistinct(ctx, "name", caseId, masterType, masterId, myFilter, func(list []*model.Record) {
		for _, item := range list {
			if _, ok := keys[item.Name]; !ok {
				keys[item.Name] = true
				result := &query.DistinctNameResult{
					Name: item.Name,
				}
				res = append(res, result)
			}
		}
	})
	if err != nil {
		return nil, err
	}

	err = s.findDistinct(ctx, "opp_name", caseId, masterType, masterId, oppFilter, func(list []*model.Record) {
		for _, item := range list {
			if _, ok := keys[item.OppAcct]; !ok {
				keys[item.OppName] = true
				record := &query.DistinctNameResult{
					Name: item.OppName,
				}
				res = append(res, record)
			}
		}
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DistinctName
// @Description: 通过流水数据统计公司
// @receiver s
// @param ctx
// @param caseId
// @param name
// @return []*query.RecordAccountFindByNameResult
// @return error
func (s *RecordDao) DistinctName(ctx context.Context, caseId, masterType, masterId, myFilter, oppFilter string) ([]*query.DistinctNameResult, error) {
	var res []*query.DistinctNameResult
	keys := map[string]bool{}
	err := s.findDistinct(ctx, "name", caseId, masterType, masterId, myFilter, func(list []*model.Record) {
		for _, item := range list {
			if _, ok := keys[item.Name]; !ok {
				keys[item.Name] = true
				result := &query.DistinctNameResult{
					Name: item.Name,
				}
				res = append(res, result)
			}
		}
	})
	if err != nil {
		return nil, err
	}

	err = s.findDistinct(ctx, "opp_name", caseId, masterType, masterId, oppFilter, func(list []*model.Record) {
		for _, item := range list {
			if _, ok := keys[item.OppName]; !ok {
				keys[item.OppName] = true
				record := &query.DistinctNameResult{
					Name: item.OppName,
				}
				res = append(res, record)
			}
		}
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DistinctCompany
// @Description: 通过流水数据统计公司
// @receiver s
// @param ctx
// @param caseId
// @param name
// @return []*query.RecordAccountFindByNameResult
// @return error
func (s *RecordDao) DistinctCompany(ctx context.Context, caseId, masterType, masterId, myFilter, oppFilter string) ([]*query.DistinctNameResult, error) {
	var res []*query.DistinctNameResult
	keys := map[string]bool{}
	err := s.findDistinct(ctx, "name", caseId, masterType, masterId, myFilter, func(list []*model.Record) {
		for _, item := range list {
			if _, ok := keys[item.Name]; !ok {
				keys[item.Name] = true
				result := &query.DistinctNameResult{
					Name: item.Name,
				}
				res = append(res, result)
			}
		}
	})
	if err != nil {
		return nil, err
	}

	err = s.findDistinct(ctx, "opp_name", caseId, masterType, masterId, oppFilter, func(list []*model.Record) {
		for _, item := range list {
			if _, ok := keys[item.OppName]; !ok {
				keys[item.OppName] = true
				record := &query.DistinctNameResult{
					Name: item.OppName,
				}
				res = append(res, record)
			}
		}
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DistinctAccount
// @Description: 通过流水数据统计银行账号
// @receiver s
// @param ctx
// @param caseId
// @param name
// @return []*query.RecordAccountFindByNameResult
// @return error
func (s *RecordDao) DistinctAccount(ctx context.Context, caseId, masterType, masterId, myFilter, oppFilter string, findOpp bool) ([]*query.DistinctAccountResult, error) {
	var accounts []*query.DistinctAccountResult
	keys := map[string]bool{}
	err := s.findDistinct(ctx, "acct, bank_name", caseId, masterType, masterId, myFilter, func(list []*model.Record) {
		for _, item := range list {
			if _, ok := keys[item.Acct]; !ok {
				keys[item.Acct] = true
				result := &query.DistinctAccountResult{
					Account:   item.Acct,
					BankName:  item.BankName,
					OwnerName: item.Name,
				}
				accounts = append(accounts, result)
			}
		}
	})
	if err != nil {
		return nil, err
	}

	if findOpp {
		err = s.findDistinct(ctx, "opp_acct, opp_bank_name", caseId, masterType, masterId, oppFilter, func(list []*model.Record) {
			for _, item := range list {
				if _, ok := keys[item.OppAcct]; !ok {
					keys[item.OppAcct] = true
					record := &query.DistinctAccountResult{
						Account:   item.OppAcct,
						BankName:  item.OppBankName,
						OwnerName: item.OppName,
					}
					accounts = append(accounts, record)
				}
			}
		})
		if err != nil {
			return nil, err
		}
	}

	return accounts, nil
}

func (s *RecordDao) findDistinct(ctx context.Context, fields string, caseId, masterType, masterId, filter string, setData func(list []*model.Record)) error {
	sql, err := s.getDistinctRSQL(ctx, caseId, masterType, masterId, filter)
	if err != nil {
		return err
	}
	qry := store.NewFindDistinctQuery()
	qry.SetFields(fields)
	qry.SetFilter(sql)
	res := s.FindDistinct(ctx, qry)
	if res.GetError() != nil {
		return res.GetError()
	}

	if setData != nil {
		setData(res.GetData())
	}
	return nil
}

func (s *RecordDao) getDistinctRSQL(ctx context.Context, caseId, masterType, masterId string, filter string) (string, error) {
	var ands1 []rsql.Condition
	if caseId != "" {
		ands1 = append(ands1, rsql.Eq("case_id", caseId))
	}
	if masterType != "" {
		ands1 = append(ands1, rsql.Eq("master_type", masterType))
	}
	if masterId != "" {
		ands1 = append(ands1, rsql.Eq("master_id", masterId))
	}
	str1 := rsql.NewBuilder().And(ands1...).Build()
	if filter != "" {
		str1 = str1 + " AND " + filter
	}
	return str1, nil
}
