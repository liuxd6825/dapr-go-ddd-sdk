package service

import (
	"context"
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

type DocumentService struct {
	dao         *dao.DocumentDao
	metaService *DocumentMetaService
}

func NewDocumentService() *DocumentService {
	return &DocumentService{
		dao: dao.NewDocumentDao(DBKey),
	}
}

func (s *DocumentService) SetInit(metaService *DocumentMetaService) {
	s.metaService = metaService
}

func (s *DocumentService) GetConfig() *idao.DaoConfig {
	return s.dao.GetConfig()
}

func (s *DocumentService) HasDocumentByDocument(ctx context.Context, folderId string) (bool, error) {
	count, err := s.dao.CountByRSQL(ctx, "folder_id=='"+folderId+"'")
	return count > 0, err
}

//func (t *DocumentService) Create(ctx context.Context, cmd *command.DocumentCreateCommand) error {
//	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
//		return t.dao.Create(ctx, &cmd.Data).GetError()
//	})
//}

func (t *DocumentService) CreateData(ctx context.Context, data *model.Document) error {
	return t.dao.Create(ctx, data).GetError()
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

func (t *DocumentService) FindByFolderId(ctx context.Context, folderId string) ([]*model.Document, error) {
	return t.dao.FindByRSQL(ctx, fmt.Sprintf("folder_id=='%s'", folderId))
}

func (t *DocumentService) FindDocumentCountByFolderId(ctx context.Context, folderId string) (int64, error) {
	return t.dao.CountByRSQL(ctx, fmt.Sprintf("folder_id=='%s'", folderId))
}

func (t *DocumentService) UpdateDocumentMetas(ctx context.Context, folderId string, folderMeta *model.FolderMeta) error {
	docs, err := t.FindByFolderId(ctx, folderId)
	if err != nil {
		return err
	}
	if docs == nil {
		return nil
	}
	dMetas := &[]string{}
	uMetas := &[]*model.DocumentMeta{}
	cMetas := &[]*model.DocumentMeta{}

	for _, doc := range docs {
		err = t.metaService.BuildUpdateModels(ctx, doc.Id, folderMeta, dMetas, uMetas, cMetas)
		if err != nil {
			return err
		}
	}
	if len(*dMetas) > 0 {
		err = t.metaService.DeleteByIds(ctx, *dMetas)
		if err != nil {
			return err
		}
	}
	if len(*cMetas) > 0 {
		err = t.metaService.CreateMany(ctx, *cMetas)
		if err != nil {
			return err
		}
	}
	if len(*uMetas) > 0 {
		err = t.metaService.UpdateMany(ctx, *uMetas)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *DocumentService) Document2DocumentView(document *model.Document) *model.DocumentView {
	view := &model.DocumentView{}
	view.Id = document.Id
	view.TenantId = document.TenantId
	view.CaseId = document.CaseId
	view.CreatedTime = document.CreatedTime
	view.CreatorId = document.CreatorId
	view.CreatorName = document.CreatorName
	view.UpdatedTime = document.UpdatedTime
	view.UpdaterId = document.UpdaterId
	view.UpdaterName = document.UpdaterName
	/*	view.DeletedTime = document.DeletedTime
		view.DeleterId = document.DeleterId
		view.DeleterName = document.DeleterName
		view.IsDeleted = document.IsDeleted*/
	view.Remark = document.Remark
	view.BusId = document.BusId
	view.EntityId = document.EntityId
	view.RootId = document.RootId
	view.RootPath = document.RootPath
	view.FolderId = document.FolderId
	view.FileId = document.FileId
	view.ObjectName = document.ObjectName
	view.ExtName = document.ExtName
	view.Size = document.Size
	view.SizeTitle = document.SizeTitle
	view.DownloadTotal = document.DownloadTotal
	view.DownloadUrl = document.DownloadUrl
	view.PreviewUrl = document.PreviewUrl
	view.Thumbnail = document.Thumbnail
	view.Md5 = document.Md5
	view.TagName = document.TagName
	view.TagId = document.TagId
	view.TagColor = document.TagColor
	view.FsKey = document.FsKey
	view.Name = document.Name
	view.DisabledFrontEdit = document.DisabledFrontEdit

	return view
}

func (t *DocumentService) DocumentView2Document(view *model.DocumentView) *model.Document {
	document := &model.Document{}

	document.Id = view.Id
	document.TenantId = view.TenantId
	document.CaseId = view.CaseId
	document.CreatedTime = view.CreatedTime
	document.CreatorId = view.CreatorId
	document.CreatorName = view.CreatorName
	document.UpdatedTime = view.UpdatedTime
	document.UpdaterId = view.UpdaterId
	document.UpdaterName = view.UpdaterName
	/*	document.DeletedTime = view.DeletedTime
		document.DeleterId = view.DeleterId
		document.DeleterName = view.DeleterName
		document.IsDeleted = view.IsDeleted*/
	document.Remark = view.Remark
	document.BusId = view.BusId
	document.EntityId = view.EntityId
	document.RootId = view.RootId
	document.RootPath = view.RootPath
	document.FolderId = view.FolderId
	document.FileId = view.FileId
	document.ObjectName = view.ObjectName
	document.ExtName = view.ExtName
	document.Size = view.Size
	document.SizeTitle = view.SizeTitle
	document.DownloadTotal = view.DownloadTotal
	document.DownloadUrl = view.DownloadUrl
	document.PreviewUrl = view.PreviewUrl
	document.Thumbnail = view.Thumbnail
	document.Md5 = view.Md5
	document.TagName = view.TagName
	document.TagId = view.TagId
	document.TagColor = view.TagColor
	document.FsKey = view.FsKey
	document.Name = view.Name
	document.DisabledFrontEdit = view.DisabledFrontEdit

	return document
}
