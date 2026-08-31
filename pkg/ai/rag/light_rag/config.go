package light_rag

import (
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
)

type LightRagConfig struct {
	Url    string `yaml:"url"`
	ApiKey string `yaml:"apiKey"`
}

func LoadLightRagConfig() *LightRagConfig {
	e := env.GetEnv()
	ragMeta := e.App.Meta["lightRag"]
	var err error
	if ragMeta == nil {
		panic("lightRag not found in env.app ")
	}

	cfg, err := readLightRagConfig(ragMeta)
	if err != nil {
		panic("read LightRagConfig error" + err.Error())
	}
	if cfg.Url == "" {
		panic("LightRAG Url can't be empty")
	}
	if cfg.ApiKey == "" {
		panic("LightRAG ApiKey can't be empty")
	}
	return cfg
}

func readLightRagConfig(ragMeta any) (*LightRagConfig, error) {
	jsonData, err := json.Marshal(ragMeta)
	if err != nil {
		return nil, err
	}
	var ragCfg LightRagConfig
	if err := json.Unmarshal(jsonData, &ragCfg); err != nil {
		return nil, err
	}
	return &ragCfg, nil
}
