package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/embedding"
	"sync"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/llm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/storage"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/sirupsen/logrus"
)

type RagService struct {
	graphRag           *my_rag.GraphRag
	messageService     *MessageService
	chatSummaryService *ChatSummaryService
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

func (s *RagService) SetMessageService(messageService *MessageService) {
	s.messageService = messageService
}

func (s *RagService) SetChatSummaryService(chatSummaryService *ChatSummaryService) {
	s.chatSummaryService = chatSummaryService
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

	if s.chatSummaryService != nil && s.messageService != nil && query.ChatId != "" {
		chatSummary, err := s.chatSummaryService.GetByChatId(ctx, query.ChatId)
		if err == nil && chatSummary != nil && chatSummary.Summary != "" {
			s.graphRag.SetInitialSummary(chatSummary.Summary)
		}

		chatHistory, err := s.messageService.FindByChatIdAfterPosition(ctx, query.ChatId, 0)
		if err == nil && len(chatHistory) > 0 {
			history := make([]*my_rag.Content, 0, len(chatHistory))
			for _, msg := range chatHistory {
				history = append(history, &my_rag.Content{
					Role:    msg.Role,
					Content: msg.Content,
				})
			}
			query.ConversationHistory = history
		}
	}

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

	var llmModel model.ToolCallingChatModel
	switch ragCfg.LLM.Type {
	case "ollama":
		timeout := 30 * time.Second
		m, err := llm.NewOllama(ctx, ollama.ChatModelConfig{
			BaseURL: ragCfg.LLM.BaseUrl,
			Model:   ragCfg.LLM.Model, // 使用的模型版本
			//Thinking: &ollama.ThinkValue{"low"},
			Timeout: timeout,
		})

		if err != nil {
			panic("open rag model error" + err.Error())
		}
		llmModel = m
	case "ark":
		m, err := llm.NewArk(ctx, &ark.ChatModelConfig{
			BaseURL: ragCfg.LLM.BaseUrl,
			APIKey:  ragCfg.LLM.APIKey, // 填写你的 API Key
			Model:   ragCfg.LLM.Model,  // 填写你的推理接入点 ID (例如: ep-2024...)
			// Region: "cn-beijing",          // 可选，默认为 cn-beijing
		})
		if err != nil {
			panic("open rag model error" + err.Error())
		}
		llmModel = m
	case "openai":
		m, err := llm.NewOpenAI(ctx, openai.ChatModelConfig{
			BaseURL: ragCfg.LLM.BaseUrl,
			Model:   ragCfg.LLM.Model, // 使用的模型版本
			APIKey:  ragCfg.LLM.APIKey,
		})

		if err != nil {
			panic("open rag model error" + err.Error())
		}

		llmModel = m
	default:
		panic("rag config llm.type is null ")
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

	//store := storage.NewStorageNoEmbedder(graphStorage, kv)
	store := storage.NewStorage(graphStorage, kv, vectorStorage, embedder)
	return my_rag.NewGraphRag(llmModel, store, ragConfig, logrus.New())
}
