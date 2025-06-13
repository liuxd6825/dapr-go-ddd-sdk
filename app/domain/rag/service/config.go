package service

import "encoding/json"

type RagConfig struct {
	LLM      LLMConfig      `yaml:"llm"`
	Embedder EmbedderConfig `yaml:"embedder"`
	Vector   VectorConfig   `yaml:"vector"`
}

type LLMConfig struct {
	Model   string `json:"model"`
	APIKey  string `json:"apiKey"`
	BaseUrl string `json:"baseUrl"`
}

type EmbedderConfig struct {
	BaseURL string `json:"baseUrl"`
	Model   string `json:"model"`
	ApiKey  string `json:"apiKey"`
}

type VectorConfig struct {
	Addr           string `json:"addr"`
	CollectionName string `json:"collectionName"`
	Dim            int    `json:"dim"`
}

func ReadRagConfig(ragMeta any) (*RagConfig, error) {
	jsonData, err := json.Marshal(ragMeta)
	if err != nil {
		return nil, err
	}
	var ragCfg RagConfig
	if err := json.Unmarshal(jsonData, &ragCfg); err != nil {
		return nil, err
	}
	return &ragCfg, nil
}
