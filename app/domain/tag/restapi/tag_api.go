package restapi

import (
	"context"
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type TagAPI struct {
	env            *envs.Env
	tagService     *service.TagService
	tagRelationSvc *service.TagRelationService
	rootPath       string
	mongoDbKey     string
}

func NewTagAPI(env *envs.Env, rootPath string, mongoDbKey string) *TagAPI {
	tagService := service.NewTagService(mongoDbKey)
	tagRelationSvc := service.NewTagRelationService(mongoDbKey)
	return &TagAPI{
		mongoDbKey:     mongoDbKey,
		env:            env,
		tagService:     tagService,
		tagRelationSvc: tagRelationSvc,
		rootPath:       rootPath,
	}
}

func (s *TagAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	s.tagService = service.NewTagService(s.mongoDbKey)
	ctl := restapi.NewController(app, s.rootPath+"/sys", "sys.TagAPI", s)
	ctl.Post("/tag", "Create")
	ctl.Put("/tag", "Update")
	ctl.Delete("/tag", "Delete", restapi.WithParamsInBody(true))
	ctl.GetPaging("/tag", "FindPaging")
	return ctl
}

func (s *TagAPI) Create(ctx context.Context, cmd *command.TagCreateCommand) error {
	return s.tagService.Create(ctx, cmd)
}

func (s *TagAPI) Update(ctx context.Context, cmd *command.TagUpdateCommand) error {
	return tx.StartTx(ctx, tx.NewTxCfg(s.mongoDbKey), func(ctx context.Context, options ...*store.SessionOptions) error {
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
			err = s.tagRelationSvc.UpdateMany(ctx, arrTagRelation, opts)
			if err != nil {
				return err
			}
		}

		opts := idao.NewCallOptions()
		opts.SetUpdateFields([]string{"name", "color"})
		return s.tagService.Update(ctx, &cmd.Data, opts)
	})
}

func (s *TagAPI) Delete(ctx context.Context, cmd *command.TagDeleteCommand) error {
	return tx.StartTx(ctx, tx.NewTxCfg(s.mongoDbKey), func(ctx context.Context, options ...*store.SessionOptions) error {
		err := s.tagRelationSvc.DeleteByRSQL(ctx, fmt.Sprintf("tag_id=='%s'", cmd.Data.Id))
		if err != nil {
			return err
		}
		return s.tagService.Delete(ctx, cmd)
	})
}

func (s *TagAPI) FindPaging(ctx context.Context, qry *query.FindPagingTagQuery) (idao.FindPagingResult[*model.Tag], error) {
	var filter string
	if qry.ETag == true {
		filter = fmt.Sprintf("tenant_id=='%s' and is_e_tag==%v", qry.TenantId, qry.ETag)
	} else {
		filter = fmt.Sprintf("tenant_id=='%s' and ((case_id=='%s' and is_e_tag==%v) or is_e_tag==%s)", qry.TenantId, qry.CaseId, qry.ETag, "true")
	}

	qry1 := store.NewFindPagingQueryRequest()
	qry1.PageNum = 0
	qry1.PageSize = 99999999999999
	qry1.Filter = filter
	qry1.Sort = "created_time:desc"
	qry1.IsTotalRows = true
	return s.tagService.FindPaging(ctx, qry1)
}
