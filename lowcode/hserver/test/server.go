package test

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/test"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/spf13/afero"
)

type TestServer struct {
	FsPkg  pkg.FsPkg
	SrcFs  afero.Fs
	EnvCfg common.IEnvConfig
}

type Options func(*TestServer)

func NewServer(srcPath string, opts ...Options) (pkg.Server, error) {
	fs, err := test.NewFsManager()
	if err != nil {
		return nil, err
	}

	envCfg := NewEnvConfig(fs, srcPath)
	server, err := NewTestServer(envCfg)
	for _, opt := range opts {
		opt(server)
	}

	if err != nil {

		return nil, err
	}
	return server, nil
}

func (s *TestServer) GetEnvConfig() common.IEnvConfig {
	return s.EnvCfg
}

func NewTestServer(cfg common.IEnvConfig) (*TestServer, error) {
	return &TestServer{
		FsPkg:  nil,
		SrcFs:  afero.NewOsFs(),
		EnvCfg: cfg,
	}, nil
}

func (s *TestServer) GetFsPkg() pkg.FsPkg {
	return s.FsPkg
}

func (s *TestServer) GetSrcFs() afero.Fs {
	return s.SrcFs
}

func (s *TestServer) ReadFile(filename string, opts ...*fsopts.Options) ([]byte, error) {
	data := s.FsPkg.ReadFile(filename, opts...)
	return data, nil
}

func (s *TestServer) GetCacheEnable() bool {
	return false
}

func (s *TestServer) SetRunValue(key string, value any) {

}
