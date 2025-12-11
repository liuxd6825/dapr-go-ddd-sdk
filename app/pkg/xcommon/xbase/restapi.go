package xbase

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

/*type RestApi[T any] struct {
	config RestApiConfig
	env    *env.Env
}

type RestApiConfig struct {
	RootPath string // 根路径名
	Domain   string // 领域名称，路径名
	Resource string // 资源名称
}

func NewRestApi[T any](env *env.Env, config RestApiConfig) *RestApi[T] {
	return &RestApi[T]{
		config: config,
		env:    env,
	}
}

func (s *RestApi[T]) NewAPIController(app *iris.Application) *restapi.ApiController {
	url := s.config.RootPath
	if s.config.Domain != "" {
		url = url + "/" + s.config.Domain
	}
	if s.config.Resource != "" {
		url = url + "/" + s.config.Resource
	}

	appName := s.config.Domain + "." + s.config.Resource
	controller := restapi.NewController(app, url, appName, s)
	controller.Post("", "Create", restapi.NewAPIOptions())
	controller.Put("", "Update", restapi.NewAPIOptions())
	controller.GetOne("/{id}", "FindById", restapi.NewAPIOptions())
	controller.GetPaging("", "FindPaging")
	return controller
}

func (s *RestApi[T]) Create(ctx context.Context, cmd *Command[T]) (T, error) {
	return s.service.Create(ctx, cmd)
}

func (s *RestApi[T]) Update(ctx context.Context, cmd *Command[T]) (T, error) {
	return s.service.Update(ctx, cmd)
}

func (s *RestApi[T]) Delete(ctx context.Context, cmd *DeleteByIdCommand) error {
	return s.service.Delete(ctx, cmd)
}

func (s *RestApi[T]) FindById(ctx context.Context, qry *restapi.FindByIdRequest) (T, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *RestApi[T]) FindPaging(ctx context.Context, qry *restapi.FindPagingQueryRequest) store.FindPagingResult[T] {
	return s.service.FindPaging(ctx, qry)
}

func (s *RestApi[T]) FindByAll(ctx context.Context) ([]T, error) {
	return s.service.FindAll(ctx)
}

func (s *RestApi[T]) GetRootPath() string {
	return s.config.RootPath
}

func (s *RestApi[T]) GetEnv() *env.Env {
	return s.env
}
*/

type ApiConfig struct {
	RootPath string
	Domain   string
	Resource string
	Methods  []string
}

func NewAPIController(app *iris.Application, apiObj any, cfg ApiConfig) *restapi.ApiController {
	url := cfg.RootPath
	if cfg.Domain != "" {
		url = url + "/" + cfg.Domain
	}
	if cfg.Resource != "" {
		url = url + "/" + cfg.Resource
	}

	appName := cfg.Domain + "." + cfg.Resource
	controller := restapi.NewController(app, url, appName, apiObj)
	controller.Post("", "Create")
	controller.Post("", "CreateMany")
	controller.Put("", "Update")
	controller.Put("", "UpdateMany")
	controller.Delete("", "Delete")
	controller.Put("", "Update")
	controller.GetOne("/{id}", "FindById")
	controller.GetPaging("", "FindPaging")
	return controller
}
