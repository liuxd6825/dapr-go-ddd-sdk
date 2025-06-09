package embedding

import (
	"context"
	"testing"
)

func Test_OllamaEmbedder_EmbedTexts(t *testing.T) {
	e, err := NewOllamaEmbedder(OllamaConfig{
		BaseURL:        "http://localhost:11434",
		EmbeddingModel: "bge-m3:latest", // 使用的模型版本
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	ctx := context.Background()
	texts := []string{
		"张三",
		"胡适",
		"顾炎武",
	}
	data, err := e.EmbedTexts(ctx, texts)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(data)
}
