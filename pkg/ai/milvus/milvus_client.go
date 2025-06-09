package milvus

import (
	"context"
	"fmt"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type MilvusClient struct {
	client client.Client
}

func NewMilvusClient(addr string) (*MilvusClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c, err := client.NewClient(ctx, client.Config{
		Address: addr,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Milvus: %w", err)
	}

	return &MilvusClient{client: c}, nil
}

func (mc *MilvusClient) CreateCollection(ctx context.Context, collectionName string, dim int) error {
	schema := &entity.Schema{
		CollectionName: collectionName,
		Description:    "RAG Knowledge Base",
		Fields: []*entity.Field{
			{
				Name:       "id",
				DataType:   entity.FieldTypeInt64,
				PrimaryKey: true,
				AutoID:     true,
			},
			{
				Name:     "vector",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": fmt.Sprintf("%d", dim),
				},
			},
			{
				Name:     "text",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					"max_length": "65535",
				},
			},
			{
				Name:     "source",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					"max_length": "256",
				},
			},
		},
	}

	err := mc.client.CreateCollection(ctx, schema, entity.DefaultShardNumber)
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	// 创建索引
	index := entity.NewGenericIndex("vector_idx", entity.IvfFlat, map[string]string{
		"nlist": "1024",
	})
	err = mc.client.CreateIndex(ctx, collectionName, "vector", index, false)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}
