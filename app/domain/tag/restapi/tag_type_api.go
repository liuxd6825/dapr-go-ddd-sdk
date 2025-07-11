package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type TagTypeAPI struct {
	env            *env.Env
	tagTypeService *service.TagTypeService
}

func NewTagTypeAPI(env *env.Env, rootPath string) *TagTypeAPI {
	tagTypeService := service.NewTagTypeService()
	return &TagTypeAPI{
		env:            env,
		tagTypeService: tagTypeService,
	}
}

func (s *TagTypeAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodPost, "/sys/tag-type", "Create")
	b.Handle(iris.MethodPut, "/sys/tag-type", "Update")
	b.Handle(iris.MethodDelete, "/sys/tag-type", "Delete")
	b.Handle(iris.MethodGet, "/sys/tag-type", "FindPaging")
}

func (s *TagTypeAPI) Create(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		err := tx.StartTx(ctx, []string{s.tagTypeService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
			var cmd *command.TagTypeCreateCommand
			if err := ictx.ReadJSON(&cmd); err != nil {
				return err
			}

			s.tagTypeService.Create(ctx, &cmd.Data)

			return nil
		})

		return err
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *TagTypeAPI) Update(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		var cmd *command.TagTypeUpdateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		s.tagTypeService.Update(ctx, &cmd.Data)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *TagTypeAPI) Delete(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		err := tx.StartTx(ctx, []string{s.tagTypeService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
			var cmd *command.TagTypeDeleteCommand
			if err := ictx.ReadJSON(&cmd); err != nil {
				return err
			}

			//has := s.folderService.HasChildren(ctx, cmd.Data.Id)
			//if has {
			//	return errors.New("存在子目录")
			//}
			//has = s.docService.HasDocumentByFolder(ctx, cmd.Data.Id)
			//if has {
			//	return errors.New("存在文档")
			//}

			s.tagTypeService.DeleteById(ctx, cmd.Data.Id)
			return nil
		})

		return err
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *TagTypeAPI) FindPaging(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		parentId := ictx.URLParam("folder-id")
		qry := store.NewFindPagingQueryRequest()
		qry.PageNum = 0
		qry.PageSize = 99999999999999
		qry.Filter = "parent_id=='" + parentId + "'"
		qry.Sort = "created_time:desc"
		qry.IsTotalRows = true
		res := s.tagTypeService.FindPaging(ctx, qry)
		return restapi.SetData(ictx, res)
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}
