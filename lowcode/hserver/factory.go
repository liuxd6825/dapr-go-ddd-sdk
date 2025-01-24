package hserver

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element/base"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element/script"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element/server"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element/service/request"
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

func (f *factory) NewScript(config *element.ScriptConfig, logger logrus.FieldLogger, reader fs.Reader, pkg *types.CMap[any]) (element.Script, error) {
	return script.NewScript(config, logger, reader, pkg)
}

func (f *factory) NewServer(app *iris.Application, srcFileName string, srcFs afero.Fs, factory element.Factory, env common.IEnvConfig, opts ...element.NewServerOptions) (element.Server, error) {
	return server.NewServer(app, srcFileName, srcFs, factory, env, opts...)
}

func (f *factory) NewService(server element.Server, html []byte, srcFileName string, data map[string]any) (element.Service, error) {
	return service.NewService(server, html, srcFileName, data)
}

func (f *factory) NewRequest(server element.Server, service element.Service, srcFileName string, config *element.RequestConfig) (element.Request, error) {
	return request.NewRequest(server, service, srcFileName, config)
}

func (f *factory) NewScriptManager(logger logrus.FieldLogger, reader fs.Reader) element.ScriptManager {
	return script.NewScriptManager(logger, reader)
}

func (f *factory) NewWebContext(ctx context.Context, ictx iris.Context, request element.Request) element.WebContext {
	return NewWebContext(ctx, ictx, request)
}
