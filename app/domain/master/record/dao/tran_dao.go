package dao

import (
	"context"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
)

type TranDao struct {
	idao.Dao[*model.Tran]
}

func NewTranDao(dbKey string) *TranDao {
	tableName := "master_tran"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Tran{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.Tran](newCfg)
	daoVal := &TranDao{Dao: baseDao}
	return daoVal
}

func (s *TranDao) FindByAccountOppAccount(ctx context.Context, account string, startDate, endDate *time.Time) ([]*model.Tran, error) {

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
	qry.SetPageSize(store.PagingMaxPageSize)

	res := s.FindPaging(ctx, qry)
	return res.GetData(), res.GetError()
}
