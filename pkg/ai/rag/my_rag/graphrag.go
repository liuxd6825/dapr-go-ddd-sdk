package my_rag

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/embedding"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/entity"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/graph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/llm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/vector"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"log"
	"strings"
	"sync"
	"time"
)

type GraphRag struct {
	llm        llm.LLM
	embedder   embedding.Embedder
	vector     vector.VectorStorage
	graphStore graph.Storage
}

const MaxRetrieveContexts = 5 // 最大检索上下文数量

func NewGraphRag(llm llm.LLM, embedder embedding.Embedder, vector vector.VectorStorage, graphStore graph.Storage) *GraphRag {
	return &GraphRag{
		llm:        llm,
		embedder:   embedder,
		vector:     vector,
		graphStore: graphStore,
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

type GoResult struct {
	Data []string
	Err  error
}

func (g *GraphRag) getGraphContext(ctx context.Context, query *QueryParam) *GoResult {
	nodeKeys, err := g.getKeys(ctx, query.Query)
	if err != nil {
		return &GoResult{Err: errors.New("获取查询关键字时出错：%s", err.Error())}
	}
	logs.Info(ctx, logs.Fields{"keys": nodeKeys})
	graphContext, err := g.graphStore.GetKnowledge(ctx, query.TenantId, query.CaseId, nodeKeys, query.MaxDeep, query.Limit)
	if err != nil {
		return &GoResult{Err: errors.New("取图知识时出错：%s", err.Error())}
	}
	return &GoResult{
		Data: graphContext,
	}
}

func (g *GraphRag) getDocumentContext(ctx context.Context, query *QueryParam) *GoResult {
	queryEmbed, err := g.embedder.EmbedTexts(ctx, []string{query.Query})
	if err != nil {
		return &GoResult{Err: errors.New("将查询内容转为向量数据时出错：%s", err.Error())}
	}
	contexts, err := g.vector.Search(ctx, queryEmbed[0], MaxRetrieveContexts)
	if err != nil {
		return &GoResult{Err: errors.New("查找向量数据时出错：%s", err.Error())}
	}
	return &GoResult{
		Data: contexts,
	}
}

type QueryParam struct {
	TenantId     string   `json:"tenantId"`
	CaseId       string   `json:"caseId"`
	Query        string   `json:"query"`
	SystemPrompt string   `json:"systemPrompt"`
	Context      []string `json:"context"`
	MaxDeep      int      `json:"maxDeep"`
	Limit        int      `json:"limit"`
}

func (g *GraphRag) Query(ctx context.Context, query QueryParam, streams ...func(txt string)) (string, error) {
	if query.MaxDeep <= 0 {
		query.MaxDeep = MaxRetrieveContexts
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// 使用结构体通道传递结果和错误
	resultCh := make(chan *GoResult, 2)

	logs.Info(ctx, logs.Fields{"query": query.Query})
	// 协程1：取图关系中知识
	go func() {
		defer wg.Done()
		resultCh <- g.getGraphContext(ctx, &query)
	}()

	// 协程2：取向量数据库中的知道
	go func() {
		defer wg.Done()
		resultCh <- g.getDocumentContext(ctx, &query)
	}()

	// 等待协程完成并关闭通道
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	contexts := query.Context
	for res := range resultCh {
		if res.Err != nil {
			return "", fmt.Errorf("协程执行失败: %w", res.Err)
		}
		contexts = append(contexts, res.Data...)
	}

	logs.Info(ctx, logs.Fields{"contexts": contexts})

	prompt := buildRAGPrompt(query.Query, contexts)
	msgList := []*schema.Message{
		{Role: schema.User, Content: prompt},
	}
	if query.SystemPrompt != "" {
		msgList = append(msgList, &schema.Message{
			Role:    schema.System,
			Content: query.SystemPrompt,
		})
		logs.Info(ctx, logs.Fields{"systemPrompt": query.SystemPrompt})
	}

	resp, err := g.llm.Stream(ctx, msgList)
	if err != nil {
		return "", fmt.Errorf("大模型问题分析时出错: %w", err)
	}

	sb, err := Reader(resp, streams...)
	if err != nil {
		return "", fmt.Errorf("读取返回结果时出错: %w", err)
	}

	return sb.String(), nil
}

// getKeys 取得关键字
func (g *GraphRag) getKeys(ctx context.Context, query string) ([]string, error) {
	msgList := []*schema.Message{
		{Role: schema.User, Content: fmt.Sprintf("### 指令：提取以下文本中的实体并用逗号分隔\n### 文本：{%s}", query)},
		{Role: schema.System, Content: `
			从文本中提取所有人名、地名和组织名，用英文逗号分隔各实体，不要包含其他符号或说明。\n
			按以下格式输出：实体1,实体2,实体3
			示例：张三,李四,北京,上海
		`},
	}
	resp, err := g.llm.Stream(ctx, msgList)
	sb, err := Reader(resp)
	if err != nil {
		return nil, fmt.Errorf("no response from LLM")
	}
	str := sb.String()
	if index := strings.Index(str, "</think>"); index != -1 {
		str = str[index+8:]
		str = strings.Replace(str, " ", "", -1)
		str = strings.Replace(str, "\n", "", -1)
	}
	list := strings.Split(str, ",")
	return list, nil
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
