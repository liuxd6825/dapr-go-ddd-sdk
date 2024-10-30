package rs_server

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/localfs"
	common2 "github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/fileutils"

	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/db"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/server"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/watcher"
	"github.com/liuxd6825/k6server/js"
	"github.com/liuxd6825/k6server/lib"
	"github.com/liuxd6825/k6server/loader"
	"github.com/liuxd6825/k6server/metrics"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/url"
	"os/exec"
	"strings"
)

type JsServer struct {
	app      *iris.Application
	mainFile string
	mainPath string
	reload   bool
	watcher  watcher.Watcher
	srcFs    *SrcFsConfig // source code file system
	data     map[string]any
	envCfg   common2.IEnvConfig
}

func NewServer(app *iris.Application, data map[string]any, srcFs *SrcFsConfig, envCfg common2.IEnvConfig, mainFile string, reload bool) (*JsServer, error) {
	errs := errors.NewParamsError("rs_server.NewServer()")
	errs.AddNil("srcFs", srcFs)
	errs.AddNil("app", app)
	errs.AddNil("mainFile", mainFile)

	if errs.HasError() {
		return nil, errs
	}

	var mainPath = "/"
	i := strings.LastIndex(mainFile, "/")
	if i > -1 {
		mainPath = mainFile[:i]
	}

	return &JsServer{
		app:      app,
		mainPath: mainPath,
		mainFile: mainFile,
		srcFs:    srcFs,
		reload:   reload,
		data:     data,
		envCfg:   envCfg,
	}, nil
}

type RunOption = func(vu lib.VU) error

func (s *JsServer) Run(options ...RunOption) error {
	if err := s.run(options...); err != nil {
		return err
	}
	if s.reload && s.srcFs.FileFs.Name() == localfs.Name() {
		s.fileWatcher()
	}
	return nil
}

func (s *JsServer) run(options ...RunOption) error {
	file, err := s.srcFs.FileFs.Open(s.mainFile)
	if err != nil {
		return err
	}

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	/*fsManager, err := s.envCfg.GetFsManager()
	if err != nil {
		return err
	}*/

	piState := getTestPreInitState(logrus.New())
	bundle, err := newBundle(s.mainPath, s.mainFile, fileData, piState, s.app, s.data, s.srcFs, s.envCfg)
	if err != nil {
		return err
	}

	runner, err := js.NewFromBundle(piState, bundle)
	if err != nil {
		return err
	}

	ctx, _ := context.WithCancel(context.Background())
	//defer cancel()

	vu, err := runner.NewVU(ctx, 1, 1, make(chan metrics.SampleContainer, 1), func(vu lib.VU) error {
		modules := k6.NewModules(s.app, s.data, s.envCfg)
		for k, m := range modules {
			if err = vu.GetRuntime().Set(k, m); err != nil {
				return err
			}
		}
		for _, opt := range options {
			if err := opt(vu); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	params := &lib.VUActivationParams{RunContext: ctx}
	err = vu.Activate(params).RunOnce()
	if scriptErr, ok := err.(*js.ScriptExceptionError); ok {
		for _, v := range scriptErr.Inner().Stack() {
			return NewCodeError(err, &v)
		}
	}
	return err
}

type CodeError struct {
	fileName string
	lines    []string
	position int
	message  string
}

func NewCodeError(err error, v *goja.StackFrame) *CodeError {
	program := v.Program()
	var errList []string
	errList = append(errList, err.Error())
	pos := v.Position().Line
	if pos != 0 {
		pos += 1
		lines := strings.Split(program.Src().Source(), "\n")
		start := pos - 10
		end := pos + 20
		if end > len(lines) {
			end = len(lines)
			start = end - 30
		}
		if start < 0 {
			start = 0
		}
		for i := start; i < end; i++ {
			var line string
			if i == pos {
				line = fmt.Sprintf("=>%03d : ", i)
			} else {
				line = fmt.Sprintf("  %03d : ", i)
			}
			line += lines[i]
			errList = append(errList, line)
		}
	}
	fileName := ""
	if program != nil && program.Src() != nil {
		fileName = program.Src().Name()
	}
	return &CodeError{fileName: fileName, lines: errList, position: pos}
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("%s:%d\n", e.fileName, e.position) + strings.Join(e.lines, "\n")
}

func (s *JsServer) fileWatcher() {
	f := s.srcFs.FileFs
	if f.Name() != localfs.Name() {
		return
	}
	rootPath := s.mainPath
	if pathFs, ok := f.(fs.PathFs); ok {
		rootPath = fileutils.AbsPath(pathFs.BasePath(), rootPath)
	}

	s.watcher = watcher.NewFileWatcher(rootPath)
	err := s.watcher.Start(func(rootPath, fileName string, eventType watcher.EventType) error {
		reload := false
		if strings.HasSuffix(fileName, ".ts") {
			s.tsc()
			return nil
		}

		switch eventType {
		case watcher.WriteEvent:
			reload = true
		default:
			return nil
		}
		if reload {
			if err := s.run(); err != nil {
				s.app.Logger().Error(err)
			}
			if err := s.app.RefreshRouter(); err != nil {
				s.app.Logger().Error(err)
			}

			s.app.Logger().Infof("restart RsServer, file: %s/%s。", rootPath, fileName)
		}
		return nil
	})
	if err != nil {
		s.app.Logger().Error(err)
	}
}

func (s *JsServer) tsc() {
	cmd := exec.Command("tsc", "--build", "tsconfig.json")
	err := cmd.Run()
	if err != nil {
		fmt.Println("tsc failed:" + err.Error())
	}
}

func getTestPreInitState(tb logrus.FieldLogger) *lib.TestPreInitState {
	reg := metrics.NewRegistry()
	return &lib.TestPreInitState{
		Logger:         tb,
		RuntimeOptions: lib.RuntimeOptions{},
		Registry:       reg,
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(reg),
	}
}

func newBundle(rootPath, filename string, jsCodeData []byte, piState *lib.TestPreInitState, app *iris.Application, data map[string]any, srcFs *SrcFsConfig, cfg common2.IEnvConfig) (*js.Bundle, error) {
	jsModules := map[string]any{
		"k6/server": server.New(app, data, cfg),
		"k6/db":     db.New(cfg),
		"k6/schema": schema.New(),
		"k6/common": common.New(),
	}

	return js.NewBundleFormJsModules(
		piState,
		&loader.SourceData{
			Data: jsCodeData,
			URL:  &url.URL{Path: filename, Scheme: "file"},
			PWD:  &url.URL{Path: rootPath, Scheme: "file"},
		},
		srcFs.ToMap(),
		jsModules,
	)
}
