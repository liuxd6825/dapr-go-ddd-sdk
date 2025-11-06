package memoryfs

import (
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/maputils/mapstructure"
)

type Config struct {
	Name string   `yaml:"name"`
	Tags []string `yaml:"tags"`
}

func NewConfig(metadata map[string]interface{}) (*Config, error) {
	cfg := &Config{}
	err := mapstructure.Decode(metadata, cfg)
	if err != nil {
		return nil, err
	}
	vErr := errors.NewVerifyError()
	vErr.Message = fmt.Sprintf("config: %s", cfg.Name)
	if cfg.Name == "" {
		vErr.AppendField("id", "missing id")
	}
	if vErr.HasError() {
		return nil, vErr
	}
	return cfg, err
}
