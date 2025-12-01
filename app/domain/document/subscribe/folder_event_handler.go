package subscribe

import (
	"context"
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
	fsService         *service.FsService
}

func NewFolderEventSubHandler(env *env.Env, baseUrl string) *FolderEventSubHandler {
	return &FolderEventSubHandler{
		rootPath:          baseUrl,
		env:               env,
		folderService:     service.NewFolderService(),
		folderMetaService: service.NewFolderMetaService(),
		fsService:         service.NewFsService(),
	}
}

func (s *FolderEventSubHandler) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, "subscribe/document/event", "FolderEventSubHandler", s)
	ctl.EventHandle("folder-create-event", "FolderCreateEvent")
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
