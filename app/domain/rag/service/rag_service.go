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
	llm, err := llm2.NewOpenAI(ctx, openai.ChatModelConfig{
		BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
		Model:   "deepseek-r1-distill-llama-70b", // 使用的模型版本
		APIKey:  "sk-4a999651298047efaaf38aea633ba636",
	})

	if err != nil {
		panic(err)
	}
	embedder, err := embedding.NewOllamaEmbedder(embedding.OllamaConfig{
		BaseURL:        "http://localhost:11434",
		EmbeddingModel: "bge-m3:latest",
	})

	if err != nil {
		panic(err)
	}
	vectorStorage, err := vector.NewMilvusVector(vector.MilvusConfig{
		Addr:           "127.0.0.1:19530",
		CollectionName: "graph",
		Dim:            1024,
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
