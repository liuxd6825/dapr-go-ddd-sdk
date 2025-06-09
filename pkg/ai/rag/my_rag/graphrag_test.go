package my_rag

import (
	"context"
	_ "embed"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/embedding"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/entity"
	llm2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/llm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/vector"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/randomutils"
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
	ctx := context.Background()
	rag := newGraphRag(ctx, false)
	gp.Try(func() error {
		res, err := rag.Query(ctx, "孙悟空是谁", []string{})
		if err == nil {
			t.Log(res)
		}
		return err
	}).Catch(func(e error) {
		t.Error(e)
	})
}

func newGraphRag(ctx context.Context, isDrop bool) *GraphRag {
	llm, err := llm2.NewOllama(ctx, ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "modelscope.cn/Qwen/Qwen3-14B-GGUF:latest", // 使用的模型版本
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
	return NewGraphRag(llm, embedder, vectorStorage)
}
