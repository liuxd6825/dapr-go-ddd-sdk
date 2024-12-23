package pkg

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"github.com/spf13/afero"
)

type Server interface {
	GetFsPkg() FsPkg
	GetSrcFs() afero.Fs
	ReadFile(filename string, opts ...*fsopts.Options) ([]byte, error)
	GetCacheEnable() bool
	SetRunValue(key string, value any)
	GetEnvConfig() common.IEnvConfig
	GetSchemaLoader() schema.URLLoader
}
