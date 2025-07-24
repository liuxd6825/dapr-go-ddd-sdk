package restapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/service"
	tagModel "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	tagSvc "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"io"
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"

	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type DocumentAPI struct {
	env             *env.Env
	documentService *service.DocumentService
	fileService     *service.FileService
	fsService       *service.FsService
	folderService   *service.FolderService
	tagRelationSvc  *tagSvc.TagRelationService
}

func NewDocumentAPI(env *env.Env, rootPath string) *DocumentAPI {
	documentService := service.NewDocumentService()
	fileService := service.NewFileService()
	fsService := service.NewFsService()
	folderService := service.NewFolderService()
	tagRelationSvc := tagSvc.NewTagRelationService()
	return &DocumentAPI{
		env:             env,
		documentService: documentService,
		fileService:     fileService,
		fsService:       fsService,
		folderService:   folderService,
		tagRelationSvc:  tagRelationSvc,
	}
}

func (s *DocumentAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodPost, "/doc/document:upload-chunk", "UploadChunk")
	b.Handle(iris.MethodGet, "/doc/document:download", "Download")
	b.Handle(iris.MethodPost, "/doc/document", "Create")
	b.Handle(iris.MethodPut, "/doc/document", "Update")
	b.Handle(iris.MethodPut, "/doc/document:rename", "Rename")
	b.Handle(iris.MethodPut, "/doc/document:move", "Move")
	b.Handle(iris.MethodPut, "/doc/document:update-tags", "UpdateTags")
	b.Handle(iris.MethodDelete, "/doc/document", "Delete")
	b.Handle(iris.MethodGet, "/doc/document", "FindPaging")
	b.Handle(iris.MethodGet, "/doc/document/:id", "FindById")
	b.Handle(iris.MethodGet, "/doc/document:by-tag", "FindByTagId")
}

func (s *DocumentAPI) UploadChunk(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		chunk, _, err := ictx.FormFile("chunk")
		chunkIndex := ictx.FormValue("chunkIndex")
		chunkSize := ictx.FormValue("chunkSize")
		objectName := ictx.FormValue("objectName")
		folderPath := ictx.FormValue("folderPath")
		if err != nil {
			return err
		}

		has := s.fsService.Exists(folderPath + "/" + objectName)

		if chunkIndex == "0" && !has {
			file := s.fsService.Create(folderPath + "/" + objectName)
			if file != nil {
				file.Close()
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
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *DocumentAPI) Download(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		fileId := ictx.URLParam("file-id")
		if fileId == "" {
			return errors.New("Url参数file-id不能为空")
		}

		file, err := s.fileService.FindById(ctx, fileId)
		if err != nil {
			return err
		}
		if file == nil {
			return errors.New("没有找到文件记录,fileId=" + fileId)
		}

		folder, err := s.folderService.FindById(ctx, file.FolderId)
		if err != nil {
			return err
		}
		if folder == nil {
			return errors.New("没有找到文件目录,fileId=" + fileId)
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
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *DocumentAPI) Create(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		err := tx.StartTx(ctx, []string{s.documentService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
			var cmd *command.DocumentCreateCommand
			if err := ictx.ReadJSON(&cmd); err != nil {
				return err
			}
			res := s.fileService.Create(ctx, s.fileService.GetFile(&cmd.Data))
			if res.Error != nil {
				return res.Error
			}
			res = s.documentService.Create(ctx, &cmd.Data)
			if res.Error != nil {
				return res.Error
			}
			return nil
		})
		return err
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *DocumentAPI) Rename(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		err := tx.StartTx(ctx, []string{s.documentService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
			var cmd *command.DocumentRenameCommand
			if err := ictx.ReadJSON(&cmd); err != nil {
				return err
			}
			opts := idao.NewCallOptions()
			opts.SetUpdateFields([]string{"name", "object_name"})

			docModel := model.Document{}
			docModel.Id = cmd.Data.Id
			docModel.Name = cmd.Data.Name
			docModel.ObjectName = cmd.Data.ObjectName
			res := s.documentService.Update(ctx, &docModel, opts)
			if res.Error != nil {
				return res.Error
			}

			fileModel := model.File{}
			fileModel.Id = cmd.Data.FileId
			fileModel.Name = cmd.Data.Name
			fileModel.ObjectName = cmd.Data.ObjectName
			res = s.fileService.Update(ctx, &fileModel, opts)
			if res.Error != nil {
				return res.Error
			}

			err := s.fsService.Rename(cmd.Data.FolderPath+"/"+cmd.Data.OldName, cmd.Data.FolderPath+"/"+cmd.Data.ObjectName)
			if err != nil {
				return err
			}

			return nil
		})
		return err
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *DocumentAPI) Move(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		err := tx.StartTx(ctx, []string{s.documentService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
			var cmd *command.DocumentMoveCommand
			if err := ictx.ReadJSON(&cmd); err != nil {
				return err
			}

			files, err := s.fileService.FindByRSQL(ctx, fmt.Sprintf("document_id=='%s'", cmd.Data.Id))
			if err != nil {
				return err
			}
			for _, file := range files {
				file.FolderId = cmd.Data.FolderId
			}

			opts := idao.NewCallOptions()
			opts.SetUpdateFields([]string{"folder_id"})

			var res *idao.Result

			if len(files) > 0 {
				res = s.fileService.UpdateMany(ctx, files, opts)
				if res.Error != nil {
					return res.Error
				}
			}

			doc := model.Document{}
			doc.Id = cmd.Data.Id
			doc.FolderId = cmd.Data.FolderId

			res = s.documentService.Update(ctx, &doc, opts)
			if res.Error != nil {
				return res.Error
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
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *DocumentAPI) UpdateTags(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		err := tx.StartTx(ctx, []string{s.documentService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
			var cmd *command.DocumentUpdateCommand
			if err := ictx.ReadJSON(&cmd); err != nil {
				return err
			}

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

			res = s.documentService.Update(ctx, &cmd.Data, opts)
			if res.Error != nil {
				return res.Error
			}

			return nil
		})
		return err
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *DocumentAPI) Update(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		err := tx.StartTx(ctx, []string{s.documentService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
			var cmd *command.DocumentUpdateCommand
			if err := ictx.ReadJSON(&cmd); err != nil {
				return err
			}
			files, err := s.fileService.FindByRSQL(ctx, "document_id=='"+cmd.Data.Id+"'")
			if err != nil {
				return err
			}
			for _, file := range files {
				file.IsMain = false
			}
			var res *idao.Result
			if len(files) > 0 {
				res = s.fileService.UpdateMany(ctx, files)
				if res.Error != nil {
					return res.Error
				}
			}

			res = s.fileService.Create(ctx, s.fileService.GetFile(&cmd.Data))
			if res.Error != nil {
				return res.Error
			}

			opts := idao.NewCallOptions()
			opts.SetUpdateFields([]string{"file_id", "name", "object_name", "ext_name", "size", "size_title", "download_total", "download_url", "preview_url", "thumbnail", "md5"})
			res = s.documentService.Update(ctx, &cmd.Data, opts)
			if res.Error != nil {
				return res.Error
			}
			return nil
		})
		return err
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *DocumentAPI) Delete(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		err := tx.StartTx(ctx, []string{s.documentService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
			var cmd *command.DocumentDeleteCommand
			if err := ictx.ReadJSON(&cmd); err != nil {
				return err
			}

			res := s.tagRelationSvc.DeleteByRSQL(ctx, fmt.Sprintf("bus_id=='%s'", cmd.Data.Id))
			if res.Error != nil {
				return res.Error
			}

			res = s.documentService.DeleteById(ctx, cmd.Data.Id)
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
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *DocumentAPI) FindPaging(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		folderId := ictx.URLParam("folder-id")
		//user, _ := appctx.GetAuthUser(ctx)
		qry := store.NewFindPagingQueryRequest()
		qry.PageNum = 0
		qry.PageSize = 99999999999999
		qry.Filter = "folder_id=='" + folderId + "'"
		qry.Sort = "created_time:desc"
		qry.IsTotalRows = true
		res := s.documentService.FindPaging(ctx, qry)
		return restapi.SetData(ictx, res)
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *DocumentAPI) FindById(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		id := ictx.Params().GetString("id")
		res, err := s.documentService.FindById(ctx, id)
		if err != nil {
			return err
		}
		return restapi.SetData(ictx, res)
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *DocumentAPI) FindByTagId(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		tagId := ictx.URLParam("tag-id")
		caseId := ictx.URLParam("case-id")
		//user, _ := appctx.GetAuthUser(ctx)
		qry := store.NewFindPagingQueryRequest()
		qry.PageNum = 0
		qry.PageSize = 99999999999999
		qry.Filter = fmt.Sprintf("case_id=='%s' and tag_id=contains='%s'", caseId, tagId)
		qry.Sort = "created_time:desc"
		qry.IsTotalRows = true
		res := s.documentService.FindPaging(ctx, qry)
		return restapi.SetData(ictx, res)
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}
