package localfs

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils/mapstructure"
)

type Config struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

func NewConfig(metadata map[string]any) (*Config, error) {
	cfg := &Config{}
	err := mapstructure.Decode(metadata, cfg)
	if err != nil {
		return nil, err
	}
	vErr := errors.NewVerifyError()
	vErr.Message = fmt.Sprintf("localfs.config: %s", cfg.Name)
	if cfg.Name == "" {
		vErr.AppendField("id", "missing id")
	}
	if cfg.Path == "" {
		vErr.AppendField("path", "missing url")
	}
	if vErr.HasError() {
		return nil, vErr
	}
	return cfg, err
}
