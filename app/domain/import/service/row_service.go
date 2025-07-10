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
	repos *dao.RowDao
}

func NewRowService() *RowService {
	return singleutils.CreateObj[*RowService](func() *RowService {
		return &RowService{
			repos: dao.NewRowDao(config.DBKey),
		}
	})
}

func (f *RowService) Create(ctx context.Context, m *model.Row) {
	f.repos.Create(ctx, m)
}

func (f *RowService) CreateMany(ctx context.Context, m []*model.Row) {
	f.repos.CreateMany(ctx, m)
}

func (f *RowService) Update(ctx context.Context, m *model.Row) {
	f.repos.Update(ctx, m)
}

func (f *RowService) DeleteById(ctx context.Context, id string) {
	f.repos.DeleteById(ctx, id)
}

func (f *RowService) FindById(ctx context.Context, id string) (*model.Row, bool, error) {
	f.repos.FindById(ctx, id)
}

func (f *RowService) FindBySheet(ctx context.Context, qry *query.FindRowBySheetQuery) []*model.Row {
	b := idao.NewFindPagingQueryBuilder().SetPageNum(0).SetPageSize(1000)
	b.SetTenantId(qry.TenantId)
	b.SetMustFilter(fmt.Sprintf("case_id=='%v'", qry.CaseId))
	b.SetFilter(fmt.Sprintf("sheet_id=='%s'", qry.SheetId))
	b.SetSort("row_num:asc")
	res := f.FindPaging(ctx, b.Build())
	return res.Data
}

func (f *RowService) FindPaging(ctx context.Context, qry idao.FindPagingQuery) *idao.FindPagingResult[*model.Row] {
	return f.repos.FindPaging(ctx, qry)
}
