package hserver

import (
	"context"

	"github.com/PuerkitoBio/goquery"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	element2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/element/base"
	funcs2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/element/funcs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/element/server"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/element/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type factory struct {
}

func NewFactory() element2.Factory {
	return &factory{}
}

func (f *factory) NewBase(srcFileName string, logger logrus.FieldLogger, reader fs.Reader, factory element2.Factory, fsOpts *fsopts.Options) (element2.Base, error) {
	return base.NewBase(srcFileName, logger, reader, factory, fsOpts)
}

func (f *factory) NewFunc(server element2.Server, config *element2.FuncConfig, logger logrus.FieldLogger, reader fs.Reader, pkg *types.CMap[any]) (element2.Func, error) {
	return funcs2.NewFunc(server, config, logger, reader, pkg)
}

func (f *factory) NewServer(httpServer *restapp.HttpServer, srcFileName string, srcFs afero.Fs, factory element2.Factory, env *env.Env, opts ...element2.NewServerOptions) (element2.Server, error) {
	return server.NewServer(httpServer, srcFileName, srcFs, factory, env, opts...)
}

func (f *factory) NewService(server element2.Server, sel *goquery.Selection, srcFileName string, data map[string]any) (element2.Service, error) {
	return service.NewService(server, sel, srcFileName, data)
}

func (f *factory) NewFuncManager(logger logrus.FieldLogger, reader fs.Reader) element2.FuncManager {
	return funcs2.NewFuncManager(logger, reader)
}

func (f *factory) NewWebContext(ctx context.Context, ictx iris.Context) element2.WebContext {
	return NewWebContext(ctx, ictx)
}
