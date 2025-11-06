package restapi

import (
	"context"
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type TagAPI struct {
	env            *env.Env
	tagService     *service.TagService
	tagRelationSvc *service.TagRelationService
}

func NewTagAPI(env *env.Env) *TagAPI {
	tagService := service.NewTagService()
	tagRelationSvc := service.NewTagRelationService()
	return &TagAPI{
		env:            env,
		tagService:     tagService,
		tagRelationSvc: tagRelationSvc,
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
		err := tx.StartTx(ctx, []string{s.tagService.GetConfig().DBKey}, func(ctx context.Context, options ...*store2.SessionOptions) error {
			var cmd *command.TagCreateCommand
			if err := ictx.ReadJSON(&cmd); err != nil {
				return err
			}
			res := s.tagService.Create(ctx, &cmd.Data)
			return res.Error
		})
		return err
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *TagAPI) Update(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		err := tx.StartTx(ctx, []string{s.tagService.GetConfig().DBKey}, func(ctx context.Context, options ...*store2.SessionOptions) error {
			var cmd *command.TagUpdateCommand
			if err := ictx.ReadJSON(&cmd); err != nil {
				return err
			}

			arrTagRelation, err := s.tagRelationSvc.FindByRSQL(ctx, fmt.Sprintf("tag_id=='%s'", cmd.Data.Id))
			if err != nil {
				return err
			}
			for _, v := range arrTagRelation {
				v.TagName = cmd.Data.Name
				v.TagColor = cmd.Data.Color
				v.ChangedSource = "TagCenter"
			}
			if len(arrTagRelation) > 0 {
				opts := idao.NewCallOptions()
				opts.SetUpdateFields([]string{"tag_name", "tag_color", "changed_source"})
				s.tagRelationSvc.UpdateMany(ctx, arrTagRelation, opts)
			}

			opts := idao.NewCallOptions()
			opts.SetUpdateFields([]string{"name", "color"})
			res := s.tagService.Update(ctx, &cmd.Data, opts)
			return res.Error
		})
		return err
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *TagAPI) Delete(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		err := tx.StartTx(ctx, []string{s.tagService.GetConfig().DBKey}, func(ctx context.Context, options ...*store2.SessionOptions) error {
			var cmd *command.TagDeleteCommand
			if err := ictx.ReadJSON(&cmd); err != nil {
				return err
			}

			res := s.tagRelationSvc.DeleteByRSQL(ctx, fmt.Sprintf("tag_id=='%s'", cmd.Data.Id))
			if res.Error != nil {
				return res.Error
			}

			res = s.tagService.DeleteById(ctx, cmd.Data.Id)
			return res.Error
		})
		return err
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *TagAPI) FindPaging(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		tenantId := ictx.URLParam("tenant-id")
		caseId := ictx.URLParam("case-id")
		etag := ictx.URLParam("etag")
		//user, _ := appctx.GetAuthUser(ctx)

		var filter string
		if etag == "true" {
			filter = fmt.Sprintf("tenant_id=='%s' and is_e_tag==%s", tenantId, etag)
		} else {
			filter = fmt.Sprintf("tenant_id=='%s' and ((case_id=='%s' and is_e_tag==%s) or is_e_tag==%s)", tenantId, caseId, etag, "true")
		}

		qry := store2.NewFindPagingQueryRequest()
		qry.PageNum = 0
		qry.PageSize = 99999999999999
		qry.Filter = filter
		qry.Sort = "created_time:desc"
		qry.IsTotalRows = true
		res := s.tagService.FindPaging(ctx, qry)
		return restapi.SetData(ictx, res)
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}
