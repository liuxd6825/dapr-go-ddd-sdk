package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/singleutils"
)

type FileService struct {
	repos *dao.ExcelFileDao
}

func NewFileService() *FileService {
	return singleutils.CreateObj[*FileService](func() *FileService {
		return &FileService{
			repos: dao.NewExcelFileDao(config.DBKey),
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
	build := rsql.NewBuilder().And(
		rsql.Eq("case_id", caseId),
		rsql.Eq("field_id", fileId),
	)
	return f.repos.FindOneByRSQL(ctx, build.Build())
}

func (f *FileService) FindPaging(ctx context.Context, qry idao.FindPagingQuery) store.FindPagingResult[*model.ExcelFile] {
	return f.repos.FindPaging(ctx, qry)
}
