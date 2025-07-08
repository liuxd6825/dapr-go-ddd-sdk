package storage

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/doc_extract"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/embedding"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"testing"
)

func Test_VectorInsertDoc(t *testing.T) {
	logger := logrus.New()
	extract := doc_extract.NewExtract(logger)
	fs := afero.NewOsFs()
	fileName := "/Users/lxd/Projects/liuxd6825/dapr/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/xtest/天眼查-刘建新.pdf"
	txt, err := extract.Extract(fs, fileName)
	if err != nil {
		t.Error(err)
		return
	}

	embedder := embedding.NewOllamaEmbedder(embedding.OllamaConfig{
		BaseURL:        "http://192.168.120.224:11434",
		EmbeddingModel: "bge-m3:latest",
	})
	vector := NewMilvusVector(embedder, MilvusConfig{
		Addr:           "192.168.120.224:19530",
		CollectionName: "tenant",
		Dim:            1024,
	})
	ctx := context.Background()
	data := InsertDocData{
		CaseId:   "1001",
		TenantId: "test",
		ChunkId:  0,
		FileName: "a.txt",
		DocId:    "docId",
		Content:  []string{txt},
	}
	err = vector.VectorInsertDoc(ctx, data, Options{})
	if err != nil {
		t.Error(err)
		return
	}
}
