package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type TagRelationAPI struct {
	env            *env.Env
	tagRelationSvc *service.TagRelationService
}

func NewTagRelationAPI(env *env.Env, rootPath string) *TagRelationAPI {
	tagRelationSvc := service.NewTagRelationService()
	return &TagRelationAPI{
		env:            env,
		tagRelationSvc: tagRelationSvc,
	}
}

func (s *TagRelationAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodPost, "/sys/tag-relation", "Create")
	b.Handle(iris.MethodPut, "/sys/tag-relation", "Update")
	b.Handle(iris.MethodDelete, "/sys/tag-relation", "Delete")
	b.Handle(iris.MethodGet, "/sys/tag-relation", "FindPaging")
}

func (s *TagRelationAPI) Create(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		var cmd *command.TagRelationCreateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		res := s.tagRelationSvc.Create(ctx, &cmd.Data)
		return res.Error
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *TagRelationAPI) Update(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		var cmd *command.TagRelationUpdateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		res := s.tagRelationSvc.Update(ctx, &cmd.Data)
		return res.Error
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *TagRelationAPI) Delete(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		var cmd *command.TagRelationDeleteCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		res := s.tagRelationSvc.DeleteById(ctx, cmd.Data.Id)
		return res.Error
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *TagRelationAPI) FindPaging(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		caseId := ictx.URLParam("case-id")
		user, _ := appctx.GetAuthUser(ctx)
		qry := store.NewFindPagingQueryRequest()
		qry.PageNum = 0
		qry.PageSize = 99999999999999
		qry.Filter = "case_id=='" + caseId + "' and creator_id=='" + user.GetId() + "'"
		qry.Sort = "created_time:desc"
		qry.IsTotalRows = true
		res := s.tagRelationSvc.FindPaging(ctx, qry)
		return restapi.SetData(ictx, res)
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}
