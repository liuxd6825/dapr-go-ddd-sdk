package embedding

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
)

type OllamaConfig struct {
	BaseURL        string `json:"baseUrl"`
	ApiKey         string `json:"apiKey"`
	ApiVersion     string `json:"apiVersion"`
	EmbeddingModel string `json:"embeddingModel"`
}
type OllamaEmbedder struct {
	embedder embeddings.Embedder
}

func NewOllamaEmbedder(cfg OllamaConfig) (*OllamaEmbedder, error) {
	llm, err := ollama.New(
		ollama.WithServerURL(cfg.BaseURL),
		ollama.WithModel(cfg.EmbeddingModel),
		ollama.WithRunnerVocabOnly(true),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM: %w", err)
	}
	e, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedder: %w", err)
	}
	return &OllamaEmbedder{embedder: e}, nil
}

func (e *OllamaEmbedder) EmbedTexts(ctx context.Context, texts []string) ([][]float32, error) {
	vectors, err := e.embedder.EmbedDocuments(ctx, texts)
	if err != nil {
		return nil, fmt.Errorf("failed to embed documents: %w", err)
	}

	// 转换为 float32 切片
	return lo.Map(vectors, func(v []float32, _ int) []float32 {
		return lo.Map(v, func(f float32, _ int) float32 {
			return float32(f)
		})
	}), nil
}
