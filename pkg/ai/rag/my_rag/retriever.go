package my_rag

import (
	"context"
	"fmt"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type Retriever struct {
	milvusClient client.Client
	collection   string
}

func NewRetriever(milvusClient client.Client, collection string) *Retriever {
	return &Retriever{
		milvusClient: milvusClient,
		collection:   collection,
	}
}

func (r *Retriever) Search(ctx context.Context, vector []float32, topK int) ([]string, error) {
	sp, _ := entity.NewIndexIvfFlatSearchParam(16) // 搜索参数

	results, err := r.milvusClient.Search(
		ctx,
		r.collection,
		[]string{},
		"",
		[]string{"text"},
		[]entity.Vector{entity.FloatVector(vector)},
		"vector",
		entity.COSINE,
		topK,
		sp,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search Milvus: %w", err)
	}

	texts := make([]string, 0, topK)
	for _, result := range results {
		for _, field := range result.Fields {
			if field.Name() == "text" {
				text, ok := field.(*entity.ColumnVarChar)
				if !ok {
					continue
				}
				for i := 0; i < text.Len(); i++ {
					val, _ := text.ValueByIdx(i)
					texts = append(texts, val)
				}
			}
		}
	}

	return texts, nil
}
