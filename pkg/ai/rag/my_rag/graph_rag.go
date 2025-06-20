package my_rag

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/entity"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/llm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/storage"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/sirupsen/logrus"
	"log"
	"strings"
	"sync"
	"time"
)

type GraphRag struct {
	LLM    llm.LLM
	store  storage.Storage
	logger *logrus.Logger
	config storage.Config
}

const MaxRetrieveContexts = 5 // 最大检索上下文数量

func NewGraphRag(llm llm.LLM, store storage.Storage, config storage.Config, logger *logrus.Logger) *GraphRag {
	return &GraphRag{
		LLM:    llm,
		store:  store,
		logger: logger,
		config: config,
	}
}

// IngestDocuments 将文档摄取到 Milvus
func (g *GraphRag) IngestDocuments(ctx context.Context, docs []*entity.Document, opts storage.Options) ([]string, int, []string) {
	const (
		batchSize = 100 // 批量处理大小
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
			docId = GenerateDocumentID(doc.FileName)
		}
		documentIDs = append(documentIDs, docId)

		// 分块处理文本
		chunks, err := g.config.GetChunksDocument(doc.TenantId, doc.CaseId, doc.Id, doc.Text)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}
		if len(chunks) == 0 {
			errors = append(errors, fmt.Sprintf("文档 %s 没有有效内容", docId))
			continue
		}

		if err := storage.InsertDocument(ctx, doc, g.config, g.store, g.LLM, g.logger); err != nil {
			errors = append(errors, fmt.Sprintf("导入文档 %s 时出错%s", docId, err.Error()))
			continue
		}

		fileName := truncateString(doc.FileName, 256)
		// 批量处理文本块
		for i := 0; i < len(chunks); i += batchSize {
			end := i + batchSize
			if end > len(chunks) {
				end = len(chunks)
			}
			batch := chunks[i:end]

			var strList []string
			for _, source := range batch {
				strList = append(strList, source.Content)
			}

			insertData := storage.InsertDocData{
				ChunkId:  i,
				TenantId: doc.TenantId,
				CaseId:   doc.CaseId,
				FileName: fileName,
				DocId:    doc.Id,
				Content:  strList,
			}

			maxRetries := g.config.GetMaxRetries()
			// 插入数据到 Milvus（带重试）
			for attempt := 1; attempt <= maxRetries; attempt++ {
				err := g.store.VectorInsertDoc(ctx, insertData, opts)
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
	if err := g.store.VectorFlush(ctx, opts); err != nil {
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
	qry := storage.GraphQueryParam{
		Keys:    nodeKeys,
		MaxDeep: query.MaxDeep,
		Limit:   query.TopK,
	}
	opt := storage.Options{
		TenantId: query.TenantId,
		CaseId:   query.CaseId,
	}
	graphContext, err := g.store.GraphQuery(ctx, qry, opt)
	if err != nil {
		return &GoResult{Err: errors.New("取图知识时出错：%s", err.Error())}
	}
	return &GoResult{
		Data: graphContext,
	}
}

func (g *GraphRag) getDocumentContext(ctx context.Context, query *QueryParam) *GoResult {
	opts := storage.Options{
		TenantId: query.TenantId,
		CaseId:   query.CaseId,
	}

	queryEmbed, err := g.store.EmbedTexts(ctx, []string{query.Query})
	if err != nil {
		return &GoResult{Err: errors.New("将查询内容转为向量数据时出错：%s", err.Error())}
	}

	contexts, err := g.store.VectorSearchDoc(ctx, queryEmbed[0], query.TopK, opts)
	if err != nil {
		return &GoResult{Err: errors.New("查找向量数据时出错：%s", err.Error())}
	}

	return &GoResult{
		Data: contexts,
	}
}

func (g *GraphRag) CreateTenant(ctx context.Context, tenantId string) error {
	return g.store.CreateTenant(ctx, tenantId)
}

func (g *GraphRag) CreateCase(ctx context.Context, tenantId, caseId string) error {
	return g.store.CreateCase(ctx, tenantId, caseId)
}

func (g *GraphRag) DeleteTenant(ctx context.Context, tenantId string) error {
	return g.store.DeleteTenant(ctx, tenantId)
}

func (g *GraphRag) DeleteCase(ctx context.Context, tenantId, caseId string) error {
	return g.store.DeleteCase(ctx, tenantId, caseId)
}

func (g *GraphRag) DeleteDoc(ctx context.Context, tenantId, caseId, docId string) error {
	return g.store.DeleteDoc(ctx, tenantId, caseId, docId)
}

func (g *GraphRag) LoadTenant(ctx context.Context, tenantId string) error {
	return g.store.LoadTenant(ctx, tenantId)
}

func (g *GraphRag) Query(ctx context.Context, query QueryParam, streams ...func(txt string)) (string, error) {
	if query.Query == "" {
		return "", errors.New("query is empty")
	}
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

	contexts := []string{query.Query}
	for res := range resultCh {
		if res.Err != nil {
			return "", fmt.Errorf("协程执行失败: %w", res.Err)
		}
		contexts = append(contexts, res.Data...)
	}

	logs.Info(ctx, logs.Fields{"contexts": contexts})

	prompt := buildRAGPrompt(query.Query, contexts)
	messages := []*schema.Message{
		{Role: schema.User, Content: prompt},
	}
	for _, item := range query.ConversationHistory {
		messages = append(messages, &schema.Message{
			Role:    schema.RoleType(item.Role),
			Content: item.Content,
		})
	}

	resp, err := g.LLM.Stream(ctx, messages)
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
	resp, err := g.LLM.Stream(ctx, msgList)
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
