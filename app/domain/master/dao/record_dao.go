package dao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
	"time"
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
