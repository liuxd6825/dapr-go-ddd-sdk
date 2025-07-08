package llm

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

type LLM = model.BaseChatModel

// NewOpenAI component initialization function of node 'ChatModel1' in graph 'graphAgent'
/*
	config := &openai.ChatModelConfig{
		BaseURL: "http://localhost:1234/v1",
		Model:   "text2cypher-gemma-2-9b-it-finetuned-2024v1", // 使用的模型版本
	}
*/
func NewOpenAI(ctx context.Context, cfg openai.ChatModelConfig) (cm model.ToolCallingChatModel, err error) {
	cm, err = openai.NewChatModel(ctx, &cfg)
	return cm, err
}

// NewOllama component initialization function of node 'ChatModel2' in graph 'graphAgent'
/*
	config := &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "modelscope.cn/Qwen/Qwen3-14B-GGUF:latest", // 使用的模型版本
	}
*/
func NewOllama(ctx context.Context, cfg ollama.ChatModelConfig) (cm model.ToolCallingChatModel, err error) {
	cm, err = ollama.NewChatModel(ctx, &cfg)
	return cm, err
}
