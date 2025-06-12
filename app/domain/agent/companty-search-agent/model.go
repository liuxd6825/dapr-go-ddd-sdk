package companty_search_agent

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

// newDeepSeekR1ForAli 阿里云
func newDeepSeekR1ForAli(ctx context.Context, model string) (cm model.ToolCallingChatModel, err error) {
	apiKey := "sk-4a999651298047efaaf38aea633ba636"

	config := &openai.ChatModelConfig{
		BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
		Model:   model,
		APIKey:  apiKey,
	}
	cm, err = openai.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}
