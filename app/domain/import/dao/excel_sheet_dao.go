package dao

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type ExcelSheetDao struct {
	idao.Dao[*model.ExcelSheet]
}

func NewExcelSheetDao(dbKey string) *ExcelSheetDao {
	tableName := "import_excel_sheet"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.ExcelSheet{}, tableName)
	newCfg := &dao.NewConfig{
		DBKey:      dbKey,
		IsPubEvent: dao.IsFalse(),
		TableName:  tableName,
		DBSchema:   dbSch,
	}
	daoVal := &ExcelSheetDao{Dao: dao.NewDao[*model.ExcelSheet](newCfg)}
	return daoVal
}

func (f *ExcelSheetDao) FindByName(ctx context.Context, qry *query.FindSheetByNameQuery) store.FindPagingResult[*model.ExcelSheet] {
	q := store.NewFindPagingQueryRequest()
	filter := fmt.Sprintf(`case_id=="%v" and doc_id=="%v" and file_id=="%v" and name=="%v" `, qry.CaseId, qry.DocId, qry.FileId, qry.Name)
	q.SetMustFilter(filter)
	return f.Dao.FindPaging(ctx, q)
}
