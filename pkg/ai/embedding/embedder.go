package embedding

import "context"

type Embedder interface {
	EmbedTexts(ctx context.Context, texts []string) ([][]float32, error)
}
