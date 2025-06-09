package embedding

import (
	"context"
	"testing"
)

func Test_OpenAIEmbedder_EmbedTexts(t *testing.T) {
	e, err := NewOpenAIEmbedder(OpenAIConfig{
		BaseURL:        "http://localhost:11434/api",
		EmbeddingModel: "bge-m3:latest", // 使用的模型版本
		ApiVersion:     "",
		ApiKey:         "NONE",
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
