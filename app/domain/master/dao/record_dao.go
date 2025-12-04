package dao

import (
	"context"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/query"
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

// AccountFindByName
// @Description: 通过流水数据统计银行账号
// @receiver s
// @param ctx
// @param caseId
// @param name
// @return []*query.RecordAccountFindByNameResult
// @return error
func (s *RecordDao) AccountFindByName(ctx context.Context, caseId, name, masterType, masterId string) ([]*query.RecordAccountFindByNameResult, error) {
	sql1, sql2, err := s.getAccountFindByNameRSQL(ctx, caseId, name, masterType, masterId)
	if err != nil {
		return nil, err
	}
	qry1 := store.NewFindDistinctQuery()
	qry1.SetFields("acct, bank_name")
	qry1.SetFilter(sql1)
	res1 := s.FindDistinct(ctx, qry1)
	if res1.GetError() != nil {
		return nil, res1.GetError()
	}

	qry2 := store.NewFindDistinctQuery()
	qry2.SetFields("opp_acct, opp_bank_name")
	qry2.SetFilter(sql2)
	res2 := s.FindDistinct(ctx, qry2)
	if res2.GetError() != nil {
		return nil, res2.GetError()
	}

	accounts := []*query.RecordAccountFindByNameResult{}
	keys := map[string]bool{}
	for _, item := range res1.GetData() {
		if _, ok := keys[item.Acct]; !ok {
			keys[item.Acct] = true
			result := &query.RecordAccountFindByNameResult{
				Account:   item.Acct,
				BankName:  item.BankName,
				OwnerName: item.Name,
			}
			accounts = append(accounts, result)
		}
	}
	for _, item := range res2.GetData() {
		if _, ok := keys[item.OppAcct]; !ok {
			keys[item.OppAcct] = true
			record := &query.RecordAccountFindByNameResult{
				Account:   item.OppAcct,
				BankName:  item.OppBankName,
				OwnerName: item.OppName,
			}
			accounts = append(accounts, record)
		}
	}

	return accounts, nil
}

func (s *RecordDao) getAccountFindByNameRSQL(ctx context.Context, caseId, name, masterType, masterId string) (string, string, error) {
	var ands1 []rsql.Condition
	var ands2 []rsql.Condition
	if caseId != "" {
		ands1 = append(ands1, rsql.Eq("case_id", caseId))
		ands2 = append(ands2, rsql.Eq("case_id", caseId))
	}
	if masterType != "" {
		ands1 = append(ands1, rsql.Eq("master_type", masterType))
		ands2 = append(ands2, rsql.Eq("master_type", masterType))
	}
	if masterId != "" {
		ands1 = append(ands1, rsql.Eq("master_id", masterId))
		ands2 = append(ands2, rsql.Eq("master_id", masterId))
	}
	if name != "" {
		ands1 = append(ands1, rsql.Eq("name", name))
		ands2 = append(ands2, rsql.Eq("opp_name", name))
	}
	str1 := rsql.NewBuilder().And(ands1...).Build()
	str2 := rsql.NewBuilder().And(ands2...).Build()
	return str1, str2, nil
}
