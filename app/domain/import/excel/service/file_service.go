package service

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/excel/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/excel/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/singleutils"
)

type FileService struct {
	repos *dao.FileDao
}

func NewFileService() *FileService {
	return singleutils.CreateObj[*FileService](func() *FileService {
		return &FileService{
			repos: dao.NewFileDao(config.DBKey),
		}
	})
}

func (f *FileService) Create(ctx context.Context, m *model.ExcelFile) error {
	return f.repos.Create(ctx, m).GetError()
}

func (f *FileService) Update(ctx context.Context, m *model.ExcelFile) error {
	return f.repos.Update(ctx, m).GetError()
}

func (f *FileService) DeleteById(ctx context.Context, id string) {
	f.repos.DeleteById(ctx, id)
}

func (f *FileService) FindById(ctx context.Context, caseId, fileId string) (*model.ExcelFile, error) {
	return f.repos.FindOneByRSQL(ctx, fmt.Sprintf("case_id='%s' and file_id='%s'", caseId, fileId))
}

func (f *FileService) FindPaging(ctx context.Context, qry idao.FindPagingQuery) store.FindPagingResult[*model.ExcelFile] {
	return f.repos.FindPaging(ctx, qry)
}
