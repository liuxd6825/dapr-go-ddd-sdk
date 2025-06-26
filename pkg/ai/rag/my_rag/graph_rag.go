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
	"strings"
	"sync"
	"time"
)

type GraphRag struct {
	LLM       llm.LLM
	store     storage.Storage
	logger    *logrus.Logger
	config    storage.Config
	docHandle *storage.GraphHandle
}

const MaxRetrieveContexts = 5 // 最大检索上下文数量

type GoResult struct {
	Data []string
	Err  error
}

func NewGraphRag(llm llm.LLM, store storage.Storage, config storage.Config, logger *logrus.Logger) *GraphRag {
	docHandle := storage.NewGraphHandle(config, store, llm, logger)
	return &GraphRag{
		LLM:       llm,
		store:     store,
		logger:    logger,
		config:    config,
		docHandle: docHandle,
	}
}

func (g *GraphRag) SetOnEvent() {

}

func (g *GraphRag) SetOnEvents(setEvents func(e *storage.DocEvents)) {
	g.docHandle.SetOnEvents(setEvents)
}

// IngestDocuments 将文档摄取到 Milvus
func (g *GraphRag) IngestDocuments(ctx context.Context, docs []*entity.Document) (documentIDs []string, chunkCount int, errorList []error) {
	// 处理所有文档
	for _, doc := range docs {
		documentIDs = append(documentIDs, doc.Id)
		chunkCountVal, err := g.IngestDocument(ctx, doc)
		if err != nil {
			errorList = append(errorList, err)
		}
		chunkCount += chunkCountVal
	}
	return documentIDs, chunkCount, errorList
}

// IngestDocument
func (g *GraphRag) IngestDocument(ctx context.Context, doc *entity.Document) (chunkCount int, err error) {
	// 生成文档ID（如果未提供）
	docId := doc.Id
	if docId == "" {
		docId = GenerateDocumentID(doc.FileName)
	}

	// 分块处理文本
	chunks, err := g.config.GetChunksDocument(doc.TenantId, doc.CaseId, doc.Id, doc.Text)
	if err != nil {
		return 0, err
	}
	chunkCount = len(chunks)
	if chunkCount == 0 {
		return chunkCount, errors.New("文档 %s 没有有效内容", doc.FileName)
	}

	if chunkCount, err = g.SaveVector(ctx, doc); err != nil {
		err = errors.New("导入文档%s生成向量数据时出错, %s。", doc.FileName, err.Error())
	} else if err = g.docHandle.SaveGraph(ctx, doc); err != nil {
		err = errors.New("导入文档%s生成图数据时出错, %s。", doc.FileName, err.Error())
	}

	if err != nil {
		if delErr := g.DeleteDoc(ctx, doc.TenantId, doc.CaseId, doc.Id); delErr != nil {
			err = errors.New("%s 撤销文件%s时失败:%s", err.Error(), doc.Id, delErr.Error())
		}
	}
	return chunkCount, err
}

func (g *GraphRag) SaveVector(ctx context.Context, doc *entity.Document) (chunkCount int, err error) {
	// 分块处理文本
	chunks, err := g.config.GetChunksDocument(doc.TenantId, doc.CaseId, doc.Id, doc.Text)
	if err != nil {
		return 0, err
	}
	chunkCount = len(chunks)
	if chunkCount == 0 {
		return chunkCount, errors.New("文档 %s 没有有效内容", doc.FileName)
	}
	logger := g.logger
	batchSize := g.config.GetBatchSize()
	fileName := truncateString(doc.FileName, 1000)
	for i := 0; i < chunkCount; i += batchSize {
		end := i + batchSize
		if end > chunkCount {
			end = chunkCount
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
		opts := storage.Options{
			TenantId: doc.TenantId,
			CaseId:   doc.CaseId,
		}
		// 插入数据到 Milvus（带重试）
		for attempt := 1; attempt <= maxRetries; attempt++ {
			err = g.store.VectorInsertDoc(ctx, insertData, opts)
			if err == nil {
				break
			}

			if attempt == maxRetries {
				return 0, errors.New("文档 %s 块 %d-%d 插入失败: %v", doc.FileName, i, end, err)
			} else {
				logger.Printf("插入失败（尝试 %d/%d），重试中...: %v", attempt, maxRetries, err)
				time.Sleep(time.Duration(attempt) * time.Second)
			}
		}
	}
	return
}

// getGraphContext 取得图数据上下文
func (g *GraphRag) getGraphContext(ctx context.Context, query *QueryParam) *GoResult {
	nodeKeywords, err := g.getKeywords(ctx, query.Query)
	if err != nil {
		return &GoResult{Err: errors.New("获取查询关键字时出错：%s", err.Error())}
	}
	logs.Info(ctx, logs.Fields{"node keywords": nodeKeywords})
	qry := storage.GraphQueryParam{
		Keys:    nodeKeywords,
		MaxDeep: query.MaxDeep,
		Limit:   query.TopK,
	}
	opt := storage.Options{
		TenantId:  query.TenantId,
		CaseId:    query.CaseId,
		NodeLabel: storage.NodeLabel_All,
	}
	graphContext, err := g.store.GraphQuery(ctx, qry, opt)
	if err != nil {
		return &GoResult{Err: errors.New("取图知识时出错：%s", err.Error())}
	}
	return &GoResult{
		Data: graphContext,
	}
}

// getVectorContext 取得向量数据上下文
func (g *GraphRag) getVectorContext(ctx context.Context, query *QueryParam) *GoResult {
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

func (g *GraphRag) Query(ctx context.Context, query *QueryParam, streams ...func(txt string)) (string, error) {
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

	logger := g.logger
	logger.Info("query", query.Query)

	// 协程1：取图关系中知识
	go func() {
		defer wg.Done()
		resultCh <- g.getGraphContext(ctx, query)
	}()

	// 协程2：取向量数据库中的知道
	go func() {
		defer wg.Done()
		resultCh <- g.getVectorContext(ctx, query)
	}()

	// 等待协程完成并关闭通道
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var contexts []string
	for res := range resultCh {
		if res.Err != nil {
			return "", fmt.Errorf("协程执行失败: %w", res.Err)
		}
		contexts = append(contexts, res.Data...)
	}

	logger.Info("contexts", contexts)

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
func (g *GraphRag) getKeywords(ctx context.Context, query string) ([]string, error) {
	msgList := []*schema.Message{
		{Role: schema.User, Content: fmt.Sprintf("### 指令：提取以下文本中的实体并用逗号分隔\n### 文本：{%s}", query)},
		{Role: schema.System, Content: `
			从文本中提取所有实体名称，用英文逗号分隔各实体，不要包含其他符号或说明。\n
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
