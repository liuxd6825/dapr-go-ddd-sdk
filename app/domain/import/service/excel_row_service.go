package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/singleutils"
)

type ExcelRowService struct {
	rowDao       *dao.ExcelRowDao
	sheetService *ExcelSheetService
}

func NewExcelRowService() *ExcelRowService {
	return singleutils.CreateObj[*ExcelRowService](func() *ExcelRowService {
		return &ExcelRowService{
			rowDao:       dao.NewRowDao(config.DBKey),
			sheetService: NewExcelSheetService(),
		}
	})
}

func (f *ExcelRowService) Create(ctx context.Context, cmd *command.ExcelRowCreateCommand) error {
	m := &model.ExcelRow{
		Id:       cmd.Data.Id,
		TenantId: cmd.Data.TenantId,
		CaseId:   cmd.Data.CaseId,
		DocId:    cmd.Data.DocId,
		FileId:   cmd.Data.FileId,
		SheetId:  cmd.Data.SheetId,
		RowNum:   cmd.Data.RowNum,
		Values:   cmd.Data.Values,
	}
	return f.rowDao.Create(ctx, m).GetError()
}

func (f *ExcelRowService) CreateMany(ctx context.Context, list []*model.ExcelRow) error {
	return f.rowDao.CreateMany(ctx, list).GetError()
}

// FindRows
// @Description: 获取指定excel的预览数据
// @receiver e
// @param ctx
// @param cmd
// @return *query.ExcelRowFindPreviewQueryResult
// @return error
func (f *ExcelRowService) FindRows(ctx context.Context, qry *query.ExcelRowFindPreviewQuery) (*query.ExcelRowFindPreviewQueryResult, error) {
	return xbase.DoQuery2[*query.ExcelRowFindPreviewQueryResult](ctx, qry, func(ctx context.Context) (*query.ExcelRowFindPreviewQueryResult, error) {

		res := &query.ExcelRowFindPreviewQueryResult{}
		if sheet, err := f.sheetService.FindById(ctx, qry.SheetId); err != nil {
			return nil, err
		} else if sheet != nil {
			res.SheetName = sheet.Name
			res.Columns = sheet.Columns
			res.MaxRow = sheet.MaxRow
			res.MaxCol = sheet.MaxCol
		} else {
			return nil, errors.ErrNotFound
		}

		var data []map[string]any
		rows, err := f.FindBySheet(ctx, qry)
		for _, row := range rows {
			data = append(data, row.Values)
		}
		res.Rows = data
		return res, err
	})
}

func (f *ExcelRowService) Update(ctx context.Context, data *model.ExcelRow, updateMask []string) error {
	return f.rowDao.Update(ctx, data, idao.NewCallOptions().SetUpdateFields(updateMask)).GetError()
}

func (f *ExcelRowService) DeleteById(ctx context.Context, id string) error {
	return f.rowDao.DeleteById(ctx, id).GetError()
}

func (f *ExcelRowService) FindById(ctx context.Context, id string) (*model.ExcelRow, error) {
	return f.rowDao.FindById(ctx, id)
}

func (f *ExcelRowService) DeleteByCaseId(ctx context.Context, caseId string) error {
	build := rsql.NewBuilder().And(
		rsql.Eq("case_id", caseId),
	)
	return f.rowDao.DeleteByRSQL(ctx, build.Build()).GetError()
}

func (f *ExcelRowService) DeleteByFileId(ctx context.Context, fileId string) error {
	build := rsql.NewBuilder().And(
		rsql.Eq("file_id", fileId),
	)
	return f.rowDao.DeleteByRSQL(ctx, build.Build()).GetError()
}

func (f *ExcelRowService) DeleteBySheetId(ctx context.Context, sheetId string) error {
	build := rsql.NewBuilder().And(
		rsql.Eq("sheet_id", sheetId),
	)
	return f.rowDao.DeleteByRSQL(ctx, build.Build()).GetError()
}

func (f *ExcelRowService) FindBySheet(ctx context.Context, qry *query.ExcelRowFindPreviewQuery) ([]*model.ExcelRow, error) {
	rb := rsql.NewBuilder().And(
		rsql.Eq("case_id", qry.CaseId),
		rsql.Eq("sheet_id", qry.SheetId),
	)
	pagingBuilder := idao.NewFindPagingQueryBuilder().SetPageNum(0).SetPageSize(qry.MaxRows)
	pagingBuilder.SetMustFilter(rb.Build())
	pagingBuilder.SetSort("row_num:asc")
	res := f.FindPaging(ctx, pagingBuilder.Build())
	return res.GetData(), nil
}

func (f *ExcelRowService) FindPaging(ctx context.Context, qry idao.FindPagingQuery) idao.FindPagingResult[*model.ExcelRow] {
	return f.rowDao.FindPaging(ctx, qry)
}
