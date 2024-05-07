package server

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/html/template/loader"
	"github.com/liuxd6825/dapr-go-ddd-sdk/html/template/loader/file"
	"github.com/liuxd6825/dapr-go-ddd-sdk/html/template/loader/gitea"
	"github.com/liuxd6825/dapr-go-ddd-sdk/power/server/config"
)

func Start(app *iris.Application, cfg *config.JsServerConfig, settings ...Setting) error {
	if cfg == nil {
		return errors.New("cfg  is nil")
	}

	if cfg.Service != nil {
		if cfg.Service.Html != nil && enable(cfg.Service.Html.Enable) {
			cfg.Init()
			if err := InitConfig(cfg.Repos); err != nil {
				return err
			}
			url := cfg.Service.Html.Url
			if len(url) == 0 {
				return errors.New("power.servivce.html.url is nil")
			}
			if err := AddHtmlHandler(app, url); err != nil {
				return err
			}
		}
		if cfg.Service.Api != nil && enable(cfg.Service.Api.Enable) {
			path := cfg.Service.Api.Path
			if len(path) == 0 {
				return errors.New("power.servivce.api.path is nil")
			}
			if err := AddServerHandler(app, path, settings...); err != nil {
				return err
			}
		}
	}

	return nil
}

func enable(val *bool) bool {
	if val == nil {
		return true
	}
	b := *val
	return b
}

func InitConfig(cfg *config.Repos) (err error) {
	for _, p := range cfg.Gitea {
		l := gitea.NewLoader(p.Name, p.Repo, p.Branch)
		l.Init("gitea:"+p.Name, p.Repo, p.Branch)
		if err = l.Connect(p.Url, p.User, p.Password); err != nil {
			return err
		}
		if err = loader.AddLoader(l); err != nil {
			return err
		}
	}
	for _, p := range cfg.File {
		l := file.NewLoader("file:"+p.Name, p.Path)
		if err = l.Connect(); err != nil {
			return err
		}
		if err = loader.AddLoader(l); err != nil {
			return err
		}
	}
	return err
}
