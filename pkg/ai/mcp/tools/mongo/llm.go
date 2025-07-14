package mongo

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	llm2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/llm"
)

func newLLM(ctx context.Context, cfg *LLMConfig) (llm llm2.LLM, err error) {
	switch cfg.Type {
	case LLMType_Ollama:
		llm, err = llm2.NewOllama(ctx, ollama.ChatModelConfig{
			BaseURL: cfg.BaseURL,
			Model:   cfg.Model,
		})
	case LLMType_OpenAi:
		llm, err = llm2.NewOpenAI(ctx, openai.ChatModelConfig{
			BaseURL: cfg.BaseURL,
			Model:   cfg.Model,
			APIKey:  cfg.ApiKey,
		})
	}
	return llm, err
}
