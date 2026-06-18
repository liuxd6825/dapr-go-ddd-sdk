package miniofs

import (
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/maputils/mapstructure"
)

type Config struct {
	Name     string   `yaml:"name"`
	Bucket   string   `yaml:"bucket"`
	RootPath string   `yaml:"rootPath"`
	Tags     []string `yaml:"tags"`
}

func NewConfig(metadata map[string]any) (*Config, error) {
	cfg := &Config{}
	err := mapstructure.Decode(metadata, cfg)
	if err != nil {
		return nil, err
	}
	vErr := errors.NewVerifyError()
	vErr.Message = fmt.Sprintf("miniofs.config: %s", cfg.Name)
	if cfg.Name == "" {
		vErr.AppendField("name", "missing name")
	}
	if cfg.Bucket == "" {
		vErr.AppendField("bucket", "missing bucket")
	}
	if vErr.HasError() {
		return nil, vErr
	}
	return cfg, err
}
