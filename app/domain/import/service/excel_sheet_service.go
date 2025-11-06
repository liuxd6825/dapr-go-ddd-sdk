package service

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/pkg/readexcel"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/singleutils"
)

type ExcelSheetService struct {
	repos idao.Dao[*model.ExcelSheet]
}

func NewExcelSheetService() *ExcelSheetService {
	return singleutils.CreateObj[*ExcelSheetService](func() *ExcelSheetService {
		return &ExcelSheetService{
			repos: dao.NewExcelSheetDao(config.DBKey),
		}
	})
}

func (f *ExcelSheetService) Create(ctx context.Context, m *model.ExcelSheet) error {
	return f.repos.Create(ctx, m).GetError()
}

func (f *ExcelSheetService) CreateMany(ctx context.Context, m []*model.ExcelSheet) error {
	return f.repos.CreateMany(ctx, m).GetError()
}

func (f *ExcelSheetService) Update(ctx context.Context, m *model.ExcelSheet) error {
	return f.repos.Update(ctx, m).GetError()
}

func (f *ExcelSheetService) DeleteById(ctx context.Context, id string) error {
	return f.repos.DeleteById(ctx, id).GetError()
}

func (f *ExcelSheetService) DeleteByFileId(ctx context.Context, fileId string) error {
	build := rsql.NewBuilder().And(
		rsql.Eq("file_Id", fileId),
	)
	return f.repos.DeleteByRSQL(ctx, build.Build()).GetError()
}

func (f *ExcelSheetService) DeleteByDocId(ctx context.Context, docId string) error {
	build := rsql.NewBuilder().And(
		rsql.Eq("doc_Id", docId),
	)
	return f.repos.DeleteByRSQL(ctx, build.Build()).GetError()
}

func (f *ExcelSheetService) FindById(ctx context.Context, id string) (*model.ExcelSheet, error) {
	return f.repos.FindById(ctx, id)
}

func (f *ExcelSheetService) FindByFileId(ctx context.Context, fileId string) ([]*model.ExcelSheet, error) {
	build := rsql.NewBuilder().And(
		rsql.Eq("file_id", fileId),
	)
	return f.repos.FindByRSQL(ctx, build.Build())
}

func (f *ExcelSheetService) FindByDocFileId(ctx context.Context, docFileId string) ([]*model.ExcelSheet, error) {
	build := rsql.NewBuilder().And(
		rsql.Eq("doc_file_id", docFileId),
	)
	return f.repos.FindByRSQL(ctx, build.Build())
}

func (f *ExcelSheetService) FindPaging(ctx context.Context, qry idao.FindPagingQuery) (idao.FindPagingResult[*model.ExcelSheet], error) {
	res := f.repos.FindPaging(ctx, qry)
	return res, res.GetError()
}

func newRows(file *model.ExcelFile, sheetId string, items []readexcel.ReviewItem) []*model.ExcelRow {
	var rows []*model.ExcelRow
	for rowNum, item := range items {
		row := &model.ExcelRow{
			Id:       idutils.NewId(),
			TenantId: file.TenantId,
			CaseId:   file.CaseId,
			DocId:    file.DocId,
			FileId:   file.Id,
			SheetId:  sheetId,
			RowNum:   int64(rowNum),
			Values:   item,
		}
		rows = append(rows, row)
	}
	return rows
}

func newSheet(file *model.ExcelFile, sheetName string, maxRow int64, maxCol int64, columns []string) *model.ExcelSheet {
	sheet := &model.ExcelSheet{
		Id:        idutils.NewId(), //file.Id,
		TenantId:  file.TenantId,
		CaseId:    file.CaseId,
		DocId:     file.DocId,
		DocFileId: file.DocFileId,
		FileId:    file.Id,
		FileName:  file.Name,
		Name:      sheetName,
		MaxRow:    maxRow,
		MaxCol:    maxCol,
		Columns:   columns,
	}
	return sheet
}
