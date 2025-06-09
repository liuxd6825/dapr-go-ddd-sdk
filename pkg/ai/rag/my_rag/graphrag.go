package my_rag

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/embedding"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/entity"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/llm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/vector"
	"log"
	"strings"
	"time"
)

type GraphRag struct {
	llm      llm.LLM
	embedder embedding.Embedder
	vector   vector.VectorStorage
}

const MaxRetrieveContexts = 3 // 最大检索上下文数量

func NewGraphRag(llm llm.LLM, embedder embedding.Embedder, vector vector.VectorStorage) *GraphRag {
	return &GraphRag{
		llm:      llm,
		embedder: embedder,
		vector:   vector,
	}
}

// IngestDocuments 将文档摄取到 Milvus
func (g *GraphRag) IngestDocuments(ctx context.Context, docs []*entity.Document) ([]string, int, []string) {
	const (
		chunkSize  = 512 // 每个文本块的最大词数
		overlap    = 50  // 块之间的重叠词数
		batchSize  = 100 // 批量处理大小
		maxRetries = 3   // 最大重试次数
	)

	var (
		documentIDs []string
		chunkCount  int
		errors      []string
	)

	// 处理所有文档
	for _, doc := range docs {
		// 生成文档ID（如果未提供）
		docId := doc.Id
		if docId == "" {
			docId = GenerateDocumentID(doc.Source)
		}
		documentIDs = append(documentIDs, docId)

		// 分块处理文本
		chunks := ChunkText(doc.Text, chunkSize, overlap)
		if len(chunks) == 0 {
			errors = append(errors, fmt.Sprintf("文档 %s 没有有效内容", docId))
			continue
		}

		// 批量处理文本块
		for i := 0; i < len(chunks); i += batchSize {
			end := i + batchSize
			if end > len(chunks) {
				end = len(chunks)
			}
			batch := chunks[i:end]

			// 向量化文本块
			vectors, err := g.embedder.EmbedTexts(ctx, batch)
			if err != nil {
				errors = append(errors, fmt.Sprintf("文档 %s 块 %d-%d 向量化失败: %v", docId, i, end, err))
				continue
			}

			source := truncateString(doc.Source, 256)
			insertData := vector.InsertData{
				TenantId: doc.TenantId,
				CaseId:   doc.CaseId,
				Source:   source,
				DocId:    doc.Id,
				Batch:    batch,
				Vectors:  vectors,
			}
			// 插入数据到 Milvus（带重试）
			for attempt := 1; attempt <= maxRetries; attempt++ {
				err := g.vector.Insert(ctx, insertData)

				if err == nil {
					chunkCount += len(batch)
					break
				}

				if attempt == maxRetries {
					errors = append(errors, fmt.Sprintf("文档 %s 块 %d-%d 插入失败: %v", docId, i, end, err))
				} else {
					log.Printf("插入失败（尝试 %d/%d），重试中...: %v", attempt, maxRetries, err)
					time.Sleep(time.Duration(attempt) * time.Second)
				}
			}
		}
	}

	// 刷新数据确保可搜索
	if err := g.vector.Flush(ctx); err != nil {
		errors = append(errors, fmt.Sprintf("刷新失败: %v", err))
	}

	return documentIDs, chunkCount, errors
}

func (g *GraphRag) Query(ctx context.Context, query string, context []string) (string, error) {
	queryEmded, err := g.embedder.EmbedTexts(ctx, []string{query})
	if err != nil {
		return "", err
	}

	contexts, err := g.vector.Search(ctx, queryEmded[0], MaxRetrieveContexts)
	if err != nil {
		return "", err
	}

	prompt := buildRAGPrompt(query, append(context, contexts...))

	resp, err := g.llm.Stream(ctx, []*schema.Message{
		{Role: schema.User, Content: prompt},
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	sb, err := Reader(resp)
	if err != nil {
		return "", fmt.Errorf("no response from LLM")
	}

	return sb.String(), nil
}

func buildRAGPrompt(query string, context []string) string {
	sb := strings.Builder{}
	sb.WriteString("你是一个知识助手，根据以下上下文回答问题：\n\n")

	for i, text := range context {
		sb.WriteString(fmt.Sprintf("上下文 %d: %s\n\n", i+1, text))
	}

	sb.WriteString(fmt.Sprintf("问题: %s\n\n回答:", query))
	return sb.String()
}
