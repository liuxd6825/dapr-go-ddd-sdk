package localfs

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils/mapstructure"
)

type Config struct {
	Id   string `yaml:"id"`
	Path string `yaml:"path"`
}

func NewConfig(metadata map[string]any) (*Config, error) {
	cfg := &Config{}
	err := mapstructure.Decode(metadata, cfg)
	if err != nil {
		return nil, err
	}
	vErr := errors.NewVerifyError()
	vErr.Message = fmt.Sprintf("config: %s", cfg.Id)
	if cfg.Id == "" {
		vErr.AppendField("id", "missing name")
	}
	if cfg.Path == "" {
		vErr.AppendField("path", "missing url")
	}
	if vErr.HasError() {
		return nil, vErr
	}
	return cfg, err
}
