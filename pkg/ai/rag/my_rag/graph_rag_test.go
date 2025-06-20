package my_rag

import (
	"context"
	_ "embed"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/embedding"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/entity"
	llm2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/llm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/storage"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"github.com/sirupsen/logrus"
	"testing"
)

//go:embed xtest/xyj.txt
var xyjTxt string

//go:embed xtest/graph.txt
var graphTxt string

const tenantId = "test"
const caseId = "1001"
const docId = "doc1001"

func Test_GraphRag_CreateTenant(t *testing.T) {
	ctx := context.Background()
	rag := newGraphRag(ctx, true)
	if err := rag.CreateTenant(ctx, tenantId); err != nil {
		t.Error(err)
		return
	}
	if err := rag.CreateCase(ctx, tenantId, caseId); err != nil {
		t.Error(err)
	}

	if err := rag.LoadTenant(ctx, tenantId); err != nil {
		t.Error(err)
		return
	}
}

func Test_GraphRag_DeleteTenant(t *testing.T) {
	ctx := context.Background()
	rag := newGraphRag(ctx, true)
	if err := rag.DeleteTenant(ctx, tenantId); err != nil {
		t.Error(err)
		return
	}
}

func Test_GraphRag_DeleteCase(t *testing.T) {
	ctx := context.Background()
	rag := newGraphRag(ctx, true)
	if err := rag.DeleteCase(ctx, tenantId, caseId); err != nil {
		t.Error(err)
		return
	}
}

func Test_GraphRag_DeleteDoc(t *testing.T) {
	ctx := context.Background()
	rag := newGraphRag(ctx, true)
	if err := rag.DeleteDoc(ctx, tenantId, caseId, docId); err != nil {
		t.Error(err)
		return
	}
}

func Test_GraphRag_LoadTenant(t *testing.T) {
	ctx := context.Background()
	rag := newGraphRag(ctx, true)
	if err := rag.LoadTenant(ctx, tenantId); err != nil {
		t.Error(err)
		return
	}
}

func Test_GraphRag_IngestDocuments(t *testing.T) {

	fileName := "/xtest_file/xyjTxt.txt"
	doc := &entity.Document{
		Id:       docId,
		FileName: fileName,
		Text:     xyjTxt,
		TenantId: tenantId,
		CaseId:   caseId,
	}
	ctx := context.Background()
	rag := newGraphRag(ctx, true)
	gp.Try(func() error {
		rag.IngestDocuments(ctx, []*entity.Document{doc}, storage.Options{
			TenantId: tenantId,
			CaseId:   caseId,
		})
		return nil
	}).Catch(func(e error) {
		t.Error(e)
	})
}

func Test_GraphRag_Query(t *testing.T) {
	gp.Try(func() error {
		ctx := context.Background()
		rag := newGraphRag(ctx, false)
		query := QueryParam{
			TenantId:            tenantId,
			CaseId:              caseId,
			UserPrompt:          "谁与孙悟空的公司有关系",
			ConversationHistory: []*Content{{Role: "system", Content: "详细回答相关人与公司的关司与技能"}},
		}

		res, err := rag.Query(ctx, query)
		if err == nil {
			t.Log(res)
		}
		return err
	}).Catch(func(e error) {
		t.Error(e)
	})
}

func newGraphRag(ctx context.Context, isDrop bool) *GraphRag {
	env.SetEnv(xtest.NewEnvConfig_Neo4j())
	llm := llm2.NewOpenAI(ctx, openai.ChatModelConfig{
		BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
		Model:   "deepseek-r1-distill-llama-70b", // 使用的模型版本
		APIKey:  "sk-4a999651298047efaaf38aea633ba636",
	})

	embedder := embedding.NewOllamaEmbedder(embedding.OllamaConfig{
		BaseURL:        "http://localhost:11434",
		EmbeddingModel: "bge-m3:latest",
	})

	vector := storage.NewMilvusVector(embedder, storage.MilvusConfig{
		Addr:           "127.0.0.1:19530",
		CollectionName: "tenant",
		Dim:            1024,
	})

	keyValue := storage.NewRedisKeyValueStorage()

	graph := storage.NewNeo4jGraphStorage("neo4j")
	config := NewConfigDefault(3)
	store := storage.NewStorage(graph, vector, keyValue, embedder)
	return NewGraphRag(llm, store, config, logrus.New())
}
