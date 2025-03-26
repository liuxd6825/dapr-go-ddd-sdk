package test

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/utils/schema_utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/test"
	"github.com/spf13/afero"
)

type TestServer struct {
	FsPkg        pkg.FsPkg
	SrcFs        afero.Fs
	EnvCfg       common.IEnvConfig
	SchemaLoader *schema_utils.SchemaLoader
}

type Options func(*TestServer) error

func NewServer(rootPath string, opts ...Options) (pkg.Server, error) {
	fs, err := test.NewFsManager(rootPath)
	if err != nil {
		return nil, err
	}

	envCfg := NewEnvConfig(fs, rootPath)
	server, err := NewTestServer(envCfg)
	if err != nil {
		return nil, err
	}
	server.SchemaLoader = schema_utils.NewSchemaLoader(fs)
	for _, opt := range opts {
		if err = opt(server); err != nil {
			return nil, err
		}
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

func (s *TestServer) GetSchemaLoader() schema.URLLoader {
	return s.SchemaLoader
}
