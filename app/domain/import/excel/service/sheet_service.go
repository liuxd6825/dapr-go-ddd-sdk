package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/excel/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/excel/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/recordie/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/singleutils"
)

type SheetService struct {
	repos idao.Dao[*model.ExcelSheet]
}

func NewSheetService() *SheetService {
	return singleutils.CreateObj[*SheetService](func() *SheetService {
		return &SheetService{
			repos: dao.NewSheetDao(config.DBKey),
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
