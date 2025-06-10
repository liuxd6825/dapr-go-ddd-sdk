package my_rag

import (
	"fmt"
	"github.com/cloudwego/eino/schema"
	"io"
	"strings"
	"time"
)

func Reader(reader *schema.StreamReader[*schema.Message], streams ...func(txt string)) (*strings.Builder, error) {
	content := &strings.Builder{}
	defer reader.Close()
	for {
		chunk, err := reader.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			// 错误处理
			return nil, err
		}
		if len(streams) > 0 {
			stream := streams[len(streams)-1]
			stream(chunk.Content)
		}
		// 响应片段处理
		content.WriteString(chunk.Content)
	}
	return content, nil
}

// GenerateDocumentID 生成文档ID
func GenerateDocumentID(source string) string {
	if source == "" {
		return fmt.Sprintf("doc-%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%s-%d", source, time.Now().UnixNano())
}

func ChunkText(text string, chunkSize int, overlap int) []string {
	var chunks []string
	words := strings.Fields(text)
	totalWords := len(words)

	if totalWords == 0 {
		return chunks
	}

	start := 0
	for start < totalWords {
		end := start + chunkSize
		if end > totalWords {
			end = totalWords
		}
		chunks = append(chunks, strings.Join(words[start:end], " "))

		if end == totalWords {
			break
		}

		// 重叠处理
		start = end - overlap
		if start < 0 {
			start = 0
		}
	}

	return chunks
}

// truncateString 截断字符串到指定长度
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
