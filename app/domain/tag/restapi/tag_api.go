package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type TagAPI struct {
	env        *env.Env
	tagService *service.TagService
}

func NewTagAPI(env *env.Env) *TagAPI {
	tagService := service.NewTagService()
	return &TagAPI{
		env:        env,
		tagService: tagService,
	}
}

func (s *TagAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodPost, "/sys/tag", "Create")
	b.Handle(iris.MethodPut, "/sys/tag", "Update")
	b.Handle(iris.MethodDelete, "/sys/tag", "Delete")
	b.Handle(iris.MethodGet, "/sys/tag", "FindPaging")
}

func (s *TagAPI) Create(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		err := tx.StartTx(ctx, []string{s.tagService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
			var cmd *command.TagCreateCommand
			if err := ictx.ReadJSON(&cmd); err != nil {
				return err
			}
			s.tagService.Create(ctx, &cmd.Data)
			return nil
		})
		return err
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *TagAPI) Update(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		err := tx.StartTx(ctx, []string{s.tagService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
			var cmd *command.TagUpdateCommand
			if err := ictx.ReadJSON(&cmd); err != nil {
				return err
			}

			opts := idao.NewCallOptions()
			opts.SetUpdateFields([]string{"name", "color"})

			s.tagService.Update(ctx, &cmd.Data, opts)
			return nil
		})
		return err
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *TagAPI) Delete(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		err := tx.StartTx(ctx, []string{s.tagService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
			var cmd *command.TagDeleteCommand
			if err := ictx.ReadJSON(&cmd); err != nil {
				return err
			}
			s.tagService.DeleteById(ctx, cmd.Data.Id)
			return nil
		})
		return err
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *TagAPI) FindPaging(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		etag := ictx.URLParam("etag")
		//user, _ := appctx.GetAuthUser(ctx)
		qry := store.NewFindPagingQueryRequest()
		qry.PageNum = 0
		qry.PageSize = 99999999999999
		qry.Filter = "is_e_tag==" + etag
		qry.Sort = "created_time:desc"
		qry.IsTotalRows = true
		res := s.tagService.FindPaging(ctx, qry)
		return restapi.SetData(ictx, res)
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}
