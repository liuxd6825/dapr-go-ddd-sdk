package context

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/llm"
	"github.com/cloudwego/eino/schema"
)

type Compressor struct {
	llm    llm.LLM
	config SummaryConfig
}

func NewCompressor(llm llm.LLM, config SummaryConfig) *Compressor {
	return &Compressor{
		llm:    llm,
		config: config,
	}
}

func (c *Compressor) GenerateSummary(ctx context.Context, oldSummary string, recentHistory []*Content) (string, error) {
	if len(recentHistory) == 0 && oldSummary == "" {
		return "", nil
	}

	var prompt strings.Builder
	prompt.WriteString("【任务目标】\n")
	prompt.WriteString("你是一个记忆整理助手。请根据" + "现有摘要" + "和" + "新增对话记录" + "，生成一份更新后的摘要。\n\n")

	if oldSummary != "" {
		prompt.WriteString("【现有摘要】：\n")
		prompt.WriteString(oldSummary)
		prompt.WriteString("\n\n")
	}

	if len(recentHistory) > 0 {
		prompt.WriteString("【新增对话记录】：\n")
		for _, h := range recentHistory {
			prompt.WriteString(h.Role + ": " + h.Content + "\n")
		}
		prompt.WriteString("\n")
	}

	prompt.WriteString("【更新要求】\n")
	prompt.WriteString("1. 保持客观第三人称视角（如" + "用户表示..." + "）。\n")
	prompt.WriteString("2. 保留用户最新的偏好、已确定的事实以及未解决的疑问。\n")
	prompt.WriteString("3. 剔除无效的客套话。\n")
	prompt.WriteString("\n【请输出新的摘要】：\n")

	msgList := []*schema.Message{
		{Role: schema.User, Content: prompt.String()},
	}

	resp, err := c.llm.Generate(ctx, msgList)
	if err != nil {
		return "", fmt.Errorf("failed to generate summary: %w", err)
	}

	return resp.Content, nil
}

func (c *Compressor) ExtractEntities(ctx context.Context, history []*Content) (*ExtractedEntities, error) {
	if len(history) == 0 {
		return &ExtractedEntities{
			UserProfile: make(map[string]string),
			KeyFacts:    []string{},
			Preferences: make(map[string]any),
		}, nil
	}

	prompt := `【任务】
提取用户输入中包含的关键事实和个人偏好，以JSON格式输出。

【输出格式】
{
    "user_profile": {"role": "用户角色", "topic": "讨论主题"},
    "key_facts": ["关键事实1", "关键事实2"],
    "preferences": {}
}

【对话历史】：
`
	for _, h := range history {
		prompt += fmt.Sprintf("%s: %s\n", h.Role, h.Content)
	}

	prompt += "\n【请输出JSON】："

	msgList := []*schema.Message{
		{Role: schema.User, Content: prompt},
	}

	resp, err := c.llm.Generate(ctx, msgList)
	if err != nil {
		return nil, fmt.Errorf("failed to extract entities: %w", err)
	}

	result := &ExtractedEntities{
		UserProfile: make(map[string]string),
		KeyFacts:    []string{},
		Preferences: make(map[string]any),
	}

	jsonStr := extractJSONFromResponse(resp.Content)
	if jsonStr == "" {
		return result, nil
	}

	if err := json.Unmarshal([]byte(jsonStr), result); err != nil {
		return result, nil
	}

	return result, nil
}

func (c *Compressor) ShouldCompress(historyTokens int) bool {
	return historyTokens >= c.config.Threshold
}

func (c *Compressor) GetConfig() SummaryConfig {
	return c.config
}

func extractJSONFromResponse(content string) string {
	start := strings.Index(content, "{")
	if start == -1 {
		start = strings.Index(content, "[")
		if start == -1 {
			return ""
		}
	}

	end := strings.LastIndex(content, "}")
	if end == -1 {
		end = strings.LastIndex(content, "]")
		if end == -1 {
			return ""
		}
	}

	return content[start : end+1]
}