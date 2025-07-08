package service

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/embedding"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/llm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/storage"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/sirupsen/logrus"
	"sync"
)

type RagService struct {
	graphRag *my_rag.GraphRag
}

var ragService *RagService
var ragServiceOnce sync.Once
var graphRag *my_rag.GraphRag
var graphRagOnce sync.Once

func NewRagService() *RagService {
	ragServiceOnce.Do(func() {
		ragService = &RagService{
			graphRag: NewGraphRag(),
		}
	})
	return ragService
}

func (s *RagService) CreateTenant(ctx context.Context) error {
	tenantId, _ := appctx.GetTenantId(ctx)
	return s.graphRag.CreateTenant(ctx, tenantId)
}

func (s *RagService) CreateCase(ctx context.Context, caseId string) error {
	tenantId, _ := appctx.GetTenantId(ctx)
	return s.graphRag.CreateCase(ctx, tenantId, caseId)
}

func (s *RagService) Query(ctx context.Context, query my_rag.QueryParam, streams ...func(txt string)) (string, error) {
	query.TenantId, _ = appctx.GetTenantId(ctx)
	return s.graphRag.Query(ctx, &query, streams...)
}

func NewGraphRag() *my_rag.GraphRag {
	graphRagOnce.Do(func() {
		graphRag = newGraphRag()
	})
	return graphRag
}

func newGraphRag() *my_rag.GraphRag {
	ctx := context.Background()
	e := env.GetEnv()
	ragMeta := e.App.Meta["rag"]
	if ragMeta == nil {
		panic("rag not found in env.app")
	}

	ragCfg, err := config.ReadRagConfig(ragMeta)
	if err != nil {
		panic("read RagConfig error" + err.Error())
	}

	llm, err := llm.NewOpenAI(ctx, openai.ChatModelConfig{
		BaseURL: ragCfg.LLM.BaseUrl,
		Model:   ragCfg.LLM.Model, // 使用的模型版本
		APIKey:  ragCfg.LLM.APIKey,
	})

	if err != nil {
		panic("open rag model error" + err.Error())
	}

	embedder := embedding.NewOllamaEmbedder(embedding.OllamaConfig{
		BaseURL:        ragCfg.Embedder.BaseURL,
		EmbeddingModel: ragCfg.Embedder.Model,
		ApiKey:         ragCfg.Embedder.ApiKey,
	})

	vectorStorage := storage.NewMilvusVector(embedder, storage.MilvusConfig{
		Addr:           ragCfg.Vector.Addr,
		CollectionName: ragCfg.Vector.CollectionName,
		Dim:            ragCfg.Vector.Dim,
	})

	graphStorage := storage.NewNeo4jGraphStorage("neo4j", logs.GetLogger())
	kv := storage.NewRedisKeyValueStorage()

	ragConfig := storage.NewRagConfig(func(cfg *storage.RagConfig) {

	})
	store := storage.NewStorage(graphStorage, vectorStorage, kv, embedder)
	return my_rag.NewGraphRag(llm, store, ragConfig, logrus.New())
}
