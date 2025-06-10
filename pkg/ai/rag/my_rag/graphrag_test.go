package my_rag

import (
	"context"
	_ "embed"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/embedding"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/entity"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/graph"
	llm2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/llm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/vector"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/randomutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"testing"
)

//go:embed xtest/xyj.txt
var xyjTxt string

func Test_GraphRag_IngestDocuments(t *testing.T) {
	doc := &entity.Document{
		Id:       randomutils.NewId(),
		Source:   "/xtest_file/xyj.txt",
		Text:     xyjTxt,
		TenantId: "test",
		CaseId:   "1001",
	}
	ctx := context.Background()
	rag := newGraphRag(ctx, true)
	gp.Try(func() error {
		rag.IngestDocuments(ctx, []*entity.Document{doc})
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
			TenantId:     "test",
			CaseId:       "1001",
			Query:        "谁与孙悟空的公司有关系",
			SystemPrompt: "详细回答相关人与公司的关司与技能",
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
	if isDrop {
		if err = vectorStorage.DropCollection(ctx); err != nil {
			panic(err)
		}
	}

	if err = vectorStorage.Init(ctx); err != nil {
		panic(err)
	}

	graphStorage := graph.NewNeo4jGraphStorage()
	return NewGraphRag(llm, embedder, vectorStorage, graphStorage)
}
