package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/readexcel"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/singleutils"
	"strings"
)

type SheetService struct {
	repos idao.Dao[*model.ExcelSheet]
}

func NewSheetService() *SheetService {
	return singleutils.CreateObj[*SheetService](func() *SheetService {
		return &SheetService{
			repos: dao.NewExcelSheetDao(config.DBKey),
		}
	})
}

func (f *SheetService) Create(ctx context.Context, m *model.ExcelSheet) error {
	return f.repos.Create(ctx, m).GetError()
}

func (f *SheetService) CreateMany(ctx context.Context, m []*model.ExcelSheet) error {
	return f.repos.CreateMany(ctx, m).GetError()
}

func (f *SheetService) Update(ctx context.Context, m *model.ExcelSheet) error {
	return f.repos.Update(ctx, m).GetError()
}

func (f *SheetService) DeleteById(ctx context.Context, id string) error {
	return f.repos.DeleteById(ctx, id).GetError()
}

func (f *SheetService) FindById(ctx context.Context, id string) (*model.ExcelSheet, error) {
	return f.repos.FindById(ctx, id)
}

func (f *SheetService) FindPaging(ctx context.Context, qry idao.FindPagingQuery) (idao.FindPagingResult[*model.ExcelSheet], error) {
	res := f.repos.FindPaging(ctx, qry)
	return res, res.GetError()
}

/*
func (f *SheetService) FindByName(ctx context.Context, qry *query.FindSheetByNameQuery) ([]*model.ExcelSheet, bool, error) {
	return f.repos.FindByName(ctx, qry)
}
*/

func newSheetInfos(items []*readexcel.Sheet) []*model.ExcelFileSheetInfo {
	var sheets []*model.ExcelFileSheetInfo
	for _, item := range items {
		sheet := &model.ExcelFileSheetInfo{
			Name:   item.Name,
			MaxCol: item.MaxCol,
			MaxRow: item.MaxRow,
		}
		sheets = append(sheets, sheet)
	}
	return sheets
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

func newSheetId(fileId, sheetName string) string {
	return fileId + sheetName
}

func newSheet(file *model.ExcelFile, sheetName string, maxRow int64, maxCol int64, columns []string) *model.ExcelSheet {
	sheet := &model.ExcelSheet{
		Id:       file.Id + strings.Trim(sheetName, " "),
		TenantId: file.TenantId,
		CaseId:   file.CaseId,
		DocId:    file.DocId,

		FileId:   file.Id,
		FileName: file.Name,
		Name:     sheetName,
		MaxRow:   maxRow,
		MaxCol:   maxCol,
		Columns:  columns,
	}
	return sheet
}
