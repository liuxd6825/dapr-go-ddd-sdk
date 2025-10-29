package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/dict/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/dict/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/dict/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/dict/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type DictTypeAPI struct {
	env             *env.Env
	dictTypeService *service.DictTypeService
	rootPath        string
}

func NewDictTypeAPI(env *env.Env, rootPath string) *DictTypeAPI {
	return &DictTypeAPI{
		env:             env,
		dictTypeService: service.NewDictTypeService(),
		rootPath:        rootPath,
	}
}

func (s *DictTypeAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	s.dictTypeService = service.NewDictTypeService()
	ctl := restapi.NewController(app, s.rootPath+"/sys", "sys.DictTypeApi", s)
	ctl.Post("/dictionary-type", "Create")
	ctl.Put("/dictionary-type", "Update")
	ctl.Delete("/dictionary-type", "Delete", restapi.WithParamsInBody(true))
	ctl.Delete("/dictionary-type:deleteBatch", "DeleteBatch", restapi.WithParamsInBody(true))
	ctl.GetOne("/dictionary-type/{id}", "FindById")
	ctl.GetPaging("/dictionary-type", "FindPaging")
	ctl.GetData("/dictionary-type:code/{code}", "FindByCode")
	ctl.GetData("/dictionary-type:findAllType", "FindAll")
	return ctl
}

func (s *DictTypeAPI) Create(ctx context.Context, cmd *command.DictTypeCreateCommand) error {
	return s.dictTypeService.Create(ctx, cmd)
}

func (s *DictTypeAPI) Update(ctx context.Context, cmd *command.DictTypeUpdateCommand) error {
	return s.dictTypeService.Update(ctx, cmd)
}

func (s *DictTypeAPI) Delete(ctx context.Context, cmd *command.DictTypeDeleteCommand) error {
	return s.dictTypeService.Delete(ctx, cmd)
}

func (s *DictTypeAPI) DeleteBatch(ctx context.Context, cmd *command.DictTypeDeleteBatchCommand) error {
	return s.dictTypeService.DeleteBatch(ctx, cmd)
}

func (s *DictTypeAPI) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.DictType, error) {
	return s.dictTypeService.FindById(ctx, qry)
}

func (s *DictTypeAPI) FindPaging(ctx context.Context, qry *idao.FindPagingQueryRequest) (idao.FindPagingResult[*model.DictType], error) {
	return s.dictTypeService.FindPaging(ctx, qry)
}

func (s *DictTypeAPI) FindByCode(ctx context.Context, qry *query.FindByCodeQuery) ([]*model.DictType, error) {
	return s.dictTypeService.FindByCode(ctx, qry.Code)
}

func (s *DictTypeAPI) FindAll(ctx context.Context) ([]*model.DictType, error) {
	qry := store.NewFindPagingQueryRequest()
	qry.PageNum = 0
	qry.PageSize = 99999999999999
	//qry.Sort = "created_time:desc"
	qry.IsTotalRows = true
	return s.dictTypeService.FindAll(ctx, qry)
}
