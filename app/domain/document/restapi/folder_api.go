package restapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type FolderAPI struct {
	env               *env.Env
	folderService     *service.FolderService
	docService        *service.DocumentService
	docMetaService    *service.DocumentMetaService
	fileService       *service.FileService
	fsService         *service.FsService
	folderMetaService *service.FolderMetaService
	rootPath          string
}

func NewFolderAPI(env *env.Env, rootPath string) *FolderAPI {
	folderService := service.NewFolderService()
	docService := service.NewDocumentService()
	fileService := service.NewFileService()
	fsService := service.NewFsService()
	folderMetaService := service.NewFolderMetaService()
	docMetaService := service.NewDocumentMetaService()
	return &FolderAPI{
		rootPath:          rootPath,
		env:               env,
		folderService:     folderService,
		docService:        docService,
		fileService:       fileService,
		fsService:         fsService,
		folderMetaService: folderMetaService,
		docMetaService:    docMetaService,
	}
}

func (s *FolderAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath+"/doc", "document.FolderAPI", s)
	ctl.Post("/folder:root", "CreateRoot")
	ctl.Post("/folder", "Create")
	ctl.Put("/folder", "Update")
	ctl.Put("/folder:rename", "Rename")
	ctl.Put("/folder:color", "SetColor")
	ctl.Put("/folder:move", "Move")
	ctl.Delete("/folder", "Delete", restapi.WithParamsInBody(true))
	ctl.GetPaging("/folder", "FindPaging")
	ctl.GetOne("/folder/:id", "FindById")
	ctl.GetData("/folder:tree", "FindTree")
	return ctl
}

func (s *FolderAPI) CreateRoot(ctx context.Context, cmd *command.FolderCreateCommand) error {
	//err := tx.StartTx(ctx, []string{s.folderService.GetConfig().DBKey}, func(ctx context.Context, options ...*store2.SessionOptions) error {
	//
	//	vErr := errors.NewVerifyError()
	//	if cmd.Data.TenantId == "" {
	//		vErr.AppendField("tenantId", "不能为空", "租户Id")
	//	}
	//	if cmd.Data.BusId == "" {
	//		vErr.AppendField("busId", "不能为空", "业务Id")
	//	}
	//	if cmd.Data.EntityId == "" {
	//		vErr.AppendField("entityId", "不能为空", "实体Id")
	//	}
	//	if vErr.HasError() {
	//		return vErr
	//	}
	//
	//	cmd.Data.Id = cmd.Data.TenantId + "_" + cmd.Data.BusId + "_" + cmd.Data.EntityId
	//	cmd.Data.RootId = cmd.Data.TenantId + "_" + cmd.Data.BusId + "_" + cmd.Data.EntityId
	//	cmd.Data.RootPath = "/" + cmd.Data.TenantId + "/" + cmd.Data.BusId + "/" + cmd.Data.EntityId
	//	cmd.Data.FolderPath = "/" + cmd.Data.TenantId + "/" + cmd.Data.BusId + "/" + cmd.Data.EntityId
	//	cmd.Data.Name = "文件库"
	//	err1 := s.folderService.Create(ctx, cmd)
	//	if err1 != nil {
	//		return err1
	//	}
	//
	//	s.fsService.MkdirAll(cmd.Data.FolderPath)
	//	return nil
	//})
	//return err
	return nil
}

func (s *FolderAPI) Create(ctx context.Context, cmd *command.FolderCreateCommand) error {
	err := tx.StartTx(ctx, []string{s.folderService.GetConfig().DBKey}, func(ctx context.Context, options ...*store2.SessionOptions) error {

		vErr := errors.NewVerifyError()
		if cmd.Data.TenantId == "" {
			vErr.AppendField("tenantId", "不能为空", "租户Id")
		}
		if cmd.Data.BusId == "" {
			vErr.AppendField("busId", "不能为空", "业务Id")
		}
		if cmd.Data.EntityId == "" {
			vErr.AppendField("entityId", "不能为空", "实体Id")
		}
		if vErr.HasError() {
			return vErr
		}

		count, err1 := s.folderService.GetChildrenCount(ctx, cmd.Data.ParentId)
		if err1 != nil {
			return err1
		}
		if count > 1024 {
			return errors.New("文件夹与文档总数量不能超过1024")
		}

		folder := s.folderService.FolderView2Folder(&cmd.Data)

		err := s.folderService.CreateData(ctx, folder)
		if err != nil {
			return err
		}

		meta := cmd.Data.Meta
		if meta != nil && len(meta) > 0 {
			err = s.folderMetaService.CreateMany(ctx, meta)
			if err != nil {
				return err
			}
		}

		s.fsService.MkdirAll(cmd.Data.FolderPath)
		return nil
	})

	return err
}

func (s *FolderAPI) Rename(ctx context.Context, cmd *command.FolderRenameCommand) error {
	if cmd.Data.Alias == "" {
		return s.renameFolder(ctx, cmd)
	}
	return s.updateAlias(ctx, &cmd.Data)
}

func (s *FolderAPI) renameFolder(ctx context.Context, cmd *command.FolderRenameCommand) error {
	err := tx.StartTx(ctx, []string{s.folderService.GetConfig().DBKey}, func(ctx context.Context, options ...*store2.SessionOptions) error {
		tenantId := appctx.GetTenantId2(ctx)
		arr, err := s.folderService.FindByRSQL(ctx, fmt.Sprintf("tenant_id==\"%s\" and bus_id==\"%s\" and entity_id==\"%s\"", tenantId, cmd.Data.BusId, cmd.Data.EntityId))
		if err != nil {
			return err
		}
		folders := []*model.Folder{}
		for _, f := range arr {
			if strings.HasPrefix(f.FolderPath+"/", cmd.Data.OldName+"/") {
				f.FolderPath = strings.Replace(f.FolderPath+"/", cmd.Data.OldName+"/", cmd.Data.FolderPath+"/", -1)
				f.FolderPath = strings.TrimSuffix(f.FolderPath, "/")
				folders = append(folders, f)
			}
		}

		if len(folders) > 0 {
			opts := idao.NewCallOptions()
			opts.SetUpdateFields([]string{"folder_path"})
			err = s.folderService.UpdateMany(ctx, folders, opts)
			if err != nil {
				return err
			}
		}

		opts := idao.NewCallOptions()
		opts.SetUpdateFields([]string{"name", "alias", "folder_path"})

		folder := model.Folder{}
		folder.Id = cmd.Data.Id
		folder.Name = cmd.Data.Name
		folder.FolderPath = cmd.Data.FolderPath

		err = s.folderService.Update(ctx, &folder, opts)
		if err != nil {
			return err
		}

		err = s.fsService.Rename(cmd.Data.OldName, cmd.Data.FolderPath)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *FolderAPI) updateAlias(ctx context.Context, folder *model.RenameFolder) error {
	opts := idao.NewCallOptions()
	opts.SetUpdateFields([]string{"alias"})
	newFolder := &model.Folder{}
	newFolder.Id = folder.Id
	newFolder.Alias = folder.Alias
	err := s.folderService.Update(ctx, newFolder, opts)
	return err
}

func (s *FolderAPI) SetColor(ctx context.Context, cmd *command.FolderUpdateCommand) error {
	opts := idao.NewCallOptions()
	opts.SetUpdateFields([]string{"color"})
	err := s.folderService.Update(ctx, &cmd.Data, opts)
	return err
}

func (s *FolderAPI) Move(ctx context.Context, cmd *command.FolderMoveCommand) error {
	err := tx.StartTx(ctx, []string{s.folderService.GetConfig().DBKey}, func(ctx context.Context, options ...*store2.SessionOptions) error {

		targetMetas, err := s.folderMetaService.FindByFolderId(ctx, cmd.Data.ParentId)
		if err != nil {
			return err
		}

		targetMeta := s.folderMetaService.GetSourceModel(targetMetas)

		arr, err := s.folderService.FindByRSQL(ctx, fmt.Sprintf("tenant_id==\"%s\" and bus_id==\"%s\" and entity_id==\"%s\"", cmd.Data.TenantId, cmd.Data.BusId, cmd.Data.EntityId))
		if err != nil {
			return err
		}

		dMetas := &[]string{}
		uMetas := &[]*model.FolderMeta{}
		cMetas := &[]*model.FolderMeta{}

		folders := []*model.Folder{}
		for _, f := range arr {
			if strings.HasPrefix(f.FolderPath+"/", cmd.Data.SourcePath+"/") {
				f.FolderPath = strings.Replace(f.FolderPath+"/", cmd.Data.SourcePath+"/", cmd.Data.TargetPath+"/"+cmd.Data.Name+"/", -1)
				f.FolderPath = strings.TrimSuffix(f.FolderPath, "/")

				err = s.folderMetaService.BuildUpdateModels(ctx, f.Id, targetMeta, dMetas, uMetas, cMetas)
				if err != nil {
					return err
				}

				err = s.docService.UpdateDocumentMetas(ctx, f.Id, targetMeta)
				if err != nil {
					return err
				}

				folders = append(folders, f)
			}
		}

		if len(folders) > 0 {
			opts := idao.NewCallOptions()
			opts.SetUpdateFields([]string{"folder_path"})
			err = s.folderService.UpdateMany(ctx, folders, opts)
			if err != nil {
				return err
			}

			if len(*dMetas) > 0 {
				err = s.folderMetaService.DeleteByIds(ctx, *dMetas)
				if err != nil {
					return err
				}
			}
			if len(*cMetas) > 0 {
				err = s.folderMetaService.CreateMany(ctx, *cMetas)
				if err != nil {
					return err
				}
			}
			if len(*uMetas) > 0 {
				err = s.folderMetaService.UpdateMany(ctx, *uMetas)
				if err != nil {
					return err
				}
			}
		}

		opts := idao.NewCallOptions()
		opts.SetUpdateFields([]string{"parent_id"})

		folder := model.Folder{}
		folder.Id = cmd.Data.Id
		folder.ParentId = cmd.Data.ParentId

		err = s.folderService.Update(ctx, &folder, opts)
		if err != nil {
			return err
		}

		err = s.fsService.MoveDir(cmd.Data.SourcePath, cmd.Data.TargetPath+"/"+cmd.Data.Name)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

//func (s *FolderAPI) Move(ctx context.Context, cmd *command.FolderMoveCommand) error {
//	err := tx.StartTx(ctx, []string{s.folderService.GetConfig().DBKey}, func(ctx context.Context, options ...*store2.SessionOptions) error {
//
//		arr, err := s.folderService.FindByRSQL(ctx, fmt.Sprintf("tenant_id==\"%s\" and bus_id==\"%s\" and entity_id==\"%s\"", cmd.Data.TenantId, cmd.Data.BusId, cmd.Data.EntityId))
//		if err != nil {
//			return err
//		}
//		folders := []*model.Folder{}
//		for _, f := range arr {
//			if strings.HasPrefix(f.FolderPath+"/", cmd.Data.SourcePath+"/") {
//				f.FolderPath = strings.Replace(f.FolderPath+"/", cmd.Data.SourcePath+"/", cmd.Data.TargetPath+"/"+cmd.Data.Name+"/", -1)
//				f.FolderPath = strings.TrimSuffix(f.FolderPath, "/")
//				folders = append(folders, f)
//			}
//		}
//
//		if len(folders) > 0 {
//			opts := idao.NewCallOptions()
//			opts.SetUpdateFields([]string{"folder_path"})
//			err = s.folderService.UpdateMany(ctx, folders, opts)
//			if err != nil {
//				return err
//			}
//		}
//
//		opts := idao.NewCallOptions()
//		opts.SetUpdateFields([]string{"parent_id"})
//
//		folder := model.Folder{}
//		folder.Id = cmd.Data.Id
//		folder.ParentId = cmd.Data.ParentId
//
//		err = s.folderService.Update(ctx, &folder, opts)
//		if err != nil {
//			return err
//		}
//
//		err = s.fsService.MoveDir(cmd.Data.SourcePath, cmd.Data.TargetPath+"/"+cmd.Data.Name)
//		if err != nil {
//			return err
//		}
//		return nil
//	})
//	return err
//}

func (s *FolderAPI) Update(ctx context.Context, cmd *command.FolderUpdateCommand) error {
	err := s.folderService.Update(ctx, &cmd.Data)
	return err
}

func (s *FolderAPI) Delete(ctx context.Context, cmd *command.FolderDeleteCommand) error {
	err := tx.StartTx(ctx, []string{s.folderService.GetConfig().DBKey}, func(ctx context.Context, options ...*store2.SessionOptions) error {

		arr, err := s.folderService.FindByRSQL(ctx, fmt.Sprintf("tenant_id==\"%s\" and bus_id==\"%s\" and entity_id==\"%s\"", cmd.Data.TenantId, cmd.Data.BusId, cmd.Data.EntityId))
		if err != nil {
			return err
		}

		var res *idao.Result
		for _, f := range arr {
			if strings.HasPrefix(f.FolderPath+"/", cmd.Data.FolderPath+"/") {
				res = s.fileService.DeleteByRSQL(ctx, fmt.Sprintf("folder_id=='%s'", f.Id))
				if res.Error != nil {
					return res.Error
				}

				data, err := s.docService.FindByRSQL(ctx, fmt.Sprintf("folder_id=='%s'", f.Id))
				if err != nil {
					return err
				}
				if data != nil && len(data) > 0 {
					for _, d := range data {
						res = s.docMetaService.DeleteByDocumentId(ctx, d.Id)
						if res.Error != nil {
							return res.Error
						}
					}
				}

				res = s.docService.DeleteByRSQL(ctx, fmt.Sprintf("folder_id=='%s'", f.Id))
				if res.Error != nil {
					return res.Error
				}

				res = s.folderService.DeleteById(ctx, f.Id)
				if res.Error != nil {
					return res.Error
				}

				res = s.folderMetaService.DeleteByFolderId(ctx, f.Id)
				if res.Error != nil {
					return res.Error
				}
			}
		}
		s.fsService.RemoveAll(cmd.Data.FolderPath)
		return nil
	})

	return err
}

func (s *FolderAPI) FindPaging(ctx context.Context, query *query.FindFolderByFolderIdQuery) (store2.FindPagingResult[*model.FolderView], error) {
	qry := store2.NewFindPagingQueryRequest()
	qry.PageNum = 0
	qry.PageSize = 99999999999999
	qry.Filter = "parent_id=='" + query.FolderId + "'"
	if query.TenantId != "" {
		qry.Filter = qry.Filter + " and tenant_id=='" + query.TenantId + "'"
	}
	if query.BusId != "" {
		qry.Filter = qry.Filter + " and bus_id=='" + query.BusId + "'"
	}
	if query.EntityId != "" {
		qry.Filter = qry.Filter + " and entity_id=='" + query.EntityId + "'"
	}
	qry.Sort = "created_time:desc"
	qry.IsTotalRows = true
	folders, err := s.folderService.FindPaging(ctx, qry)
	if err != nil {
		return nil, err
	}
	if len(folders.GetData()) == 0 {
		return store2.NewFindPagingResult([]*model.FolderView{}, 0, qry, nil), nil
	}
	folderIds := make([]string, 0)
	for _, folder := range folders.GetData() {
		folderIds = append(folderIds, folder.Id)
	}
	allMetas, err := s.folderMetaService.FindByFolderIds(ctx, folderIds)
	if err != nil {
		return nil, err
	}

	fvs := make([]*model.FolderView, 0)
	for _, f := range folders.GetData() {
		fv := s.folderService.Folder2FolderView(f)
		fv.Meta = s.getMeta(allMetas, fv.Id)
		fvs = append(fvs, fv)
	}

	res := store2.NewFindPagingResult(fvs, int64(len(fvs)), qry, nil)

	return res, nil
}

func (s *FolderAPI) getMeta(metas []*model.FolderMeta, folderId string) []*model.FolderMeta {
	folderMetas := make([]*model.FolderMeta, 0)
	if metas == nil || len(metas) == 0 {
		return []*model.FolderMeta{}
	}
	for _, meta := range metas {
		if meta.FolderId == folderId {
			folderMetas = append(folderMetas, meta)
		}
	}
	return folderMetas
}

func (s *FolderAPI) getMetaByFolderId(ctx context.Context, folderId string) ([]*model.FolderMeta, error) {
	return s.folderMetaService.FindByFolderId(ctx, folderId)
}

func (s *FolderAPI) FindTree(ctx context.Context, query *query.FindTreeByParamQuery) []model.FolderTree {
	qry := store2.NewFindPagingQueryRequest()
	qry.PageNum = 0
	qry.PageSize = 99999999999999
	qry.Filter = fmt.Sprintf("tenant_id=='%s' and bus_id=='%s' and entity_id=='%s'", query.TenantId, query.BusId, query.EntityId)
	qry.Sort = "created_time:desc"
	qry.IsTotalRows = true
	res, _ := s.folderService.FindPaging(ctx, qry)
	return s.folderService.TranslateTreeData(res.GetData())
}

func (s *FolderAPI) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.FolderView, error) {
	data, err := s.folderService.FindById(ctx, qry.Id)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	meta, err := s.folderMetaService.FindByFolderId(ctx, data.Id)
	if err != nil {
		return nil, err
	}

	view := s.folderService.Folder2FolderView(data)
	view.Meta = meta

	return view, nil
}
