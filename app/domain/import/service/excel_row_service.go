package service

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/singleutils"
)

type RowService struct {
	rowDao *dao.ExcelRowDao
}

func NewRowService() *RowService {
	return singleutils.CreateObj[*RowService](func() *RowService {
		return &RowService{
			rowDao: dao.NewRowDao(config.DBKey),
		}
	})
}

func (f *RowService) Create(ctx context.Context, m *model.ExcelRow) {
	f.rowDao.Create(ctx, m)
}

func (f *RowService) CreateMany(ctx context.Context, m []*model.ExcelRow) error {
	return f.rowDao.CreateMany(ctx, m).GetError()
}

func (f *RowService) Update(ctx context.Context, m *model.ExcelRow) {
	f.rowDao.Update(ctx, m)
}

func (f *RowService) DeleteById(ctx context.Context, id string) {
	f.rowDao.DeleteById(ctx, id)
}

func (f *RowService) FindById(ctx context.Context, id string) (*model.ExcelRow, error) {
	return f.rowDao.FindById(ctx, id)
}

func (f *RowService) FindBySheet(ctx context.Context, qry *query.FindRowBySheetQuery) ([]*model.ExcelRow, error) {
	b := idao.NewFindPagingQueryBuilder().SetPageNum(0).SetPageSize(1000)
	b.SetTenantId(qry.TenantId)
	b.SetMustFilter(fmt.Sprintf("case_id=='%v'", qry.CaseId))
	b.SetFilter(fmt.Sprintf("sheet_id=='%s'", qry.SheetId))
	b.SetSort("row_num:asc")
	res := f.FindPaging(ctx, b.Build())
	return res.GetData(), nil
}

func (f *RowService) FindPaging(ctx context.Context, qry idao.FindPagingQuery) idao.FindPagingResult[*model.ExcelRow] {
	return f.rowDao.FindPaging(ctx, qry)
}
