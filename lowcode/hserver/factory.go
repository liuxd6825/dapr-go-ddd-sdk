package hserver

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element/funcs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsopts"

	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element/base"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element/server"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type factory struct {
}

func NewFactory() element.Factory {
	return &factory{}
}

func (f *factory) NewBase(srcFileName string, logger logrus.FieldLogger, reader fs.Reader, factory element.Factory, fsOpts *fsopts.Options) (element.Base, error) {
	return base.NewBase(srcFileName, logger, reader, factory, fsOpts)
}

func (f *factory) NewFunc(server element.Server, config *element.FuncConfig, logger logrus.FieldLogger, reader fs.Reader, pkg *types.CMap[any]) (element.Func, error) {
	return funcs.NewFunc(server, config, logger, reader, pkg)
}

func (f *factory) NewServer(httpServer *restapp.HttpServer, srcFileName string, srcFs afero.Fs, factory element.Factory, env env.IEnvConfig, opts ...element.NewServerOptions) (element.Server, error) {
	return server.NewServer(httpServer, srcFileName, srcFs, factory, env, opts...)
}

func (f *factory) NewService(server element.Server, sel *goquery.Selection, srcFileName string, data map[string]any) (element.Service, error) {
	return service.NewService(server, sel, srcFileName, data)
}

func (f *factory) NewFuncManager(logger logrus.FieldLogger, reader fs.Reader) element.FuncManager {
	return funcs.NewFuncManager(logger, reader)
}

func (f *factory) NewWebContext(ctx context.Context, ictx iris.Context) element.WebContext {
	return NewWebContext(ctx, ictx)
}
