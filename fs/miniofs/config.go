package miniofs

import "github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils/mapstructure"

type Config struct {
	Name     string `yaml:"name"`
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
	return cfg, err
}
