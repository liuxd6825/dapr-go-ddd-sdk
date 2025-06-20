package my_rag

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/storage"
	"github.com/tmc/langchaingo/textsplitter"
	"time"
)

type Config struct {
	ChunkSize                  int           // 每个文本块的最大词数
	Overlap                    int           // 块之间的重叠词数
	MaxRetries                 int           // 最大重试次数
	BackoffDuration            time.Duration //补偿时间
	MaxSummariesTokenLength    int           //
	GleanCount                 int           // 调用LLM进行补偿次数
	ConcurrencyCount           int           //
	EntityExtractionPromptData storage.EntityExtractionPromptData
	ChunksDocument             func(tenantId, caseId, docId string, content string) ([]storage.Source, error)
}

func NewConfigDefault(gleanCount int) *Config {
	config := &Config{
		ChunkSize:               512, // 每个文本块的最大词数
		Overlap:                 50,  // 块之间的重叠词数
		ConcurrencyCount:        2,   // 批量处理大小
		MaxRetries:              3,   // 最大重试次数
		BackoffDuration:         60 * time.Second,
		MaxSummariesTokenLength: 1000,
		GleanCount:              gleanCount,
	}
	config.ChunksDocument = config.getChunksDocument
	return config
}

func (d *Config) GetChunksDocument(tenantId, caseId, docId string, content string) ([]storage.Source, error) {
	return d.ChunksDocument(tenantId, caseId, docId, content)
}

func (d *Config) getChunksDocument(tenantId, caseId, docId string, content string) ([]storage.Source, error) {
	list, err := LangChainGoSplitText(content, d.ChunkSize, d.Overlap)
	if err != nil {
		return nil, err
	}
	res := make([]storage.Source, 0)
	for i, s := range list {
		res = append(res, storage.Source{
			TenantId:   tenantId,
			CaseId:     caseId,
			DocId:      docId,
			Content:    s,
			OrderIndex: i,
			TokenSize:  len(s),
		})
	}
	return res, nil
}

func LangChainGoSplitText(text string, chunkSize, chunkOverlap int) ([]string, error) {
	splitter := textsplitter.NewRecursiveCharacter(func(options *textsplitter.Options) {
		options.ChunkSize = chunkSize       // 每个块最大字符数
		options.ChunkOverlap = chunkOverlap // 每个块重叠字符数
	})
	res, err := splitter.SplitText(text)
	return res, err
}

// GetEntityExtractionPromptData 实体提取提示数据
func (d *Config) GetEntityExtractionPromptData() storage.EntityExtractionPromptData {
	return storage.EntityExtractionPromptData{
		Goal:        "Extract entities",
		EntityTypes: []string{"人员", "公司", "产品", "金融机构", "合同", "事件", "其他"},
		Language:    "中文",
	}
}

// GetMaxRetries 最大尝试数
func (d *Config) GetMaxRetries() int {
	return d.MaxRetries
}

// GetConcurrencyCount 并发总数
func (d *Config) GetConcurrencyCount() int {
	return d.ConcurrencyCount
}

// GetBackoffDuration 补偿时间
func (d *Config) GetBackoffDuration() time.Duration {
	return d.BackoffDuration
}

// GetGleanCount 收集数
func (d *Config) GetGleanCount() int {
	return d.GleanCount
}

// GetMaxSummariesTokenLength 最大摘要令牌长度
func (d *Config) GetMaxSummariesTokenLength() int {
	return d.MaxSummariesTokenLength
}
