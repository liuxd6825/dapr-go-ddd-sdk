package service

import (
	"context"
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

type FolderMetaService struct {
	dao *dao.FolderMetaDao
}

func NewFolderMetaService() *FolderMetaService {
	return &FolderMetaService{
		dao: dao.NewFolderMetaDao(DBKey),
	}
}

func (t *FolderMetaService) Create(ctx context.Context, data *model.FolderMeta) error {
	return t.dao.Create(ctx, data).GetError()
}

func (t *FolderMetaService) CreateMany(ctx context.Context, entities []*model.FolderMeta) error {
	return t.dao.CreateMany(ctx, entities).GetError()
}

func (t *FolderMetaService) DeleteByFolderId(ctx context.Context, folderId string) *idao.Result {
	return t.dao.DeleteByRSQL(ctx, fmt.Sprintf("folder_id=='%s'", folderId))
}

func (t *FolderMetaService) Update(ctx context.Context, data *model.FolderMeta, opts ...idao.CallOptions) error {
	return t.dao.Update(ctx, data, opts...).GetError()
}

func (t *FolderMetaService) Submit(ctx context.Context, data *model.FolderMeta, opts ...idao.CallOptions) error {
	metas, err := t.dao.FindByRSQL(ctx, fmt.Sprintf("folder_id=='%s' and source=='%s' and name=='%s'", data.FolderId, data.Source, data.Name))
	if err != nil {
		return err
	}
	if metas != nil && len(metas) == 1 {
		return t.dao.Update(ctx, data, opts...).GetError()
	} else {
		return t.dao.Create(ctx, data, opts...).GetError()
	}
}

func (t *FolderMetaService) FindBySource(ctx context.Context, sourceId string) ([]*model.FolderMeta, error) {
	return t.dao.FindByRSQL(ctx, fmt.Sprintf("source=='%s'", sourceId))
}

func (t *FolderMetaService) FindByFolderId(ctx context.Context, folderId string) ([]*model.FolderMeta, error) {
	return t.dao.FindByRSQL(ctx, fmt.Sprintf("folder_id=='%s'", folderId))
}
