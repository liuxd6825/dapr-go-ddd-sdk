package restapp

import (
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/mapperutils"
)

func NewEnv(cfg *EnvConfig) *envs.Env {
	env := envs.NewEnv()
	env.Fs = cfg.Fs
	env.Name = cfg.Name

	env.App = NewApp(cfg.App)
	env.Fsm = NewFsm(cfg.Fs, cfg.App.HServer.SrcName)
	env.Log = NewLog(cfg.Log)
	env.Dapr = NewDapr(cfg.Dapr)
	env.Services = cfg.Services
	env.DBs = cfg.DBs
	env.OSS = cfg.OSS
	env.Graphdb = cfg.Graphdb
	env.IMDB = cfg.IMDB

	env.Init()
	envs.SetEnv(env)

	return env
}

func NewDapr(cfg *DaprConfig) *envs.Dapr {
	dapr := envs.NewDapr()
	copy(cfg, dapr)
	return dapr
}

/*
func NewAuth(cfg *AuthConfig) *env.Auth {
	auth := env.NewAuth()
	copy(cfg, auth)
	return auth
}
*/

func NewApp(cfg *AppConfig) *envs.App {
	app := envs.NewApp()
	copy(cfg, app)
	return app
}

func NewLog(cfg *LogConfig) *envs.Log {
	log := envs.NewLog()
	copy(cfg, log)
	return log
}

func NewFsm(fs []map[string]any, defFsName string) *fsm.Manager {
	if len(fs) != 0 {
		fsManager, err := fsm.NewManagerWithConfigs(fs, defFsName)
		if err != nil {
			panic(errors.New("fs.NewManagerWithConfigs() err: %s", err.Error()))
		}
		return fsManager
	}
	return fsm.NewManager()
}

func copy(fromObj any, toObj any) {
	err := mapperutils.Mapper(fromObj, toObj)
	if err != nil {
		panic(fmt.Sprintf("NewApp() error:%s", err.Error()))
	}
}
