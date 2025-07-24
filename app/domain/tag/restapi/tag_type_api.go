package restapi

import (
	"context"
	"fmt"
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

type TagTypeAPI struct {
	env            *env.Env
	tagTypeService *service.TagTypeService
	tagService     *service.TagService
}

func NewTagTypeAPI(env *env.Env) *TagTypeAPI {
	tagTypeService := service.NewTagTypeService()
	tagService := service.NewTagService()
	return &TagTypeAPI{
		env:            env,
		tagTypeService: tagTypeService,
		tagService:     tagService,
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

			res := s.tagTypeService.Create(ctx, &cmd.Data)

			return res.Error
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

		opts := idao.NewCallOptions()
		opts.SetUpdateFields([]string{"name", "color"})

		res := s.tagTypeService.Update(ctx, &cmd.Data, opts)
		return res.Error
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

			res := s.tagTypeService.DeleteById(ctx, cmd.Data.Id)

			//err := s.DeleteByParentId(ctx, cmd.Data.Id)

			return res.Error
		})

		return err
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
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

func (s *TagTypeAPI) FindPaging(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		tenantId := ictx.URLParam("tenant-id")
		caseId := ictx.URLParam("case-id")
		etag := ictx.URLParam("etag")

		var filter string
		if etag == "true" {
			filter = fmt.Sprintf("tenant_id=='%s' and is_e_tag==%s", tenantId, etag)
		} else {
			filter = fmt.Sprintf("tenant_id=='%s' and ((case_id=='%s' and is_e_tag==%s) or is_e_tag==%s)", tenantId, caseId, etag, "true")
		}

		qry := store.NewFindPagingQueryRequest()
		qry.PageNum = 0
		qry.PageSize = 99999999999999
		qry.Filter = filter
		qry.Sort = "created_time:desc"
		qry.IsTotalRows = true
		res := s.tagTypeService.FindPaging(ctx, qry)
		return restapi.SetData(ictx, res)
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}
