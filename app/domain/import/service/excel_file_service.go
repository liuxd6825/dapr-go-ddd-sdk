package service

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/pkg/readexcel"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
	db "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/singleutils"
)

// ExcelFileService excel文件服务 从文档中心导入的文件
type ExcelFileService struct {
	dao            *dao.ExcelFileDao
	docFileService *service.FileService
	sheetService   *ExcelSheetService
	rowService     *ExcelRowService
}

func NewExcelFileService() *ExcelFileService {
	return singleutils.CreateObj[*ExcelFileService](func() *ExcelFileService {
		return &ExcelFileService{
			dao:            dao.NewExcelFileDao(config.DBKey),
			docFileService: service.NewFileService(),
			sheetService:   NewExcelSheetService(),
			rowService:     NewExcelRowService(),
		}
	})
}

func (s *ExcelFileService) Create(ctx context.Context, cmd *command.ExcelFileCreateCommand) (*model.ExcelFile, error) {
	return xbase.DoCommand2[*model.ExcelFile](ctx, cmd, func(ctx context.Context) (*model.ExcelFile, error) {
		file, err := s.FindByDocFileId(ctx, cmd.Data.DocFileId)
		if err != nil {
			return nil, err
		} else if file != nil {
			return file, err
		}

		file = &model.ExcelFile{
			Id:        cmd.Data.Id,
			Name:      cmd.Data.FileName,
			DocId:     cmd.Data.DocId,
			DocFileId: cmd.Data.DocFileId,
			CaseId:    cmd.Data.CaseId,
		}
		bytes, err := s.readExcelFile(ctx, cmd.Data.DocFileId)
		if err != nil {
			return nil, err
		}

		views, err := readexcel.ReadBytesToMap(bytes, 1000)
		if err != nil {
			return nil, err
		}

		// 需要开启事物
		err = db.StartTx(ctx, db.NewTxCfg(config.DBKey), func(txCtx context.Context, options ...*store2.SessionOptions) error {
			res := s.dao.Create(ctx, file)
			if res.Error != nil {
				return res.Error
			}
			for _, item := range views.Sheets {
				sheet := newSheet(file, item.Name, item.MaxRow, item.MaxCol, item.Columns)
				sheet.FileName = file.Name
				sheet.FileId = file.Id
				sheet.DocId = cmd.Data.DocId
				sheet.DocFileId = cmd.Data.DocFileId
				if err = s.sheetService.Create(ctx, sheet); err != nil {
					return err
				}
				rows := newRows(file, sheet.Id, item.Items)
				if err = s.rowService.CreateMany(ctx, rows); err != nil {
					return err
				}
			}
			return nil
		})
		return file, err
	})

}

func (s *ExcelFileService) Update(ctx context.Context, m *command.ExcelFileUpdateCommand) error {
	return xbase.DoCommand(ctx, m, func(ctx context.Context) error {
		return s.dao.Update(ctx, m.Data).GetError()
	})
}

func (s *ExcelFileService) Delete(ctx context.Context, m *command.ExcelFileDeleteCommand) error {
	return xbase.DoCommand(ctx, m, func(ctx context.Context) error {
		return s.dao.Delete(ctx, m.Data).GetError()
	})
}

func (s *ExcelFileService) DeleteById(ctx context.Context, m *command.ExcelFileDeleteByIdCommand) error {
	return xbase.DoCommand(ctx, m, func(ctx context.Context) error {
		if err := s.sheetService.DeleteByFileId(ctx, m.Data.Id); err != nil {
			return err
		}
		return s.dao.DeleteById(ctx, m.Data.Id).GetError()
	})
}
func (s *ExcelFileService) FindById(ctx context.Context, fileId string) (*model.ExcelFile, error) {
	return s.dao.FindById(ctx, fileId)
}

func (s *ExcelFileService) FindByFileId(ctx context.Context, fileId string) ([]*model.ExcelFile, error) {
	build := rsql.NewBuilder().And(
		rsql.Eq("file_id", fileId),
	)
	return s.dao.FindByRSQL(ctx, build.Build())
}

func (s *ExcelFileService) FindByDocFileId(ctx context.Context, docFileId string) (*model.ExcelFile, error) {
	build := rsql.NewBuilder().And(
		rsql.Eq("doc_file_id", docFileId),
	)
	return s.dao.FindOneByRSQL(ctx, build.Build())
}

func (s *ExcelFileService) FindPaging(ctx context.Context, qry idao.FindPagingQuery) store2.FindPagingResult[*model.ExcelFile] {
	return s.dao.FindPaging(ctx, qry)
}

func (s *ExcelFileService) readExcelFile(ctx context.Context, fileId string) ([]byte, error) {
	bytes, err := s.docFileService.ReadByteByFileId(ctx, fileId)
	return bytes, err
}
