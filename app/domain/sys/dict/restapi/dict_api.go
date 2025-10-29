package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/dict/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/dict/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/dict/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/dict/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type DictAPI struct {
	env         *env.Env
	dictService *service.DictService
	rootPath    string
}

func NewDictAPI(env *env.Env, rootPath string) *DictAPI {
	return &DictAPI{
		env:         env,
		dictService: service.NewDictService(),
		rootPath:    rootPath,
	}
}

func (s *DictAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	s.dictService = service.NewDictService()
	ctl := restapi.NewController(app, s.rootPath+"/sys", "sys.DictApi", s)
	ctl.Post("/dictionary", "Create")
	ctl.Put("/dictionary", "Update")
	ctl.Delete("/dictionary", "Delete", restapi.WithParamsInBody(true))
	ctl.Delete("/dictionary:deleteBatch", "DeleteBatch", restapi.WithParamsInBody(true))
	ctl.GetOne("/dictionary/{id}", "FindById")
	ctl.GetPaging("/dictionary", "FindPaging")
	ctl.GetData("/dict-type/{dictType}/dictionary:code/{code}", "FindByDictTypeAndCode")
	ctl.GetData("/dict-type/{dictTypeCode}/dictionary", "FindByDictTypeCodeAndParams")
	return ctl
}

func (s *DictAPI) Create(ctx context.Context, cmd *command.DictCreateCommand) error {
	return s.dictService.Create(ctx, cmd)
}

func (s *DictAPI) Update(ctx context.Context, cmd *command.DictUpdateCommand) error {
	return s.dictService.Update(ctx, cmd)
}

func (s *DictAPI) Delete(ctx context.Context, cmd *command.DictDeleteCommand) error {
	return s.dictService.Delete(ctx, cmd)
}

func (s *DictAPI) DeleteBatch(ctx context.Context, cmd *command.DictDeleteBatchCommand) error {
	return s.dictService.DeleteBatch(ctx, cmd)
}

func (s *DictAPI) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.Dict, error) {
	return s.dictService.FindById(ctx, qry)
}

func (s *DictAPI) FindPaging(ctx context.Context, qry *query.FindPagingByCaseIdQuery) (idao.FindPagingResult[*model.Dict], error) {
	return s.dictService.FindPaging(ctx, qry.CaseId, qry)
}

func (s *DictAPI) FindByDictTypeAndCode(ctx context.Context, qry *query.FindByDictTypeAndCodeQuery) ([]*model.Dict, error) {
	return s.dictService.FindByDictTypeAndCode(ctx, qry.DictType, qry.Code)
}

func (s *DictAPI) FindByDictTypeCodeAndParams(ctx context.Context, qry *query.FindByDictTypeCodeAndParamsQuery) ([]*model.Dict, error) {
	return s.dictService.FindByDictTypeCodeAndParams(ctx, qry.DictTypeCode, qry.CaseTypeId, qry.CaseId)
}
