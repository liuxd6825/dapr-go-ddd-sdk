package storage

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/llm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/chunks/xlsx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/entity"
	"github.com/tmc/langchaingo/textsplitter"
	"path"
	"strings"
	"time"
)

type RagConfig struct {
	ChunkSize                  int           // 每个文本块的最大词数
	Overlap                    int           // 块之间的重叠词数
	MaxRetries                 int           // 最大重试次数
	BackoffDuration            time.Duration //补偿时间
	MaxSummariesTokenLength    int           //
	GleanCount                 int           // 调用LLM进行补偿次数
	ConcurrencyCount           int           //
	EntityExtractionPromptData *EntityExtractionPromptData
	ChunksDocument             func(llm llm.LLM, doc *entity.Document) ([]Source, error)
	BatchSize                  int
	xlsxChunk                  *xlsx.XlsxChunk
}

func NewRagConfig(opts ...func(cfg *RagConfig)) *RagConfig {
	config := &RagConfig{
		ChunkSize:               1024, // 每个文本块的最大词数
		Overlap:                 50,   // 块之间的重叠词数
		ConcurrencyCount:        5,    // 并行批量处理数
		MaxRetries:              3,    // 最大重试次数
		BackoffDuration:         60 * time.Second,
		MaxSummariesTokenLength: 1000,
		GleanCount:              3,
		BatchSize:               1,
		EntityExtractionPromptData: &EntityExtractionPromptData{
			Goal:        "提取实体与之间的关系",
			EntityTypes: []string{"人员", "公司", "产品", "机构", "合同", "交易", "案件"},
			Language:    "中文",
		},
		xlsxChunk: xlsx.NewXlsxChunk(),
	}
	config.ChunksDocument = config.getChunksDocument
	for _, opt := range opts {
		if opt != nil {
			opt(config)
		}
	}
	return config
}

func (d *RagConfig) GetBatchSize() int {
	return d.BatchSize
}

func (d *RagConfig) GetChunksDocument(llm llm.LLM, doc *entity.Document) ([]Source, error) {
	return d.ChunksDocument(llm, doc)
}

func (d *RagConfig) getChunksDocument(llm llm.LLM, doc *entity.Document) ([]Source, error) {
	if path.Ext(doc.FileName) == ".xlsx" {
		sheets, err := d.xlsxChunk.GetSheets(llm, doc.Text)
		if err != nil {
			return nil, err
		}
		res := make([]Source, 0)
		pageSize := 100
		for _, sheet := range sheets {
			rows := sheet.GetRows()
			total := len(rows)
			count := total / pageSize
			if count%pageSize > 0 {
				count++
			}
			for j := 0; j < count; j++ {
				sb := strings.Builder{}
				sb.WriteString(fmt.Sprintf("# 文件名:%s\n", doc.FileName))
				sb.WriteString(fmt.Sprintf("### Sheet名称:%s 第%d页内容 \n", sheet.SheetName, j))
				sb.WriteString(fmt.Sprintf("```csv\n"))
				for _, line := range sheet.GetFields() {
					sb.WriteString(line)
				}
				start := j * pageSize
				end := start + pageSize
				if end > total {
					end = total
				}
				pageRows := rows[start:end]
				for _, line := range pageRows {
					sb.WriteString(line)
				}
				sb.WriteString(fmt.Sprintf("```\n"))
				s := Source{
					TenantId: doc.TenantId,
					CaseId:   doc.CaseId,
					DocId:    doc.Id,
				}
				s.Content = sb.String()
				res = append(res, s)
			}

		}
		return res, nil
	}

	list, err := LangChainGoSplitText(doc.Text, d.ChunkSize, d.Overlap)
	if err != nil {
		return nil, err
	}
	res := make([]Source, 0)
	for i, text := range list {
		text = fmt.Sprintf("#文件名:%s\n##内容%d\n%s", doc.FileName, i+1, text)
		res = append(res, Source{
			TenantId:   doc.TenantId,
			CaseId:     doc.CaseId,
			DocId:      doc.Id,
			Content:    text,
			OrderIndex: i,
			TokenSize:  len(text),
		})
	}
	return res, nil
}

// GetEntityExtractionPromptData 实体提取提示数据
func (d *RagConfig) GetEntityExtractionPromptData() EntityExtractionPromptData {
	return *d.EntityExtractionPromptData
}

// GetMaxRetries 最大尝试数
func (d *RagConfig) GetMaxRetries() int {
	return d.MaxRetries
}

// GetConcurrencyCount 并发总数
func (d *RagConfig) GetConcurrencyCount() int {
	return d.ConcurrencyCount
}

// GetBackoffDuration 补偿时间
func (d *RagConfig) GetBackoffDuration() time.Duration {
	return d.BackoffDuration
}

// GetGleanCount 收集数
func (d *RagConfig) GetGleanCount() int {
	return d.GleanCount
}

// GetMaxSummariesTokenLength 最大摘要令牌长度
func (d *RagConfig) GetMaxSummariesTokenLength() int {
	return d.MaxSummariesTokenLength
}

const NewPageTag = "<!--break-->\n"
const NewPageTagLength = len(NewPageTag)

func LangChainGoSplitText(text string, chunkSize, chunkOverlap int) ([]string, error) {
	splitter := textsplitter.NewRecursiveCharacter(func(options *textsplitter.Options) {
		options.ChunkSize = chunkSize       // 每个块最大字符数
		options.ChunkOverlap = chunkOverlap // 每个块重叠字符数
	})
	index := strings.Index(text, NewPageTag)
	if index > -1 {
		var list []string
		for {
			if index > -1 {
				/*
					chunks, err := splitter.SplitText(text[:index])
					if err != nil {
						return nil, err
					}
					list = append(list, chunks...)*/
				list = append(list, text[:index])
				text = text[index+NewPageTagLength:]
			} else {
				return list, nil
			}
			index = strings.Index(text, NewPageTag)
		}
	}
	res, err := splitter.SplitText(text)
	return res, err
}
