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
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/company-lib/ai/tools"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/company-lib/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/llm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
)

type AiCompanyService struct {
	chatModel model.ToolCallingChatModel
	config    *config.RagConfig
	agent     *adk.ChatModelAgent
}

var aiCompanyServiceOnce sync.Once

func NewAiCompanyService() *AiCompanyService {
	var aiCompanyService *AiCompanyService
	aiCompanyServiceOnce.Do(func() {
		aiCompanyService = newAiCompanyService()
		aiCompanyService.init()
		aiCompanyService.initAgent()
	})
	return aiCompanyService
}

func newAiCompanyService() *AiCompanyService {
	e := env.GetEnv()
	ragMeta := e.App.Meta["rag"]
	if ragMeta == nil {
		panic("rag not found in env.app")
	}

	ragCfg, err := config.ReadRagConfig(ragMeta)
	if err != nil {
		panic("read RagConfig error" + err.Error())
	}

	ragService := &AiCompanyService{}
	ragService.config = ragCfg
	return ragService
}

func (s *AiCompanyService) init() {
	ctx := context.Background()
	switch s.config.LLM.Type {
	case "ark":
		m, err := llm.NewArk(ctx, &ark.ChatModelConfig{
			APIKey:  s.config.LLM.APIKey,
			BaseURL: s.config.LLM.BaseUrl,
			Model:   s.config.LLM.Model, // 使用的模型版本
		})

		if err != nil {
			panic("open rag model error" + err.Error())
		}

		s.chatModel = m
	case "ollama":
		keepalive := time.Duration(-1)
		options := &ollama.Options{}
		if s.config.LLM.CtxLength < 4096 {
			options.NumCtx = 4096
		} else {
			options.NumCtx = s.config.LLM.CtxLength
		}

		var thinkValue any
		if s.config.LLM.Think == "" || s.config.LLM.Think == "true" {
			thinkValue = true
		} else if s.config.LLM.Think == "high" || s.config.LLM.Think == "low" || s.config.LLM.Think == "medium" {
			thinkValue = s.config.LLM.Think
		} else if s.config.LLM.Think != "true" {
			thinkValue = false
		}

		m, err := llm.NewOllama(ctx, ollama.ChatModelConfig{
			BaseURL:   s.config.LLM.BaseUrl,
			Model:     s.config.LLM.Model, // 使用的模型版本
			KeepAlive: &keepalive,
			Thinking:  &ollama.ThinkValue{Value: thinkValue},
			Options:   options,
		})

		if err != nil {
			panic("open rag model error" + err.Error())
		}

		s.chatModel = m
	case "openai":
		m, err := llm.NewOpenAI(ctx, openai.ChatModelConfig{
			BaseURL: s.config.LLM.BaseUrl,
			Model:   s.config.LLM.Model, // 使用的模型版本
			APIKey:  s.config.LLM.APIKey,
		})

		if err != nil {
			panic("open rag model error" + err.Error())
		}

		s.chatModel = m
	default:
		panic("rag config llm.type is null ")
	}
}

func (s *AiCompanyService) initAgent() {
	//agent, err := s.BuildAgentGraph()
	//if err != nil {
	//	panic(err)
	//}
	//s.agent = agent
	queryCompanyTool, err := tools.NewQueryCompanyTool()
	if err != nil {
		panic(err)
	}
	agent, err := adk.NewChatModelAgent(context.Background(), &adk.ChatModelAgentConfig{
		Name:        "全文检索公司信息",
		Description: "通过关键字全文检索公司信息",
		Instruction: `你是数据检索助手，必须调用工具进行检索，用户输入带有"搜索"、"查询"、"查找"、"检索"、"全文检索"几个关键词时需要提取具体搜索内容作为工具参数，如果没有关键词直接将提示词作为工具参数。
			# 示例一
			用户输入: "查找时林公司"
			你的思考: 
			- 是否有关键字: "查找"
			- 提取搜索内容: 提取搜索内容为"时林"
			调用工具: 将提取出来的"时林"传入工具参数
			
			# 示例二
			用户输入: "搜索法人邵光震"
			你的思考:
			- 是否有关键字: "搜索"
			- 提取搜索内容: 提取搜索内容为"邵光震"
			调用工具: 将提取出来的"邵光震"传入工具参数

			# 示例三
			用户输入: "全文检索地址知春路56号"
			你的思考:
			- 是否有关键字: "全文检索"
			- 提取搜索内容: 提取搜索内容为"知春路56号"
			调用工具: 将提取出来的"知春路56号"传入工具参数

			# 示例四
			用户输入: "海淀区"
			你的思考:
			- 是否有关键字: 无关键字
			- 提取搜索内容: 提取搜索内容为"海淀区"
			调用工具: 将提取出来的"海淀区"传入工具参数
					`,
		Model: s.chatModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{queryCompanyTool},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	s.agent = agent
}

func (s *AiCompanyService) Query(ctx context.Context, qry *model2.AiQueryRequest, streams ...func(txt string)) (string, error) {
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           s.agent,
		EnableStreaming: true,
	})

	userInput := qry.Prompt
	fmt.Println("用户提示词：" + userInput)

	eventIter := runner.Query(ctx, userInput)

	var result string
	for {
		event, ok := eventIter.Next()
		if !ok {
			break
		}

		if event.Err != nil {
			return result, event.Err
		}

		if event.Action != nil {
			if event.Action.Exit {
				break
			}
			if event.Action.Interrupted != nil {
				return result, fmt.Errorf("agent interrupted: %v", event.Action.Interrupted.Data)
			}
		}

		if event.Output != nil && event.Output.MessageOutput != nil {
			msgOutput := event.Output.MessageOutput
			var stream func(txt string)
			if len(streams) > 0 {
				stream = streams[len(streams)-1]
			}
			if msgOutput.Role == schema.Tool {
				result += msgOutput.Message.Content
				break
			}
			if msgOutput.Message != nil {
				if len(msgOutput.Message.ToolCalls) > 0 {
					continue
				}
				//result += msgOutput.Message.Content
				stream(msgOutput.Message.Content)
				break
			} else if msgOutput.MessageStream != nil {
				defer msgOutput.MessageStream.Close()
				for {
					msg, err := msgOutput.MessageStream.Recv()
					if err != nil {
						break
					}
					//result += msg.Content
					stream(msg.Content)
				}
			}
		}
	}

	return result, nil
}
