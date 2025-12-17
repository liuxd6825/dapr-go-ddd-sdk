package restapi

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/query"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type TagTypeAPI struct {
	env            *env.Env
	tagTypeService *service.TagTypeService
	tagService     *service.TagService
	rootPath       string
}

func NewTagTypeAPI(env *env.Env, rootPath string) *TagTypeAPI {
	tagTypeService := service.NewTagTypeService()
	tagService := service.NewTagService()
	return &TagTypeAPI{
		env:            env,
		tagTypeService: tagTypeService,
		tagService:     tagService,
		rootPath:       rootPath,
	}
}

func (s *TagTypeAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	s.tagTypeService = service.NewTagTypeService()
	ctl := restapi.NewController(app, s.rootPath+"/sys", "sys.TagTypeAPI", s)
	ctl.Post("/tag-type", "Create")
	ctl.Put("/tag-type", "Update")
	ctl.Delete("/tag-type", "Delete", restapi.WithParamsInBody(true))
	ctl.GetPaging("/tag-type", "FindPaging")
	return ctl
}

func (s *TagTypeAPI) Create(ctx context.Context, cmd *command.TagTypeCreateCommand) error {
	return s.tagTypeService.Create(ctx, cmd)
}

func (s *TagTypeAPI) Update(ctx context.Context, cmd *command.TagTypeUpdateCommand) error {
	opts := idao.NewCallOptions()
	opts.SetUpdateFields([]string{"name", "color"})
	return s.tagTypeService.Update(ctx, &cmd.Data, opts)
}

func (s *TagTypeAPI) Delete(ctx context.Context, cmd *command.TagTypeDeleteCommand) error {
	return s.tagTypeService.Delete(ctx, cmd)
}

func (s *TagTypeAPI) DeleteByParentId(ctx context.Context, id string) error {
	//rows, err := s.tagTypeService.FindByRSQL(ctx, fmt.Sprintf("parent_id=='%s'", id))
	//if err != nil {
	//	return err
	//}
	//if len(rows) > 0 {
	//	for _, row := range rows {
	//		err = s.DeleteByParentId(ctx, row.Id)
	//		if err != nil {
	//			return err
	//		}
	//	}
	//} else {
	//	res := s.tagService.DeleteByRSQL(ctx, fmt.Sprintf("tag_type_id=='%s'", id))
	//	if res.Error != nil {
	//		return res.Error
	//	}
	//	res = s.tagTypeService.DeleteById(ctx, id)
	//	if res.Error != nil {
	//		return res.Error
	//	}
	//}
	return nil
}

func (s *TagTypeAPI) FindPaging(ctx context.Context, qry *query.FindPagingTagQuery) (idao.FindPagingResult[*model.TagType], error) {
	var filter string
	if qry.ETag == true {
		filter = fmt.Sprintf("tenant_id=='%s' and is_e_tag==%v", qry.TenantId, qry.ETag)
	} else {
		filter = fmt.Sprintf("tenant_id=='%s' and ((case_id=='%s' and is_e_tag==%v) or is_e_tag==%s)", qry.TenantId, qry.CaseId, qry.ETag, "true")
	}

	qry1 := store2.NewFindPagingQueryRequest()
	qry1.PageNum = 0
	qry1.PageSize = 99999999999999
	qry1.Filter = filter
	qry1.Sort = "created_time:desc"
	qry1.IsTotalRows = true
	return s.tagTypeService.FindPaging(ctx, qry1)
}
