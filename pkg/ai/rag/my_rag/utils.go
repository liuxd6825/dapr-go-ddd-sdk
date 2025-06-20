package my_rag

import (
	"fmt"
	"github.com/cloudwego/eino/schema"
	"io"
	"slices"
	"strings"
	"text/template"
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

func newMessages(messages []string) []*schema.Message {
	list := make([]*schema.Message, len(messages))
	for i, message := range messages {
		list[i] = &schema.Message{
			Role:    "user",
			Content: message,
		}
	}
	return list
}

func CleanContent(content string) string {
	// Removes spaces and null characters.
	str := strings.TrimSpace(content)
	return strings.ReplaceAll(str, "\x00", "")
}

func PromptTemplate(name, templ string, data any) (string, error) {
	buf := strings.Builder{}
	tmpl := template.New(name).Funcs(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	})
	tmpl = template.Must(tmpl.Parse(templ))
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func AppendIfUnique(slice []string, item string) []string {
	if slices.Contains(slice, item) {
		return slice
	}
	return append(slice, item)
}

func MostFrequentItem(list []string) string {
	// Create a map to store counts
	counts := make(map[string]int)

	// Count occurrences of each string
	for _, item := range list {
		counts[item]++
	}

	// Find the item with highest count
	maxCount := 0
	var mostFreqItem string

	for item, count := range counts {
		if count > maxCount {
			maxCount = count
			mostFreqItem = item
		}
	}

	return mostFreqItem
}

func ThreeBacktick(caption string) string {
	return "```" + caption
}
