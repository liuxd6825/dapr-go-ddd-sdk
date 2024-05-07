package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
)

func LoadConfig(config string) (*JsServerConfig, error) {
	configFile, err := filepath.Abs(config)
	if err != nil {
		panic(err)
	}

	bs, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	var cfg JsServerConfig
	if err = yaml.Unmarshal(bs, &cfg); err != nil {
		panic(err)
	}
	return &cfg, err
}
