package vector

import (
	"context"
	"fmt"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type MilvusVector struct {
	client client.Client
	cfg    MilvusConfig
}

type MilvusConfig struct {
	Addr           string `json:"addr"`
	CollectionName string `json:"collectionName"`
	Dim            int    `json:"dim"`
}

func NewMilvusVector(cfg MilvusConfig) (*MilvusVector, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if cfg.Dim <= 0 {
		return nil, fmt.Errorf("dim should be greater than zero")
	}

	c, err := client.NewClient(ctx, client.Config{
		Address: cfg.Addr,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Milvus: %w", err)
	}

	return &MilvusVector{client: c, cfg: cfg}, nil
}

func (mc *MilvusVector) DropCollection(ctx context.Context) error {
	return mc.client.DropCollection(ctx, mc.cfg.CollectionName)
}

func (mc *MilvusVector) Init(ctx context.Context) error {
	schema := &entity.Schema{
		CollectionName: mc.cfg.CollectionName,
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
					"dim": fmt.Sprintf("%d", mc.cfg.Dim),
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
				Name:     "doc_id",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					"max_length": "256",
				},
			},
			{
				Name:     "tenant_id",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					"max_length": "256",
				},
			},
			{
				Name:     "case_id",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					"max_length": "256",
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
	index := entity.NewGenericIndex("vector_idx", entity.HNSW, map[string]string{
		"M":              "24", // HNSW参数
		"efConstruction": "200",
		"metric_type":    "L2", // 关键：明确指定L2/IP/COSINE:ml-citation{ref="7,10" data="citationList"}
	})

	err = mc.client.CreateIndex(ctx, mc.cfg.CollectionName, "vector", index, false)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	err = mc.client.LoadCollection(ctx, mc.cfg.CollectionName, false)

	return err
}

func (mc *MilvusVector) Flush(ctx context.Context) error {
	return mc.client.Flush(ctx, mc.cfg.CollectionName, false)
}

func (mc *MilvusVector) Insert(ctx context.Context, data InsertData) error {
	// 准备插入数据
	textCol := entity.NewColumnVarChar("text", data.Batch)
	sourceCol := entity.NewColumnVarChar("source", []string{data.Source})
	docIdCol := entity.NewColumnVarChar("doc_id", []string{data.DocId})
	tenantId := entity.NewColumnVarChar("tenant_id", []string{data.DocId})
	caseId := entity.NewColumnVarChar("case_id", []string{data.CaseId})
	vectorCol := entity.NewColumnFloatVector("vector", mc.cfg.Dim, data.Vectors)
	_, err := mc.client.Insert(ctx, mc.cfg.CollectionName, "", textCol, sourceCol, docIdCol, vectorCol, tenantId, caseId)
	return err
}

func (mc *MilvusVector) Search(ctx context.Context, vector []float32, topK int) ([]string, error) {
	sp, _ := entity.NewIndexHNSWSearchParam(200) // 搜索参数

	results, err := mc.client.Search(
		ctx,
		mc.cfg.CollectionName,
		[]string{},
		"",
		[]string{"text", "source"}, // 返回文本和来源
		[]entity.Vector{entity.FloatVector(vector)},
		"vector",
		entity.L2,
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

// truncateString 截断字符串到指定长度
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
