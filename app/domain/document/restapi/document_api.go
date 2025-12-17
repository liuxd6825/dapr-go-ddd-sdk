package restapi

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/service"
	tagModel "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	tagSvc "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type DocumentAPI struct {
	rootPath            string
	env                 *env.Env
	documentService     *service.DocumentService
	documentMetaService *service.DocumentMetaService
	fileService         *service.FileService
	fsService           *service.FsService
	folderService       *service.FolderService
	folderMetaService   *service.FolderMetaService
	tagRelationSvc      *tagSvc.TagRelationService
}

func NewDocumentAPI(env *env.Env, rootPath string) *DocumentAPI {
	documentService := service.NewDocumentService()
	documentMetaService := service.NewDocumentMetaService()
	fileService := service.NewFileService()
	fsService := service.NewFsService()
	folderService := service.NewFolderService()
	folderMetaService := service.NewFolderMetaService()
	tagRelationSvc := tagSvc.NewTagRelationService()
	return &DocumentAPI{
		rootPath:            rootPath,
		env:                 env,
		documentService:     documentService,
		fileService:         fileService,
		fsService:           fsService,
		folderService:       folderService,
		folderMetaService:   folderMetaService,
		tagRelationSvc:      tagRelationSvc,
		documentMetaService: documentMetaService,
	}
}

func (s *DocumentAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath+"/doc", "document.DocumentAPI", s)
	ctl.Post("/document:upload-chunk", "UploadChunk")
	ctl.GetData("/document:download", "Download")
	ctl.Post("/document", "Create")
	ctl.Put("/document", "Update")
	ctl.Put("/document:rename", "Rename")
	ctl.Put("/document:move", "Move")
	ctl.Put("/document:update-tags", "UpdateTags")
	ctl.Delete("/document", "Delete", restapi.WithParamsInBody(true))
	ctl.GetData("/document", "FindPaging")
	ctl.GetOne("/document/:id", "FindById")
	ctl.GetData("/document:by-tag", "FindByTagId")
	return ctl
}

func (s *DocumentAPI) UploadChunk(ctx context.Context, iCtx iris.Context) error {
	chunk, _, err := iCtx.FormFile("chunk")
	chunkIndex := iCtx.FormValue("chunkIndex")
	chunkSize := iCtx.FormValue("chunkSize")
	objectName := iCtx.FormValue("objectName")
	folderPath := iCtx.FormValue("folderPath")
	folderId := iCtx.FormValue("folderId")
	if err != nil {
		return err
	}

	count, err := s.folderService.GetChildrenCount(ctx, folderId)
	if err != nil {
		return err
	}
	if count > 1024 {
		return errors.New("文件夹与文档总数量不能超过1024")
	}

	has := s.fsService.Exists(folderPath + "/" + objectName)

	if chunkIndex == "0" && !has {
		file := s.fsService.Create(folderPath + "/" + objectName)
		if file != nil {
			_ = file.Close()
		}
	}

	data, err := io.ReadAll(chunk)
	if err != nil {
		return err
	}

	err = s.fsService.WriteAt(folderPath+"/"+objectName, data, chunkIndex, chunkSize)
	if err != nil {
		return err
	}

	return nil
}

func (s *DocumentAPI) Download(ctx context.Context, ictx iris.Context, query *query.DownloadParamQuery) error {
	file, err := s.fileService.FindById(ctx, query.FileId)
	if err != nil {
		return err
	}
	if file == nil {
		return errors.New("没有找到文件记录,fileId=" + query.FileId)
	}

	folder, err := s.folderService.FindById(ctx, file.FolderId)
	if err != nil {
		return err
	}
	if folder == nil {
		return errors.New("没有找到文件目录,fileId=" + query.FileId)
	}

	has := s.fsService.Exists(folder.FolderPath + "/" + file.ObjectName)
	if !has {
		return errors.New("没有找到文件,objectName=" + folder.FolderPath + "/" + file.ObjectName)
	}

	err = s.fsService.Download(ictx, folder.FolderPath+"/"+file.ObjectName, file.Name)
	if err != nil {
		return err
	}

	return nil
}

func (s *DocumentAPI) Create(ctx context.Context, cmd *command.DocumentCreateCommand) error {
	err := tx.StartTx(ctx, []string{s.documentService.GetConfig().DBKey}, func(ctx context.Context, options ...*store2.SessionOptions) error {

		document := s.documentService.DocumentView2Document(&cmd.Data)

		err := s.fileService.Create(ctx, s.fileService.GetFile(document))
		if err != nil {
			return err
		}

		err = s.documentService.CreateData(ctx, document)
		if err != nil {
			return err
		}

		meta := cmd.Data.Meta
		if meta != nil && len(meta) > 0 {
			err = s.documentMetaService.CreateMany(ctx, meta)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (s *DocumentAPI) Rename(ctx context.Context, cmd *command.DocumentRenameCommand) error {
	err := tx.StartTx(ctx, []string{s.documentService.GetConfig().DBKey}, func(ctx context.Context, options ...*store2.SessionOptions) error {
		opts := idao.NewCallOptions()
		opts.SetUpdateFields([]string{"name", "object_name"})

		docModel := model.Document{}
		docModel.Id = cmd.Data.Id
		docModel.Name = cmd.Data.Name
		docModel.ObjectName = cmd.Data.ObjectName
		err := s.documentService.Update(ctx, &docModel, opts)
		if err != nil {
			return err
		}

		fileModel := model.File{}
		fileModel.Id = cmd.Data.FileId
		fileModel.Name = cmd.Data.Name
		fileModel.ObjectName = cmd.Data.ObjectName
		err = s.fileService.Update(ctx, &fileModel, opts)
		if err != nil {
			return err
		}

		err = s.fsService.Rename(cmd.Data.FolderPath+"/"+cmd.Data.OldName, cmd.Data.FolderPath+"/"+cmd.Data.ObjectName)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func (s *DocumentAPI) Move(ctx context.Context, cmd *command.DocumentMoveCommand) error {
	err := tx.StartTx(ctx, []string{s.documentService.GetConfig().DBKey}, func(ctx context.Context, options ...*store2.SessionOptions) error {
		files, err := s.fileService.FindByRSQL(ctx, fmt.Sprintf("document_id=='%s'", cmd.Data.Id))
		if err != nil {
			return err
		}
		for _, file := range files {
			file.FolderId = cmd.Data.FolderId
		}

		opts := idao.NewCallOptions()
		opts.SetUpdateFields([]string{"folder_id"})

		if len(files) > 0 {
			err = s.fileService.UpdateMany(ctx, files, opts)
			if err != nil {
				return err
			}
		}

		doc := model.Document{}
		doc.Id = cmd.Data.Id
		doc.FolderId = cmd.Data.FolderId

		err = s.documentService.Update(ctx, &doc, opts)
		if err != nil {
			return err
		}

		folderMetas, err := s.folderMetaService.FindByFolderId(ctx, doc.FolderId)
		if err != nil {
			return err
		}

		folderMeta := s.folderMetaService.GetSourceModel(folderMetas)

		dMetas := &[]string{}
		uMetas := &[]*model.DocumentMeta{}
		cMetas := &[]*model.DocumentMeta{}

		err = s.documentMetaService.BuildUpdateModels(ctx, doc.Id, folderMeta, dMetas, uMetas, cMetas)

		if len(*dMetas) > 0 {
			err = s.documentMetaService.DeleteByIds(ctx, *dMetas)
			if err != nil {
				return err
			}
		}

		if len(*cMetas) > 0 {
			err = s.documentMetaService.CreateMany(ctx, *cMetas)
			if err != nil {
				return err
			}
		}

		if len(*uMetas) > 0 {
			err = s.documentMetaService.UpdateMany(ctx, *uMetas)
			if err != nil {
				return err
			}
		}

		//处理文件移动 非主版本文件也需要移动
		for _, file := range files {
			err = s.fsService.Rename(cmd.Data.SourcePath+"/"+file.ObjectName, cmd.Data.TargetPath+"/"+file.ObjectName)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return err
}

func (s *DocumentAPI) UpdateTags(ctx context.Context, cmd *command.DocumentUpdateCommand) error {
	err := tx.StartTx(ctx, []string{s.documentService.GetConfig().DBKey}, func(ctx context.Context, options ...*store2.SessionOptions) error {
		res := s.tagRelationSvc.DeleteByRSQL(ctx, fmt.Sprintf("bus_id=='%s'", cmd.Data.Id))
		if res.Error != nil {
			return res.Error
		}
		if cmd.Data.TagId != "" {
			arrTagId := strings.Split(cmd.Data.TagId, ",")
			arrTagName := strings.Split(cmd.Data.TagName, ",")
			arrTagColor := strings.Split(cmd.Data.TagColor, ",")
			trs := make([]*tagModel.TagRelation, 0)
			for index, tagId := range arrTagId {
				tr := &tagModel.TagRelation{}
				tr.Id = idutils.NewId()
				tr.TagId = tagId
				tr.TagName = arrTagName[index]
				tr.TagColor = arrTagColor[index]
				tr.CaseId = cmd.Data.CaseId
				tr.BusId = cmd.Data.Id
				tr.BusType = "文档中心"
				tr.ChangedSource = "BusinessSystem"
				trs = append(trs, tr)
			}
			res = s.tagRelationSvc.CreateMany(ctx, trs)
			if res.Error != nil {
				return res.Error
			}
		}

		opts := idao.NewCallOptions()
		opts.SetUpdateFields([]string{"tag_id", "tag_name", "tag_color"})

		err := s.documentService.Update(ctx, &cmd.Data, opts)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *DocumentAPI) Update(ctx context.Context, cmd *command.DocumentUpdateCommand) error {
	err := tx.StartTx(ctx, []string{s.documentService.GetConfig().DBKey}, func(ctx context.Context, options ...*store2.SessionOptions) error {
		files, err := s.fileService.FindByRSQL(ctx, "document_id=='"+cmd.Data.Id+"'")
		if err != nil {
			return err
		}
		for _, file := range files {
			file.IsMain = false
		}
		if len(files) > 0 {
			err = s.fileService.UpdateMany(ctx, files)
			if err != nil {
				return err
			}
		}

		err = s.fileService.Create(ctx, s.fileService.GetFile(&cmd.Data))
		if err != nil {
			return err
		}

		opts := idao.NewCallOptions()
		opts.SetUpdateFields([]string{"file_id", "name", "object_name", "ext_name", "size", "size_title", "download_total", "download_url", "preview_url", "thumbnail", "md5"})
		err = s.documentService.Update(ctx, &cmd.Data, opts)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *DocumentAPI) Delete(ctx context.Context, cmd *command.DocumentDeleteCommand) error {
	err := tx.StartTx(ctx, []string{s.documentService.GetConfig().DBKey}, func(ctx context.Context, options ...*store2.SessionOptions) error {
		//删标签关系
		res := s.tagRelationSvc.DeleteByRSQL(ctx, fmt.Sprintf("bus_id=='%s'", cmd.Data.Id))
		if res.Error != nil {
			return res.Error
		}

		res = s.documentService.DeleteById(ctx, cmd.Data.Id)
		if res.Error != nil {
			return res.Error
		}

		res = s.documentMetaService.DeleteByDocumentId(ctx, cmd.Data.Id)
		if res.Error != nil {
			return res.Error
		}

		res = s.fileService.DeleteByRSQL(ctx, fmt.Sprintf("document_id=='%s'", cmd.Data.Id))
		if res.Error != nil {
			return res.Error
		}

		s.fsService.RemoveFile(fmt.Sprintf("/%s/%s", cmd.Data.FolderPath, cmd.Data.ObjectName))
		return nil
	})
	return err
}

func (s *DocumentAPI) FindPaging(ctx context.Context, query *query.FindByFolderAndFilter) (idao.FindPagingResult[*model.DocumentView], error) {
	qry := store2.NewFindPagingQueryRequest()
	qry.PageNum = 0
	qry.PageSize = 99999999999999
	qry.Filter = "folder_id=='" + query.FolderId + "'"
	if len(query.Filter) > 0 {
		qry.Filter += " and " + query.Filter
	}
	qry.Sort = "created_time:desc"
	qry.IsTotalRows = true

	documents, err := s.documentService.FindPaging(ctx, qry)
	if err != nil {
		return nil, err
	}
	if len(documents.GetData()) == 0 {
		return store2.NewFindPagingResult([]*model.DocumentView{}, 0, qry, nil), nil
	}
	docIds := make([]string, 0)
	for _, document := range documents.GetData() {
		docIds = append(docIds, document.Id)
	}
	allMetas, err := s.documentMetaService.FindByDocumentIds(ctx, docIds)
	if err != nil {
		return nil, err
	}

	dvs := make([]*model.DocumentView, 0)
	for _, d := range documents.GetData() {
		dv := s.documentService.Document2DocumentView(d)
		dv.Meta = s.getMeta(allMetas, dv.Id)
		dvs = append(dvs, dv)
	}

	res := store2.NewFindPagingResult(dvs, int64(len(dvs)), qry, nil)

	return res, nil
}

func (s *DocumentAPI) getMeta(metas []*model.DocumentMeta, documentId string) []*model.DocumentMeta {
	docMetas := make([]*model.DocumentMeta, 0)
	if metas == nil || len(metas) == 0 {
		return []*model.DocumentMeta{}
	}
	for _, meta := range metas {
		if meta.DocumentId == documentId {
			docMetas = append(docMetas, meta)
		}
	}
	return docMetas
}

func (s *DocumentAPI) getMetaByDocumentId(ctx context.Context, documentId string) ([]*model.DocumentMeta, error) {
	return s.documentMetaService.FindByDocumentId(ctx, documentId)
}

func (s *DocumentAPI) FindById(ctx context.Context, query *query.FindByIdQuery) (*model.DocumentView, error) {
	data, err := s.documentService.FindById(ctx, query.Id)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	meta, err := s.documentMetaService.FindByDocumentId(ctx, data.Id)
	if err != nil {
		return nil, err
	}

	view := s.documentService.Document2DocumentView(data)
	view.Meta = meta

	return view, nil
}

func (s *DocumentAPI) FindByTagId(ctx context.Context, query *query.FindByTagAndCase) (idao.FindPagingResult[*model.Document], error) {
	qry := store2.NewFindPagingQueryRequest()
	qry.PageNum = 0
	qry.PageSize = 99999999999999
	qry.Filter = fmt.Sprintf("case_id=='%s' and tag_id=contains='%s'", query.CaseId, query.TagId)
	qry.Sort = "created_time:desc"
	qry.IsTotalRows = true
	return s.documentService.FindPaging(ctx, qry)
}
