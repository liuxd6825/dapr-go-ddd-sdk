package rs_server

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/localfs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/fileutils"

	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/common"
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
	fsCfg    *FsConfig
}

func New(app *iris.Application, fsCfg *FsConfig, mainFile string, reload bool) (*JsServer, error) {
	errs := errors.NewParamsError("rs_server.New()")
	errs.AddNil("fsCfg", fsCfg)
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
		fsCfg:    fsCfg,
		reload:   reload,
	}, nil
}

func (s *JsServer) Run() error {
	if err := s.run(); err != nil {
		return err
	}
	if s.reload && s.fsCfg.FileFs.Name() == localfs.Name() {
		s.fileWatcher()
	}
	return nil
}

func (s *JsServer) run() error {
	file, err := s.fsCfg.FileFs.Open(s.mainFile)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	piState := getTestPreInitState(logrus.New())
	bundle, err := newBundle(s.mainPath, s.mainFile, data, piState, s.app, s.fsCfg)
	if err != nil {
		return err
	}

	runner, err := js.NewFromBundle(piState, bundle)
	if err != nil {
		return err
	}

	ctx, _ := context.WithCancel(context.Background())
	//defer cancel()

	vu, err := runner.NewVU(ctx, 1, 10, make(chan metrics.SampleContainer, 1), func(vu lib.VU) error {
		modules := k6.NewModules(s.app)
		for k, m := range modules {
			if err = vu.GetRuntime().Set(k, m); err != nil {
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
	return err
}

func (s *JsServer) fileWatcher() {
	f := s.fsCfg.FileFs
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

func newBundle(rootPath, filename string, jsCodeData []byte, piState *lib.TestPreInitState, app *iris.Application, fs *FsConfig) (*js.Bundle, error) {
	jsModules := map[string]any{
		"k6/server": server.New(app),
		"k6/db":     db.New(),
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
		fs.ToMap(),
		jsModules,
	)
}
