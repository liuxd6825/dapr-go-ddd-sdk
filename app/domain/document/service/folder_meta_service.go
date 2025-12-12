package service

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	dao2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/goserver/pkg/orm_pkg/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
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

func (t *FolderMetaService) DeleteByRSql(ctx context.Context, rSql string) *idao.Result {
	return t.dao.DeleteByRSQL(ctx, rSql)
}

func (t *FolderMetaService) DeleteByIds(ctx context.Context, ids []string) error {
	builder := dao2.NewRSQLBuilder()
	rSql := builder.Out("id", ids).Build()
	return t.dao.DeleteByRSQL(ctx, rSql).GetError()
}

func (t *FolderMetaService) DeleteByFolderId(ctx context.Context, folderId string) *idao.Result {
	return t.dao.DeleteByRSQL(ctx, fmt.Sprintf("folder_id=='%s'", folderId))
}

func (t *FolderMetaService) Update(ctx context.Context, data *model.FolderMeta, opts ...idao.CallOptions) error {
	return t.dao.Update(ctx, data, opts...).GetError()
}

func (t *FolderMetaService) UpdateMany(ctx context.Context, entities []*model.FolderMeta, opts ...idao.CallOptions) error {
	return t.dao.UpdateMany(ctx, entities, opts...).GetError()
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

func (t *FolderMetaService) FindByFolderIds(ctx context.Context, folderIds []string) ([]*model.FolderMeta, error) {
	builder := dao2.NewRSQLBuilder()
	rSql := builder.In("folder_id", folderIds).Build()
	return t.dao.FindByRSQL(ctx, rSql)
}

func (t *FolderMetaService) FindByFolderId(ctx context.Context, folderId string) ([]*model.FolderMeta, error) {
	return t.dao.FindByRSQL(ctx, fmt.Sprintf("folder_id=='%s'", folderId))
}

func (t *FolderMetaService) BuildUpdateModels(ctx context.Context, folderId string, targetMeta *model.FolderMeta, dMeta *[]string, uMeta *[]*model.FolderMeta, cMeta *[]*model.FolderMeta) error {
	metas, err := t.FindByFolderId(ctx, folderId)
	if err != nil {
		return err
	}

	if metas == nil || len(metas) == 0 {
		if targetMeta != nil {
			newMeta := model.NewFolderMeta()
			newMeta.Id = idutils.NewId()
			newMeta.CaseId = targetMeta.CaseId
			newMeta.FolderId = folderId
			newMeta.SourceType = targetMeta.SourceType
			newMeta.Source = targetMeta.Source
			newMeta.Name = ""
			newMeta.Value = ""
			*cMeta = append(*cMeta, newMeta)
		}
	} else {
		for _, meta := range metas {
			if meta.Name == "" || meta.Name == "SourceId" {
				*dMeta = append(*dMeta, meta.Id)
				if targetMeta != nil {
					newMeta := model.NewFolderMeta()
					newMeta.Id = meta.Id
					newMeta.CaseId = targetMeta.CaseId
					newMeta.FolderId = meta.FolderId
					newMeta.SourceType = targetMeta.SourceType
					newMeta.Source = targetMeta.Source
					newMeta.Name = ""
					newMeta.Value = ""
					*cMeta = append(*cMeta, newMeta)
				}
			} else {
				if targetMeta != nil {
					meta.SourceType = targetMeta.SourceType
					meta.Source = targetMeta.Source
				} else {
					meta.SourceType = ""
					meta.Source = ""
				}
				*uMeta = append(*uMeta, meta)
			}
		}
	}

	return nil
}

func (t *FolderMetaService) GetSourceModel(folderMetas []*model.FolderMeta) *model.FolderMeta {
	if folderMetas == nil || len(folderMetas) == 0 {
		return nil
	}
	var source *model.FolderMeta
	for _, folderMeta := range folderMetas {
		if folderMeta.Name == "" || folderMeta.Name == "SourceId" {
			source = folderMeta
		}
	}
	return source
}
