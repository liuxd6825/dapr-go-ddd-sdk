package rs_server

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rs-server/modules/k6"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rs-server/modules/k6/db"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rs-server/modules/k6/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rs-server/modules/k6/server"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rs-server/watcher"
	"github.com/liuxd6825/k6server/js"
	"github.com/liuxd6825/k6server/lib"
	"github.com/liuxd6825/k6server/lib/fsext"
	"github.com/liuxd6825/k6server/loader"
	"github.com/liuxd6825/k6server/metrics"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/url"
	"os/exec"
	"strings"
)

type JsServer struct {
	app     *iris.Application
	srcPath string
	reload  bool
	watcher watcher.Watcher
}

func New(app *iris.Application, srcPath string, reload bool) *JsServer {
	return &JsServer{
		app:     app,
		srcPath: srcPath,
		reload:  reload,
	}
}

func (s *JsServer) Run() error {
	if err := s.run(); err != nil {
		return err
	}
	if s.reload {
		s.fileWatcher()
	}
	return nil
}

func (s *JsServer) run() error {
	mainFile := s.srcPath + "/main.js"
	data, err := ioutil.ReadFile(mainFile)
	if err != nil {
		return err
	}

	piState := getTestPreInitState(logrus.New())

	bundle, err := newBundle(s.srcPath, mainFile, data, piState, s.app)
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
			vu.GetRuntime().Set(k, m)
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
	s.watcher = watcher.NewFileWatcher(s.srcPath)
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

			s.app.Logger().Infof("restart JsServer, file: %s/%s。", rootPath, fileName)
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

func newBundle(rootPath string, filename string, data []byte, piState *lib.TestPreInitState, app *iris.Application) (*js.Bundle, error) {
	fs := fsext.NewOsFs()
	jsModules := map[string]any{
		"k6/server": server.New(app),
		"k6/db":     db.New(),
		"k6/schema": schema.New(),
	}

	return js.NewBundleFormJsModules(
		piState,
		&loader.SourceData{
			URL:  &url.URL{Path: filename, Scheme: "file"},
			Data: data,
			PWD:  &url.URL{Path: rootPath, Scheme: "file"},
		},
		map[string]fsext.Fs{"file": fs, "https": fsext.NewMemMapFs()},
		jsModules,
	)
}
