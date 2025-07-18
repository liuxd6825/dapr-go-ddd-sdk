package service

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/readexcel"
)

type ExcelService struct {
	fileService    *FileService
	sheetService   *SheetService
	rowService     *RowService
	docFileService *service.FileService
}

func NewExcelService() *ExcelService {
	return &ExcelService{
		fileService:    NewFileService(),
		sheetService:   NewSheetService(),
		rowService:     NewRowService(),
		docFileService: service.NewFileService(),
	}
}

func (s *ExcelService) CreateFile(ctx context.Context, cmd *command.ExcelCreateByFileCommand) (*model.ExcelFile, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	file, err := s.fileService.FindById(ctx, cmd.CaseId, cmd.FileId)
	if err != nil {
		return nil, err
	} else if file != nil {
		return file, err
	}

	file = &model.ExcelFile{
		Id:     cmd.FileId,
		DocId:  cmd.DocId,
		CaseId: cmd.CaseId,
		Name:   cmd.FileName,
	}
	if err = s.fileService.Create(ctx, file); err != nil {
		return nil, err
	}

	bytes, err := s.readExcelFile(ctx, file.Id)
	if err != nil {
		return nil, err
	}

	views, err := readexcel.ReadBytesToMap(bytes, "", 1000)
	if err != nil {
		return nil, err
	}

	file.Sheets = newSheetInfos(views.Sheets)
	if err = s.fileService.Update(ctx, file); err != nil {
		return nil, err
	}

	sheetInfo := file.GetSheet(views.OpenSheet)
	if s == nil {
		return nil, errors.New("%v 工作表不存在", views.OpenSheet)
	}

	sheet := newSheet(file, views.OpenSheet, sheetInfo.MaxRow, sheetInfo.MaxCol, views.Columns)
	if err = s.sheetService.Create(ctx, sheet); err != nil {
		return nil, err
	}

	rows := newRows(file, sheet.Id, views.Items)
	if err = s.rowService.CreateMany(ctx, rows); err != nil {
		return file, err
	}

	return file, nil
}

func (s *ExcelService) CreateRow(ctx context.Context, cmd *command.CreateRowsAppCmd) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	bytes, err := s.readExcelFile(ctx, cmd.FileId)
	if err != nil {
		return err
	}

	// 查找sheet是否存在，如果存在则退出，不需要创建。
	sheetId := fmt.Sprintf("%s(%s)", cmd.FileId, cmd.SheetName)
	if sheet, err := s.sheetService.FindById(ctx, sheetId); err != nil {
		return err
	} else if sheet != nil {
		return nil
	}

	// 加载excel文件
	views, err := readexcel.ReadBytesToMap(bytes, cmd.SheetName, 1000)

	file := &model.ExcelFile{
		Id:     cmd.FileId,
		DocId:  cmd.DocId,
		CaseId: cmd.CaseId,
		Name:   cmd.FileName,
		Sheets: newSheetInfos(views.Sheets),
	}

	if file, err := s.fileService.FindById(ctx, file.CaseId, file.Id); err != nil {
		return err
	} else if file != nil {
		if err = s.fileService.Create(ctx, file); err != nil {
			return err
		}
	}

	sheetInfo := file.GetSheet(cmd.SheetName)
	if s == nil {
		return errors.New("%v工作表不存在", cmd.SheetName)
	}
	sheet := newSheet(file, sheetInfo.Name, sheetInfo.MaxRow, sheetInfo.MaxCol, views.Columns)
	if err = s.sheetService.Create(ctx, sheet); err != nil {
		return err
	}

	rows := newRows(file, sheet.Id, views.Items)
	err = s.rowService.CreateMany(ctx, rows)
	if err != nil {
		return err
	}

	return nil
}

func (s *ExcelService) readExcelFile(ctx context.Context, fileId string) ([]byte, error) {
	bytes, err := s.docFileService.ReadByteByFileId(ctx, fileId)
	return bytes, err
}

// FindRows
// @Description: 获取指定excel的预览数据
// @receiver e
// @param ctx
// @param cmd
// @return *readexcel.Review
// @return error
func (s *ExcelService) FindRows(ctx context.Context, cmd *command.FindRowsAppQuery) (*command.FindRowsAppResult, bool, error) {
	if err := cmd.Validate(); err != nil {
		return nil, false, err
	}
	qry := &query.FindRowBySheetQuery{
		CaseId:  cmd.CaseId,
		DocId:   cmd.DocId,
		FileId:  cmd.FileId,
		SheetId: newSheetId(cmd.FileId, cmd.SheetName),
	}

	res := &command.FindRowsAppResult{}
	if sheet, err := s.sheetService.FindById(ctx, qry.SheetId); err != nil {
		return nil, false, err
	} else if sheet != nil {
		res.SheetName = sheet.Name
		res.Columns = sheet.Columns
		res.MaxRow = sheet.MaxRow
		res.MaxCol = sheet.MaxCol
	} else {
		return nil, false, errors.ErrNotFound
	}

	var data []map[string]any
	rows, err := s.rowService.FindBySheet(ctx, qry)
	for _, row := range rows {
		data = append(data, row.Values)
	}
	res.Rows = data
	return res, len(data) > 0, err
}

func (s *ExcelService) FindFile(ctx context.Context, qry *command.FindFileAppQuery) (*model.ExcelFile, error) {
	if err := qry.Validate(); err != nil {
		return nil, err
	}
	return s.fileService.FindById(ctx, qry.CaseId, qry.FileId)
}

// FindRecords
// @Description:
// @receiver e
// @param ctx
// @param appcmd
// @return []*model.RecordIe
// @return error
/*
func (e *excelViewAppService) FindRecords(ctx context.Context, appCmd *appcmd.FindRecordsAppQuery) ([]*model.RecordIe, error) {
	if err := appCmd.Validate(); err != nil {
		return nil, err
	}

	buffer, err := minioclient.GetClient().ReadFile(ctx, appCmd.TenantId, appCmd.FileId, appCmd.FileName)
	if err != nil {
		return nil, err
	}

	cmd := &record_ie.ReadByExcelFileAppCmd{
		TenantId:  appCmd.TenantId,
		CaseId:    appCmd.CaseId,
		DocId:     appCmd.DocId,
		FileId:    appCmd.FileId,
		Name: appCmd.Name,
		MaxRows:   appCmd.MaxRows,
		Template:  appCmd.Template,
		Buffer:    buffer,
	}
	return record_ie.NewRecordIeAppService().ReadByRows(ctx, cmd)
}
*/

/*
func (e *ExcelService) FindRecords(ctx context.Context, appCmd *command.FindRecordsAppQuery) ([]*model2.RecordIe, error) {
	if err := appCmd.Validate(); err != nil {
		return nil, err
	}
	cmd := &record_ie2.ReadByExcelRowAppCmd{
		TenantId:  appCmd.TenantId,
		CaseId:    appCmd.CaseId,
		DocId:     appCmd.DocId,
		FileId:    appCmd.FileId,
		SheetName: appCmd.SheetName,
		Template:  appCmd.Template,
	}
	return record_ie2.NewRecordIeAppService().ReadByRows(ctx, cmd)
}
*/
