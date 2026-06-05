package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/company-lib/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/company-lib/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type CompanyAPI struct {
	env      *env.Env
	service  *service.CompanyService
	rootPath string
}

func NewCompanyAPI(env *env.Env, rootPath string) *CompanyAPI {
	return &CompanyAPI{
		env:      env,
		service:  service.NewCompanyService(),
		rootPath: rootPath,
	}
}

func (s *CompanyAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath+"/company-lib", "company.CompanyLibAPI", s)
	ctl.GetData("company/search", "Search")
	ctl.GetData("company/full-text-search", "FullTextSearch")
	return ctl
}

func (s *CompanyAPI) Search(ctx context.Context, query *model.CompanyQuery) (*model.CompanyQueryResult, error) {
	return s.service.Search(ctx, query)
}

func (s *CompanyAPI) FullTextSearch(ctx context.Context, query *model.FullQuery) (*model.CompanyQueryResult, error) {
	return s.service.FullTextSearch(ctx, query)
}
