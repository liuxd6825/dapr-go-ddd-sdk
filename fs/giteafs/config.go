package giteafs

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils/mapstructure"
)

type Config struct {
	Id       string `yaml:"id"`
	Url      string `yaml:"url"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Email    string `yaml:"email"`
	Repo     string `yaml:"repos"`
	Branch   string `yaml:"branch"`
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
	if cfg.Url == "" {
		vErr.AppendField("url", "missing url")
	}
	if cfg.User == "" {
		vErr.AppendField("user", "missing user")
	}
	if cfg.Password == "" {
		vErr.AppendField("password", "missing password")
	}
	if cfg.Email == "" {
		vErr.AppendField("email", "missing email")
	}
	if cfg.Repo == "" {
		vErr.AppendField("repo", "missing repo")
	}
	if cfg.Branch == "" {
		vErr.AppendField("branch", "missing branch")
	}
	if vErr.HasError() {
		return nil, vErr
	}
	return cfg, err
}
