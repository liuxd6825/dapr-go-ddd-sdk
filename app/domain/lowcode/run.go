package lowcode

import (
	/*	app_pkg "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/lowcode/packages/app/pkg"
		app_domain "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/lowcode/packages/app/domain"
		pkg "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/lowcode/packages/pkg"
	*/
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/goserver"
	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/native"
	"github.com/spf13/afero"
)

var RunFun func(int, int) int

func Run(server *restapp.HttpServer) error {
	env := server.EnvConfig()
	if !env.App.HServer.Enable {
		return nil
	}

	hServer := env.App.HServer
	packages := goserver.LoadPackages()

	globals := native.Declarations{}

	buildOptions := &scriggo.BuildOptions{
		Packages: packages,
		Globals:  globals,
	}

	program, err := goserver.Build(server, hServer.BasePath, hServer.WebName, env, false, buildOptions)
	if err != nil {
		return err
	}

	if err := program.Run(nil); err != nil {
		return err
	}
	return nil
}

func build(fs afero.Fs, rootPath string, options *scriggo.BuildOptions) (*scriggo.Program, error) {
	srcFiles, err := goserver.LoadSrcFiles(fs, rootPath)

	if err != nil {
		return nil, err
	}
	loader := scriggo.Files(srcFiles)
	return scriggo.Build(loader, options)
}

func addPackages(target native.Packages, src native.Packages) {
	for key, pkg := range src {
		target[key] = pkg
	}
}

func addGlobals(target native.Declarations, src native.Declarations) {
	for key, pkg := range src {
		target[key] = pkg
	}
}
