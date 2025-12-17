package restapi

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type TagRelationAPI struct {
	env            *env.Env
	tagRelationSvc *service.TagRelationService
	rootPath       string
}

func NewTagRelationAPI(env *env.Env, rootPath string) *TagRelationAPI {
	tagRelationSvc := service.NewTagRelationService()
	return &TagRelationAPI{
		env:            env,
		tagRelationSvc: tagRelationSvc,
		rootPath:       rootPath,
	}
}

func (s *TagRelationAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	s.tagRelationSvc = service.NewTagRelationService()
	ctl := restapi.NewController(app, s.rootPath+"/sys", "sys.TagRelationAPI", s)
	ctl.Post("/tag-relation", "Create")
	ctl.Put("/tag-relation", "Update")
	ctl.Delete("/tag-relation", "Delete", restapi.WithParamsInBody(true))
	ctl.GetPaging("/tag-relation", "FindPaging")
	return ctl
}

func (s *TagRelationAPI) Create(ctx context.Context, cmd *command.TagRelationCreateCommand) error {
	return s.tagRelationSvc.Create(ctx, cmd)
}

func (s *TagRelationAPI) Update(ctx context.Context, cmd *command.TagRelationUpdateCommand) error {
	return s.tagRelationSvc.Update(ctx, &cmd.Data)
}

func (s *TagRelationAPI) Delete(ctx context.Context, cmd *command.TagRelationDeleteCommand) error {
	return s.tagRelationSvc.Delete(ctx, cmd)
}

func (s *TagRelationAPI) FindPaging(ctx context.Context, qry *query.FindPagingByCaseIdQuery) (idao.FindPagingResult[*model.TagRelation], error) {

	user, _ := appctx.GetAuthUser(ctx)
	qry1 := store.NewFindPagingQueryRequest()
	qry1.PageNum = 0
	qry1.PageSize = 99999999999999
	qry1.Filter = "case_id=='" + qry.CaseId + "' and creator_id=='" + user.GetId() + "'"
	qry1.Sort = "created_time:desc"
	qry1.IsTotalRows = true
	return s.tagRelationSvc.FindPaging(ctx, qry1)
}
