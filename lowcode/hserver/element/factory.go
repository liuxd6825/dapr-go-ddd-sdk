package element

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// Factory
// @Description: 工厂类
type Factory interface {
	//
	// NewBase
	//  @Description:  生成Base接口类
	//  @param srcFileName  源文件
	//  @param logger 日志
	//  @param reader 读
	//  @param fsOpts 参数
	//  @return Base
	//  @return error
	//
	NewBase(srcFileName string, logger logrus.FieldLogger, reader fs.Reader, factory Factory, fsOpts *fsopts.Options) (Base, error)

	//
	// NewServer
	//  @Description:
	//  @param app
	//  @param srcFileName
	//  @param srcFs
	//  @param env
	//  @param opts
	//  @return Server
	//  @return error
	//
	NewServer(app *iris.Application, srcFileName string, srcFs afero.Fs, factory Factory, env common.IEnvConfig, opts ...NewServerOptions) (Server, error)

	//
	// NewService
	//  @Description:
	//  @param server
	//  @param html
	//  @param srcFileName
	//  @param data
	//  @return Service
	//  @return error
	//
	NewService(server Server, html []byte, srcFileName string, data map[string]any) (Service, error)

	//
	// NewRequest
	//  @Description:
	//  @param server
	//  @param service
	//  @param srcFileName
	//  @param config
	//  @return Request
	//  @return error
	//
	NewRequest(server Server, service Service, srcFileName string, config *RequestConfig) (Request, error)

	//
	// NewScript
	//  @Description:
	//  @Description:
	//  @param config
	//  @param logger
	//  @param reader
	//  @param pkg
	//  @return *Script
	//  @return error
	//
	NewScript(config *ScriptConfig, logger logrus.FieldLogger, reader fs.Reader, pkg *types.CMap[any]) (Script, error)

	//
	// NewScriptManager
	//  @Description:
	//  @param logger
	//  @param reader
	//  @return *ScriptManager
	//
	NewScriptManager(logger logrus.FieldLogger, reader fs.Reader) ScriptManager

	NewWebContext(ctx context.Context, ictx iris.Context, request Request) WebContext
}
