package dao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
)

type ExcelSheetDao struct {
	idao.Dao[*model.ExcelSheet]
}

func NewExcelSheetDao(dbKey string) *ExcelSheetDao {
	tableName := "import_excel_sheet"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.ExcelSheet{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	daoVal := &ExcelSheetDao{Dao: dao.NewDao[*model.ExcelSheet](newCfg)}
	return daoVal
}

func (f *ExcelSheetDao) FindByName(ctx context.Context, qry *query.ExcelSheetFindByNameQuery) store.FindPagingResult[*model.ExcelSheet] {
	q := store.NewFindPagingQueryRequest()
	build := rsql.NewBuilder().And(
		rsql.Eq("file_id", qry.FileId),
		rsql.Eq("name", qry.Name),
	)
	q.SetMustFilter(build.Build())
	return f.Dao.FindPaging(ctx, q)
}
