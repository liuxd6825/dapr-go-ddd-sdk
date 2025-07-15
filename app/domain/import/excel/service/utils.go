package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/excel/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/readexcel"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"strings"
)

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
