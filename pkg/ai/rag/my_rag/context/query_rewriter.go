package context

import (
	"context"
	"fmt"
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/llm"
	"github.com/cloudwego/eino/schema"
)

type QueryRewriter struct {
	llm llm.LLM
}

func NewQueryRewriter(llm llm.LLM) *QueryRewriter {
	return &QueryRewriter{llm: llm}
}

func (qr *QueryRewriter) Rewrite(ctx context.Context, currentQuery string, history []*Content) (string, error) {
	if len(history) == 0 {
		return currentQuery, nil
	}

	historyText := ""
	for _, h := range history {
		historyText += fmt.Sprintf("%s: %s\n", h.Role, h.Content)
	}

	prompt := fmt.Sprintf(`根据对话历史，将当前问题改写为一个独立、完整的查询。

对话历史：
%s

当前问题：%s

改写要求：
1. 将指代词（如"它"、"这个"）替换为具体内容
2. 补全省略的主语和宾语
3. 确保问题可以独立理解

只输出改写后的查询语句，不要其他内容：`, historyText, currentQuery)

	msgList := []*schema.Message{
		{Role: schema.User, Content: prompt},
	}

	resp, err := qr.llm.Generate(ctx, msgList)
	if err != nil {
		return currentQuery, fmt.Errorf("failed to rewrite query: %w", err)
	}

	rewritten := strings.TrimSpace(resp.Content)
	if rewritten == "" {
		return currentQuery, nil
	}

	return rewritten, nil
}

func (qr *QueryRewriter) RewriteStandalone(ctx context.Context, query string) (string, error) {
	prompt := fmt.Sprintf(`将以下用户问题改写为一个独立、完整的查询语句。

要求：
1. 补全被省略的主语和宾语
2. 将指代词替换为具体内容
3. 确保问题可以独立理解，不依赖对话历史

原始问题：%s

只输出改写后的查询语句：`, query)

	msgList := []*schema.Message{
		{Role: schema.User, Content: prompt},
	}

	resp, err := qr.llm.Generate(ctx, msgList)
	if err != nil {
		return query, fmt.Errorf("failed to rewrite query: %w", err)
	}

	rewritten := strings.TrimSpace(resp.Content)
	if rewritten == "" {
		return query, nil
	}

	return rewritten, nil
}