package neo4jagent

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

// newChatModel component initialization function of node 'ChatModel1' in graph 'graphAgent'
func newText2CypherChatModel(ctx context.Context) (cm model.ToolCallingChatModel, err error) {
	config := &openai.ChatModelConfig{
		BaseURL: "http://localhost:1234/v1",
		Model:   "text2cypher-gemma-2-9b-it-finetuned-2024v1", // 使用的模型版本
	}
	cm, err = openai.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}

// newOllamaChatModel component initialization function of node 'ChatModel2' in graph 'graphAgent'
func newOllamaChatModel(ctx context.Context) (cm model.ToolCallingChatModel, err error) {
	config := &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "modelscope.cn/Qwen/Qwen3-14B-GGUF:latest", // 使用的模型版本
	}
	cm, err = ollama.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}

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
