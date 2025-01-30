package hserver

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element/api"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element/api/request"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element/base"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element/script"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element/server"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element/sub"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
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

func (f *factory) NewServer(httpServer *restapp.HttpServer, srcFileName string, srcFs afero.Fs, factory element.Factory, env common.IEnvConfig, opts ...element.NewServerOptions) (element.Server, error) {
	return server.NewServer(httpServer, srcFileName, srcFs, factory, env, opts...)
}

func (f *factory) NewApiService(server element.Server, sel *goquery.Selection, srcFileName string, data map[string]any) (element.ApiService, error) {
	return api.NewApiService(server, sel, srcFileName, data)
}

func (f *factory) NewApiRequest(server element.Server, service element.ApiService, srcFileName string, config *element.ApiRequestConfig) (element.ApiRequest, error) {
	return request.NewApiRequest(server, service, srcFileName, config)
}

func (f *factory) NewSubService(server element.Server, sel *goquery.Selection, srcFileName string, data map[string]any) (element.SubService, error) {
	return sub.NewSubService(server, sel, srcFileName, data)
}

func (f *factory) NewSubEvent(server element.Server, service element.SubService, srcFileName string, config *element.SubEventConfig) (element.SubEvent, error) {
	return sub.NewSubEvent(server, service, srcFileName, config)
}

func (f *factory) NewScriptManager(logger logrus.FieldLogger, reader fs.Reader) element.ScriptManager {
	return script.NewScriptManager(logger, reader)
}

func (f *factory) NewWebContext(ctx context.Context, ictx iris.Context, request element.ApiRequest) element.WebContext {
	return NewWebContext(ctx, ictx, request)
}
