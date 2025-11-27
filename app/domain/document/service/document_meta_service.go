package service

import (
	"context"
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

type DocumentMetaService struct {
	dao *dao.DocumentMetaDao
}

func NewDocumentMetaService() *DocumentMetaService {
	return &DocumentMetaService{
		dao: dao.NewDocumentMetaDao(DBKey),
	}
}

func (t *DocumentMetaService) Create(ctx context.Context, data *model.DocumentMeta) error {
	return t.dao.Create(ctx, data).GetError()
}

func (t *DocumentMetaService) CreateMany(ctx context.Context, entities []*model.DocumentMeta) error {
	return t.dao.CreateMany(ctx, entities).GetError()
}

func (t *DocumentMetaService) Update(ctx context.Context, data *model.DocumentMeta, opts ...idao.CallOptions) error {
	return t.dao.Update(ctx, data, opts...).GetError()
}

func (t *DocumentMetaService) DeleteByDocumentId(ctx context.Context, documentId string) *idao.Result {
	return t.dao.DeleteByRSQL(ctx, fmt.Sprintf("document_id=='%s'", documentId))
}

func (t *DocumentMetaService) Submit(ctx context.Context, data *model.DocumentMeta, opts ...idao.CallOptions) error {
	metas, err := t.dao.FindByRSQL(ctx, fmt.Sprintf("document_id=='%s' and source=='%s' and name=='%s'", data.DocumentId, data.Source, data.Name))
	if err != nil {
		return err
	}
	if metas != nil && len(metas) == 1 {
		return t.dao.Update(ctx, data, opts...).GetError()
	} else {
		return t.dao.Create(ctx, data, opts...).GetError()
	}
}

func (t *DocumentMetaService) FindBySource(ctx context.Context, sourceId string) ([]*model.DocumentMeta, error) {
	return t.dao.FindByRSQL(ctx, fmt.Sprintf("source=='%s'", sourceId))
}

func (t *DocumentMetaService) FindByDocumentId(ctx context.Context, folderId string) ([]*model.DocumentMeta, error) {
	return t.dao.FindByRSQL(ctx, fmt.Sprintf("document_id=='%s'", folderId))
}
