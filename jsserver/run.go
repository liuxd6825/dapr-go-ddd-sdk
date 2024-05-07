package jsserver

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/jsserver/modules/k6"
	"github.com/liuxd6825/dapr-go-ddd-sdk/jsserver/modules/k6/server"
	"github.com/liuxd6825/k6server/js"
	"github.com/liuxd6825/k6server/lib"
	"github.com/liuxd6825/k6server/lib/fsext"
	"github.com/liuxd6825/k6server/loader"
	"github.com/liuxd6825/k6server/metrics"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/url"
)

const script = `
export default function() {
	server.get("/test", function(ictx){
		ictx.json({"name":"test", "address":"cc"})
	})
}
`

func RunServer(app *iris.Application, srcPath string) error {
	mainFile := srcPath + "/main.js"
	data, err := ioutil.ReadFile(mainFile)
	if err != nil {
		return err
	}

	piState := getTestPreInitState(logrus.New())

	bundle, err := newBundle(srcPath, mainFile, data, piState, app)
	if err != nil {
		return err
	}

	runner, err := js.NewFromBundle(piState, bundle)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	ctx, _ := context.WithCancel(context.Background())
	//defer cancel()

	vu, err := runner.NewVU(ctx, 1, 10, make(chan metrics.SampleContainer, 1), func(vu lib.VU) error {
		modules := k6.NewModules(app)
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
