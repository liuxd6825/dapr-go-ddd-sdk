package service

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/event"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	dao2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/goserver/pkg/orm_pkg/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

type DocumentMetaService struct {
	dao           *dao.DocumentMetaDao
	docService    *DocumentService
	folderService *FolderService
}

func NewDocumentMetaService() *DocumentMetaService {
	svc := &DocumentMetaService{
		dao: dao.NewDocumentMetaDao(DBKey),
	}
	return svc
}

func (t *DocumentMetaService) SetInit(docService *DocumentService, folderService *FolderService) {
	t.docService = docService
	t.folderService = folderService
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

func (t *DocumentMetaService) UpdateMany(ctx context.Context, entities []*model.DocumentMeta) error {
	return t.dao.UpdateMany(ctx, entities).GetError()
}

func (t *DocumentMetaService) DeleteByDocumentId(ctx context.Context, documentId string) *idao.Result {
	return t.dao.DeleteByRSQL(ctx, fmt.Sprintf("document_id=='%s'", documentId))
}

func (t *DocumentMetaService) DeleteByIds(ctx context.Context, ids []string) error {
	builder := dao2.NewRSQLBuilder()
	rSql := builder.In("id", ids).Build()
	return t.dao.DeleteByRSQL(ctx, rSql).GetError()
}

func (t *DocumentMetaService) Save(ctx context.Context, cmd *command.DocumentMetaSubmitCommand, opts ...idao.CallOptions) error {
	meta, err := t.FindByDocumentIdAndName(ctx, cmd.Data.DocumentId, cmd.Data.Name)
	if err != nil {
		return err
	}
	if meta != nil {
		meta.Value = cmd.Data.Value
		opts := idao.NewCallOptions()
		opts.SetUpdateFields([]string{"value"})
		err = t.Update(ctx, meta, opts)
	} else {
		err = t.Create(ctx, &cmd.Data)
	}
	return err
}

func (t *DocumentMetaService) PublishDocumentRagEvent(ctx context.Context, cmd *command.DocumentMetaSubmitCommand) error {
	if cmd.Data.Name == "docType" && cmd.Data.Value == "知识" {
		docModel, err := t.docService.FindById(ctx, cmd.Data.DocumentId)
		if err != nil {
			return err
		}
		folder, err := t.folderService.FindById(ctx, docModel.FolderId)
		data := &event.DocumentCreateEventData{}
		data.Id = idutils.NewId()
		data.CaseId = docModel.CaseId
		data.FsKey = docModel.FsKey
		data.FileId = docModel.FileId
		data.FileName = docModel.ObjectName
		data.FilePath = folder.FolderPath
		data.SourceId = docModel.Id
		data.SourceType = "document"
		data.SourceApp = "document_service"
		evt := event.NewDocumentCreateEvent(ctx, "duxm-master-cmd-service", data)
		err = xbase.PublishEvent(ctx, evt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *DocumentMetaService) SaveStatusBySourceType(ctx context.Context, cmd *command.DocumentMetaSaveStatusBySourceType) error {
	verr := errors.NewVerifyError()
	if cmd.Data.SourceType == "" {
		verr.AppendField("SourceType", "不能为空")
	}
	if cmd.Data.Name == "" {
		verr.AppendField("Name", "不能为空")
	}
	if cmd.Data.Value == "" {
		verr.AppendField("Value", "不能为空")
	}
	if verr.HasError() {
		return verr
	}

	meta, err := t.FindBySourceTypeAndName(ctx, cmd.Data.DocumentId, cmd.Data.SourceType, cmd.Data.Name)
	if err != nil {
		return err
	}
	if meta != nil {
		meta.Value = cmd.Data.Value
		opts := idao.NewCallOptions()
		opts.SetUpdateFields([]string{"value"})
		err = t.Update(ctx, meta, opts)
	} else {
		if cmd.Data.CaseId == "" {
			verr.AppendField("CaseId", "不能为空")
			return verr
		}
		err = t.Create(ctx, &cmd.Data)
	}
	return err
}

func (t *DocumentMetaService) FindBySource(ctx context.Context, sourceId string) ([]*model.DocumentMeta, error) {
	return t.dao.FindByRSQL(ctx, fmt.Sprintf("source=='%s'", sourceId))
}

func (t *DocumentMetaService) FindByDocumentId(ctx context.Context, documentId string) ([]*model.DocumentMeta, error) {
	return t.dao.FindByRSQL(ctx, fmt.Sprintf("document_id=='%s'", documentId))
}

func (t *DocumentMetaService) FindByDocumentIdAndName(ctx context.Context, documentId, name string) (*model.DocumentMeta, error) {
	arr, err := t.dao.FindByRSQL(ctx, fmt.Sprintf("document_id=='%s' and name=='%s'", documentId, name))
	if err != nil {
		return nil, err
	}
	if arr == nil || len(arr) == 0 {
		return nil, nil
	}
	return arr[0], nil
}

func (t *DocumentMetaService) FindBySourceTypeAndName(ctx context.Context, documentId, sourceType, name string) (*model.DocumentMeta, error) {
	arr, err := t.dao.FindByRSQL(ctx, fmt.Sprintf("document_id=='%s' and source_typoe=='%s' and name=='%s'", documentId, sourceType, name))
	if err != nil {
		return nil, err
	}
	if arr == nil || len(arr) == 0 {
		return nil, nil
	}
	return arr[0], nil
}

func (t *DocumentMetaService) FindByDocumentIds(ctx context.Context, documentIds []string) ([]*model.DocumentMeta, error) {
	builder := dao2.NewRSQLBuilder()
	rSql := builder.In("document_id", documentIds).Build()
	return t.dao.FindByRSQL(ctx, rSql)
}

func (t *DocumentMetaService) BuildUpdateModels(ctx context.Context, documentId string, targetMeta *model.FolderMeta, dMeta *[]string, uMeta *[]*model.DocumentMeta, cMeta *[]*model.DocumentMeta) error {
	metas, err := t.FindByDocumentId(ctx, documentId)
	if err != nil {
		return err
	}

	if metas == nil || len(metas) == 0 {
		if targetMeta != nil {
			newMeta := &model.DocumentMeta{}
			newMeta.Id = idutils.NewId()
			newMeta.CaseId = targetMeta.CaseId
			newMeta.DocumentId = documentId
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
					newMeta := &model.DocumentMeta{}
					newMeta.Id = meta.Id
					newMeta.CaseId = targetMeta.CaseId
					newMeta.DocumentId = meta.DocumentId
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
