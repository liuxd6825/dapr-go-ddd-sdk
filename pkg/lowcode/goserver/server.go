package goserver

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/native"
	"github.com/spf13/afero"
)

type RunOptions func(options *scriggo.BuildOptions, packages native.Packages) error

// Build
//
//	@Description:
//	@param fileName
//	@param srcFsName
//	@param webFsName
//	@param httpServer
//	@return error
func Build(httpServer *restapp.HttpServer, srcFsName string, webFsName string, env *env.Env, autoRestart bool, buildOptions *scriggo.BuildOptions, opts ...RunOptions) (*scriggo.Program, error) {
	var err error
	defer func() {
		err = errors.GetRecoverError(err, recover())
		if err != nil {
			err = errors.NewErr(err, "goserver.InitHServer()")
		}
	}()

	hServer := env.App.HServer
	if !hServer.Enable {
		return nil, errors.New("goserver.InitHServer(): enable is false")
	}

	srcFs, ok := env.Fsm.GetFs(srcFsName)
	if !ok {
		return nil, fmt.Errorf("%s fs not exists", srcFsName)
	}

	/*	var webFs afero.Fs
		if webFsName != "" {
			webFs, ok = env.Fsm.GetFs(webFsName)
			if !ok {
				return fmt.Errorf(" %s fs not exists", webFsName)
			}
		}
	*/

	// irisApp := httpServer.App()

	fsys, err := LoadSrcFiles(srcFs, "/")
	if err != nil {
		return nil, err
	}
	loader := scriggo.Files(fsys)

	// 构建程序
	// 这里的 "main.go" 指的是 ./scripts/main.go
	program, err := scriggo.Build(loader, buildOptions)
	if err != nil {
		return nil, err
	}
	return program, err
}

// LoadSrcFiles 使用 afero 遍历 root 目录及其子目录，读取所有 .go 文件
// 返回的 map key 是文件路径，value 是文件内容
func LoadSrcFiles(fs afero.Fs, root string) (map[string][]byte, error) {
	// 初始化结果 map
	files := make(map[string][]byte)
	rootLen := len(root) - 1
	// 使用 afero.Walk 进行递归遍历
	// afero.Walk 的行为与 filepath.Walk 类似
	err := afero.Walk(fs, root, func(path string, info os.FileInfo, err error) error {
		// 1. 处理遍历过程中的错误
		if err != nil {
			return err
		}

		// 2. 跳过目录，只处理文件
		if info.IsDir() {
			return nil
		}

		// 3. 过滤文件后缀，只处理 .go 文件
		// 使用 strings.HasSuffix 或 filepath.Ext 均可
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".go") {
			return nil
		}

		// 4. 读取文件内容
		// 注意：必须使用传入的 fs 实例来读取，而不是 os 包
		content, err := afero.ReadFile(fs, path)
		if err != nil {
			return fmt.Errorf("读取文件 %s 失败: %w", path, err)
		}

		// 5. 存入 map
		// path 是相对于遍历起点的路径（如果 root 是相对路径）或绝对路径
		path = path[rootLen:]
		files[path] = content

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func AddPackages(target native.Packages, src native.Packages) {
	for key, pkg := range src {
		target[key] = pkg
	}
}

// LoadPackages
// @Description: 加载go包
// @param opts
// @return native.Packages
// @return error
func LoadPackages(opts ...func(native.Packages)) native.Packages {
	// 1. 定义包映射
	packages := native.Packages{
		// 模拟 "fmt" 包
		"fmt": native.Package{
			Name: "fmt",
			Declarations: native.Declarations{
				"Println":  fmt.Println,
				"Printf":   fmt.Printf,
				"Sprintf":  fmt.Sprintf,
				"Sprint":   fmt.Sprint,
				"Sprintln": fmt.Sprintln,
				"Fprintf":  fmt.Fprintf,
				"Fprintln": fmt.Fprintln,
			},
		},
		"strings": native.Package{
			Name: "strings",
			Declarations: native.Declarations{
				"Split":   strings.Split,
				"ToUpper": strings.ToUpper,
			},
		},
		"iris": native.Package{
			Name: "iris",
			Declarations: native.Declarations{
				"Application": (*iris.Application)(nil),
				"Context":     (*iris.Context)(nil),
			},
		},
	}
	for _, opt := range opts {
		opt(packages)
	}
	return packages
}

// LoadGlobals
// @Description: 加载全局变量
// @param opts
// @return native.Declarations
// @return error
func LoadGlobals(opts ...func(native.Declarations)) (native.Declarations, error) {
	globals := native.Declarations{
		"Sleep":       time.Sleep,
		"Duration":    time.Duration(0),
		"Millisecond": time.Millisecond,
	}
	for _, opt := range opts {
		opt(globals)
	}
	return globals, nil
}
