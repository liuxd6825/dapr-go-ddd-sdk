package service

import (
	"context"
	_ "embed"
	"fmt"
	"sync"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/service/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/llm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
)

type RagService struct {
	chatModel model.ToolCallingChatModel
}

var ragServiceOnce sync.Once

//go:embed prompts/import_prompt.md
var importPrompt string

func NewRagService() *RagService {
	var ragService *RagService
	ragServiceOnce.Do(func() {
		ragService = newRagService()
	})
	return ragService
}

func newRagService() *RagService {
	ctx := context.Background()
	e := env.GetEnv()
	ragMeta := e.App.Meta["rag"]
	if ragMeta == nil {
		panic("rag not found in env.app")
	}

	ragCfg, err := config.ReadRagConfig(ragMeta)
	if err != nil {
		panic("read RagConfig error" + err.Error())
	}

	//chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
	//	BaseURL: ragCfg.LLM.BaseUrl,
	//	Model:   ragCfg.LLM.Model, // 使用的模型版本
	//	APIKey:  ragCfg.LLM.APIKey,
	//})

	var llmModel model.ToolCallingChatModel
	switch ragCfg.LLM.Type {
	case "ark":
		m, err := llm.NewArk(ctx, &ark.ChatModelConfig{
			APIKey:  ragCfg.LLM.APIKey,
			BaseURL: ragCfg.LLM.BaseUrl,
			Model:   ragCfg.LLM.Model, // 使用的模型版本
		})

		if err != nil {
			panic("open rag model error" + err.Error())
		}

		llmModel = m
	case "ollama":
		keepalive := time.Duration(-1)
		options := &ollama.Options{}
		if ragCfg.LLM.CtxLength < 4096 {
			options.NumCtx = 4096
		} else {
			options.NumCtx = ragCfg.LLM.CtxLength
		}

		var thinkValue any
		if ragCfg.LLM.Think == "" || ragCfg.LLM.Think == "true" {
			thinkValue = true
		} else if ragCfg.LLM.Think == "high" || ragCfg.LLM.Think == "low" || ragCfg.LLM.Think == "medium" {
			thinkValue = ragCfg.LLM.Think
		} else if ragCfg.LLM.Think != "true" {
			thinkValue = false
		}

		m, err := llm.NewOllama(ctx, ollama.ChatModelConfig{
			BaseURL:   ragCfg.LLM.BaseUrl,
			Model:     ragCfg.LLM.Model, // 使用的模型版本
			KeepAlive: &keepalive,
			Thinking:  &ollama.ThinkValue{Value: thinkValue},
			Options:   options,
		})

		if err != nil {
			panic("open rag model error" + err.Error())
		}

		llmModel = m
	case "openai":
		m, err := llm.NewOpenAI(ctx, openai.ChatModelConfig{
			BaseURL: ragCfg.LLM.BaseUrl,
			Model:   ragCfg.LLM.Model, // 使用的模型版本
			APIKey:  ragCfg.LLM.APIKey,
		})

		if err != nil {
			panic("open rag model error" + err.Error())
		}

		llmModel = m
	default:
		panic("rag config llm.type is null ")
	}

	if err != nil {
		panic("create import.RagModel error" + err.Error())
	}
	ragService := &RagService{}
	ragService.chatModel = llmModel
	return ragService
}

func (s *RagService) Query(ctx context.Context, qry *query.RagQueryRequest, streams ...func(txt string)) (string, error) {
	userPrompt := &schema.Message{
		Role:    schema.User,
		Content: qry.UserPrompt,
	}
	input := []*schema.Message{s.getSystemPrompt(), userPrompt}
	streamReader, err := s.chatModel.Stream(ctx, input)
	if err != nil {
		return "", fmt.Errorf("大模型问题分析时出错: %w", err)
	}

	sb, err := my_rag.Reader(streamReader, streams...)
	if err != nil {
		return "", fmt.Errorf("读取返回结果时出错: %w", err)
	}
	return sb.String(), nil
}

func (s *RagService) getSystemPrompt() *schema.Message {
	return &schema.Message{
		Role:    schema.Assistant,
		Content: importPrompt,
	}
}
