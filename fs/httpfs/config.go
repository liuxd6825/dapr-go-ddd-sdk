package httpfs

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils/mapstructure"
)

type Config struct {
	Id      string `json:"id"`
	BaseUrl string `json:"baseUrl"`
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
		vErr.AppendField("id", "missing id")
	}
	if cfg.BaseUrl == "" {
		vErr.AppendField("baseUrl", "missing baseUrl")
	}
	if vErr.HasError() {
		return nil, vErr
	}
	return cfg, err
}
