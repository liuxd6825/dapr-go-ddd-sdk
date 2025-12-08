package subscribe

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"strings"
	"time"

	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/event"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

// FolderEventSubHandler
// @Description: 处理数据导入事件
type FolderEventSubHandler struct {
	rootPath          string
	env               *env.Env
	folderService     *service.FolderService
	folderMetaService *service.FolderMetaService
	docService        *service.DocumentService
	docMetaService    *service.DocumentMetaService
	fileService       *service.FileService
	fsService         *service.FsService
}

func NewFolderEventSubHandler(env *env.Env, baseUrl string) *FolderEventSubHandler {
	return &FolderEventSubHandler{
		rootPath:          baseUrl,
		env:               env,
		folderService:     service.NewFolderService(),
		folderMetaService: service.NewFolderMetaService(),
		fsService:         service.NewFsService(),
		docService:        service.NewDocumentService(),
		docMetaService:    service.NewDocumentMetaService(),
		fileService:       service.NewFileService(),
	}
}

func (s *FolderEventSubHandler) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, "subscribe/document/event", "FolderEventSubHandler", s)
	ctl.EventHandle("folder-create-event", "FolderCreateEvent")
	ctl.EventHandle("folder-update-alias-event", "FolderUpdateAliasEvent")
	ctl.EventHandle("folder-delete-event", "FolderDeleteEvent")
	ctl.Handle(iris.MethodOptions, "folder-create-event", "Check")
	return ctl
}

func (s *FolderEventSubHandler) Check(ctx context.Context) error {
	logs.Infofmt(ctx, "document/subscribe/folder/folder-create-event:check")
	return nil
}

func (s *FolderEventSubHandler) FolderCreateEvent(ctx context.Context, event *event.FolderCreateEvent) error {
	logs.Infofmt(ctx, "%s eventId:%s; occurredOn:%s; ", event.EventType, event.Id, event.CreatedTime.Format(time.DateTime))

	return tx.StartTx(ctx, tx.NewTxCfg(config.DBKey), func(ctx context.Context, options ...*store.SessionOptions) error {
		folder, _ := model2.NewFolder()
		folder.Id = event.Data.Id
		folder.BusId = event.Data.BusId
		folder.EntityId = event.Data.CaseId
		folder.RootId = event.Data.RootId
		folder.RootPath = event.Data.RootPath
		folder.FolderPath = event.Data.FolderPath
		folder.CaseId = event.Data.CaseId
		folder.Name = event.Data.Name
		folder.Alias = event.Data.Alias
		folder.ParentId = event.Data.ParentId
		folder.DisabledFrontEdit = event.Data.DisabledFrontEdit

		err := s.folderService.CreateData(ctx, folder)
		if err != nil {
			return err
		}

		var metas []*model2.FolderMeta
		for _, m := range event.Data.Meta {
			fm := model2.NewFolderMeta()
			fm.Id = idutils.NewId()
			fm.CaseId = event.Data.CaseId
			fm.FolderId = folder.Id
			fm.SourceType = m.SourceType
			fm.Source = m.Source
			fm.Name = m.Name
			fm.Value = m.Value
			metas = append(metas, fm)
		}

		err = s.folderMetaService.CreateMany(ctx, metas)
		if err != nil {
			return err
		}

		s.fsService.MkdirAll(folder.FolderPath)
		return nil
	})
}

func (s *FolderEventSubHandler) FolderUpdateAliasEvent(ctx context.Context, event *event.FolderUpdateAliasEvent) error {

	folder, _ := model2.NewFolder()
	folder.Id = event.Data.Id
	folder.Alias = event.Data.Alias

	opts := idao.NewCallOptions()
	opts.SetUpdateFields([]string{"alias"})

	return s.folderService.Update(ctx, folder, opts)
}

func (s *FolderEventSubHandler) FolderDeleteEvent(ctx context.Context, event *event.FolderDeleteEvent) error {
	logs.Infofmt(ctx, "%s eventId:%s; occurredOn:%s; ", event.EventType, event.Id, event.CreatedTime.Format(time.DateTime))

	return tx.StartTx(ctx, tx.NewTxCfg(config.DBKey), func(ctx context.Context, options ...*store.SessionOptions) error {
		arr, err := s.folderService.FindByRSQL(ctx, fmt.Sprintf("tenant_id==\"%s\" and bus_id==\"%s\" and entity_id==\"%s\"", event.Data.TenantId, event.Data.BusId, event.Data.EntityId))
		if err != nil {
			return err
		}

		if event.Data.Ids != nil && len(event.Data.Ids) > 0 {
			for _, id := range event.Data.Ids {
				folder, err := s.folderService.FindById(ctx, id)
				if err != nil {
					return err
				}
				var res *idao.Result
				for _, f := range arr {
					if strings.HasPrefix(f.FolderPath+"/", folder.FolderPath+"/") {
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
				s.fsService.RemoveAll(folder.FolderPath)
			}
		}

		return nil
	})
}
