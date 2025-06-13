package service

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/embedding"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/graph"
	llm2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/llm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/vector"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
)

type RagService struct {
	graphRag *my_rag.GraphRag
}

func NewRagService() *RagService {
	return &RagService{
		graphRag: newGraphRag(),
	}
}

func (s *RagService) Query(ctx context.Context, query my_rag.QueryParam, streams ...func(txt string)) (string, error) {
	query.TenantId, _ = appctx.GetTenantId(ctx)
	return s.graphRag.Query(ctx, query, streams...)
}

func newGraphRag() *my_rag.GraphRag {
	ctx := context.Background()
	e := env.GetEnv()
	ragMeta := e.App.Meta["rag"]
	if ragMeta == nil {
		panic("rag not found in env.app")
	}

	ragCfg, err := ReadRagConfig(ragMeta)
	if err != nil {
		panic("read RagConfig error" + err.Error())
	}

	llm, err := llm2.NewOpenAI(ctx, openai.ChatModelConfig{
		BaseURL: ragCfg.LLM.BaseUrl,
		Model:   ragCfg.LLM.Model, // 使用的模型版本
		APIKey:  ragCfg.LLM.APIKey,
	})

	if err != nil {
		panic(err)
	}
	embedder, err := embedding.NewOllamaEmbedder(embedding.OllamaConfig{
		BaseURL:        ragCfg.Embedder.BaseURL,
		EmbeddingModel: ragCfg.Embedder.Model,
		ApiKey:         ragCfg.Embedder.ApiKey,
	})

	if err != nil {
		panic(err)
	}
	vectorStorage, err := vector.NewMilvusVector(vector.MilvusConfig{
		Addr:           ragCfg.Vector.Addr,
		CollectionName: ragCfg.Vector.CollectionName,
		Dim:            ragCfg.Vector.Dim,
	})

	if err != nil {
		panic(err)
	}

	if err = vectorStorage.Init(ctx); err != nil {
		panic(err)
	}

	graphStorage := graph.NewNeo4jGraphStorage()
	return my_rag.NewGraphRag(llm, embedder, vectorStorage, graphStorage)
}
