package config

import "encoding/json"

type OryConfig struct {
	Url string `yaml:"url"`
}

func ReadOryConfig(ragMeta any) (*OryConfig, error) {
	jsonData, err := json.Marshal(ragMeta)
	if err != nil {
		return nil, err
	}
	var oryCfg OryConfig
	if err := json.Unmarshal(jsonData, &oryCfg); err != nil {
		return nil, err
	}
	return &oryCfg, nil
}
