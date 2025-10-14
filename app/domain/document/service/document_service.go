package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

type DocumentService struct {
	dao *dao.DocumentDao
}

func NewDocumentService() *DocumentService {
	return &DocumentService{
		dao: dao.NewDocumentDao(DBKey),
	}
}

func (s *DocumentService) GetConfig() *idao.DaoConfig {
	return s.dao.GetConfig()
}

func (s *DocumentService) HasDocumentByDocument(ctx context.Context, folderId string) (bool, error) {
	count, err := s.dao.CountByRSQL(ctx, "folder_id=='"+folderId+"'")
	return count > 0, err
}

func (t *DocumentService) Create(ctx context.Context, cmd *command.DocumentCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Create(ctx, &cmd.Data).GetError()
	})
}

func (t *DocumentService) Delete(ctx context.Context, cmd *command.DocumentDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.DeleteById(ctx, cmd.Data.Id).GetError()
	})
}

func (t *DocumentService) DeleteById(ctx context.Context, id string) *idao.Result {
	return t.dao.DeleteById(ctx, id)
}

func (t *DocumentService) DeleteByRSQL(ctx context.Context, rsql string) *idao.Result {
	return t.dao.DeleteByRSQL(ctx, rsql)
}

func (t *DocumentService) Update(ctx context.Context, data *model.Document, opts ...idao.CallOptions) error {
	return t.dao.Update(ctx, data, opts...).GetError()
}

func (t *DocumentService) UpdateMany(ctx context.Context, entities []*model.Document, opts ...idao.CallOptions) error {
	return t.dao.UpdateMany(ctx, entities, opts...).GetError()
}

func (t *DocumentService) FindById(ctx context.Context, id string) (*model.Document, error) {
	return t.dao.FindById(ctx, id)
}

func (t *DocumentService) FindPaging(ctx context.Context, qry store.FindPagingQuery) (idao.FindPagingResult[*model.Document], error) {
	res := t.dao.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (t *DocumentService) FindByRSQL(ctx context.Context, rsql string) ([]*model.Document, error) {
	return t.dao.FindByRSQL(ctx, rsql)
}
