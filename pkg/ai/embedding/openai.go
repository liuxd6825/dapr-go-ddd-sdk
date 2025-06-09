package embedding

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
)

type OpenAIConfig struct {
	BaseURL        string `json:"baseUrl"`
	ApiKey         string `json:"apiKey"`
	ApiVersion     string `json:"apiVersion"`
	EmbeddingModel string `json:"embeddingModel"`
}
type OpenAIEmbedder struct {
	embedder embeddings.Embedder
}

func NewOpenAIEmbedder(cfg OpenAIConfig) (*OpenAIEmbedder, error) {
	llm, err := openai.New(
		openai.WithBaseURL(cfg.BaseURL),
		openai.WithAPIVersion(cfg.ApiVersion),
		openai.WithEmbeddingModel(cfg.EmbeddingModel),
		openai.WithToken(cfg.ApiKey),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM: %w", err)
	}
	e, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedder: %w", err)
	}
	return &OpenAIEmbedder{embedder: e}, nil
}

func (e *OpenAIEmbedder) EmbedTexts(ctx context.Context, texts []string) ([][]float32, error) {
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
