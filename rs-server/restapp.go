package rs_server

import (
	"errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
)

// InitHttpServer
//
//	@Description: 添加到restapp的初始化函数 Options.Init
//	@param s
//	@return error
func InitHttpServer(s *restapp.HttpServer) error {
	env := s.EnvConfig()
	if env.App.RsServer.Enable {
		fsManger, err := env.GetFsManager()
		if err != nil {
			return err
		}
		fileFs, fileOk := fsManger.Get(env.App.RsServer.FileFsId)
		if !fileOk {
			return errors.New("rsServer file fs key not found")
		}
		httpFs, httpOk := fsManger.Get(env.App.RsServer.HttpFsId)
		if !httpOk {
			return errors.New("rsServer http fs key not found")
		}
		fsCfg := NewFsConfig(fileFs, httpFs)
		server := New(s.App(), fsCfg, env.App.RsServer.Reload)
		if err := server.Run(); err != nil {
			return err
		}
	}
	return nil
}
