package my_rag

import (
	"context"
	_ "embed"
	"fmt"
	"testing"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/doc_extract"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/embedding"
	llm2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/llm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/entity"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/storage"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

//go:embed xtest/xyj.txt
var xyjTxt string

//go:embed xtest/graph.txt
var graphTxt string

const tenantId = "test"
const caseId = "1001"
const docId = "D001"

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

func Test_LangChainGoSplitText(t *testing.T) {
	logger := logrus.New()
	extract := doc_extract.NewExtract(logger)
	fs := afero.NewOsFs()
	fileName := "/Users/lxd/Projects/liuxd6825/dapr/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/xtest/天眼查-刘建新.pdf"
	txt, err := extract.Extract(fs, fileName)
	if err != nil {
		t.Error(err)
		return
	}

	list, err := storage.LangChainGoSplitText(txt, 500, 20)
	assert.NoError(t, err)
	t.Log(list)
}

func Test_GraphRag_SaveGraph(t *testing.T) {
	logger := logrus.New()
	extract := doc_extract.NewExtract(logger)
	fs := afero.NewOsFs()
	fileName := "/Users/lxd/Projects/liuxd6825/dapr/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/xtest/天眼查-刘建新.pdf"
	txt, err := extract.Extract(fs, fileName)
	if err != nil {
		t.Error(err)
		return
	}
	doc := &entity.Document{
		Id:       docId,
		FileName: "天眼查-刘建新.pdf",
		Text:     txt,
		TenantId: tenantId,
		CaseId:   caseId,
	}
	ctx := context.Background()
	rag := newGraphRag(ctx, true)
	gp.Try(func() error {
		return rag.SaveGraph(ctx, doc)
	}).Catch(func(e error) {
		t.Error(e)
	})
}

func Test_GraphRag_IngestDocuments(t *testing.T) {
	ctx := context.Background()
	rag := newGraphRag(ctx, true)
	gp.Try(func() error {
		var files []string
		//files = append(files, "刘建新-董监高对外投资及任职报告.pdf")
		files = append(files, "张宇-董监高对外投资及任职报告.pdf")
		docs, err := getDocs(t, "pdf", files)
		if err != nil {
			return err
		}
		rag.IngestDocuments(ctx, docs)
		return nil
	}).Catch(func(e error) {
		t.Error(e)
	})
}

func Test_GraphRag_IngestXlsx(t *testing.T) {
	ctx := context.Background()
	rag := newGraphRag(ctx, true)
	gp.Try(func() error {
		var files []string
		//files = append(files, "刘建新-董监高对外投资及任职报告.pdf")
		files = append(files, "赵蕾流水汇总.xlsx")
		docs, err := getDocs(t, "xlsx", files)
		if err != nil {
			return err
		}
		rag.IngestDocuments(ctx, docs)
		return nil
	}).Catch(func(e error) {
		t.Error(e)
	})
}

func getDocs(t *testing.T, docIdType string, files []string) ([]*entity.Document, error) {
	logger := logrus.New()
	extract := doc_extract.NewExtract(logger)
	var docs []*entity.Document
	fs := afero.NewOsFs()
	for i, fileName := range files {
		pathName := "/Users/lxd/Projects/liuxd6825/dapr/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/xtest/" + fileName
		txt, err := extract.Extract(fs, pathName)
		if err != nil {
			t.Error(err)
			return nil, err
		}
		doc := &entity.Document{
			Id:       fmt.Sprintf("%s%d", docIdType, i),
			FileName: fileName,
			Text:     txt,
			TenantId: tenantId,
			CaseId:   caseId,
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

func Test_GraphRag_Query(t *testing.T) {
	gp.Try(func() error {
		ctx := context.Background()
		rag := newGraphRag(ctx, false)
		query := NewQueryParam()
		query.TenantId = tenantId
		query.CaseId = caseId
		query.Query = "悟空，唐僧，女王之间发生了什么？"
		//query.ConversationHistory = []*Content{{Role: "system", Content: "详细回答相关人与公司的关司与技能"}},

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
	env.SetEnv(xtest.NewEnvConfig_Neo4j("192.168.120.224"))

	llm := newLLM2(ctx)

	embedder := embedding.NewOllamaEmbedder(embedding.OllamaConfig{
		BaseURL:        "http://192.168.120.224:11434",
		EmbeddingModel: "bge-m3:latest",
	})

	vector := storage.NewMilvusVector(embedder, storage.MilvusConfig{
		Addr:           "192.168.120.224:19530",
		CollectionName: "tenant",
		Dim:            1024,
	})
	logger := logrus.New()
	keyValue := storage.NewRedisKeyValueStorage()

	graph := storage.NewNeo4jGraphStorage("neo4j", logger)
	config := storage.NewRagConfig()
	store := storage.NewStorage(graph, vector, keyValue, embedder)
	return NewGraphRag(llm, store, config, logger)
}

func newLLM1(ctx context.Context) llm2.LLM {
	llm, err := llm2.NewOpenAI(ctx, openai.ChatModelConfig{
		BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
		Model:   "deepseek-r1-distill-llama-70b", // 使用的模型版本
		APIKey:  "sk-4a999651298047efaaf38aea633ba636",
	})
	if err != nil {
		panic(err)
	}
	return llm
}

func newLLM2(ctx context.Context) llm2.LLM {
	llm, err := llm2.NewOllama(ctx, ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "modelscope.cn/unsloth/DeepSeek-R1-Distill-Qwen-7B-GGUF:latest",
	})
	if err != nil {
		panic(err)
	}
	llm.WithTools()
	return llm
}
