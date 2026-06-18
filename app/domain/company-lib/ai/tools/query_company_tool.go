package tools

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/company-lib/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/company-lib/service"
)

func NewQueryCompanyTool() (tool.InvokableTool, error) {
	// 2. 编写工具的具体执行函数 (入参和出参必须是结构体指针)
	executeFunc := func(ctx context.Context, args *QueryCompanyToolArgs) (*QueryCompanyToolResult, error) {
		fmt.Println("QueryCompanyTool参数：" + args.Text)
		srv := service.NewCompanyService()
		res, err := srv.FullTextSearch(ctx, &model.FullTextSearchQuery{Text: args.Text})
		return &QueryCompanyToolResult{
			CompanyQueryResult: *res,
		}, err
	}

	// 3. 使用 utils.InferTool 自动推导并创建工具
	return utils.InferTool(
		"get_full_text_search_company", // 工具名称
		"通过关键字全文检索公司信息",                // 工具描述
		executeFunc, // 执行函数
	)
}

type QueryCompanyToolArgs struct {
	Text string `json:"text" jsonschema:"description=要检索的内容,required=true"`
}

type QueryCompanyToolResult struct {
	model.CompanyQueryResult
}

// ==========================================
// 1. 定义一个自定义错误，用于通知 Eino 终止运行
// ==========================================
var ErrToolExecutionFinished = errors.New("tool_execution_finished_stop_agent")

// InterceptState 用于在 Context 中跨节点传递工具的执行结果
type InterceptState struct {
	HasExecutedTool bool
	ToolResult      string
}

// ==========================================
// 2. 编写自定义 Callback 处理器
// ==========================================
type AgentBrakeHandler struct {
	// 继承空实现，只需重写我们关心的生命周期方法
	callbacks.Handler
}

// OnAfterStart 在 Eino 节点（此处为 ToolsNode）运行完毕后触发
func (h *AgentBrakeHandler) OnAfterStart(ctx context.Context, info *callbacks.RunInfo, output interface{}, err error) (context.Context, error) {
	// 判断当前运行完毕的节点是不是工具节点 (ToolsNode)
	// 注意：Eino 内部包装的工具节点类型通常为 "ToolsNode" 或对应组件名
	if info.Component == compose.ComponentOfToolsNode || info.Name == "ToolsNode" {

		// 尝试断言输出，获取工具返回的 Message 列表
		if msgs, ok := output.([]*schema.Message); ok && len(msgs) > 0 {
			// 取出最后一个工具的返回结果
			lastMsg := msgs[len(msgs)-1]

			// 从 ctx 中取出我们在外层注入的指针，把结果“偷”出来
			if state, ok := ctx.Value("intercept_key").(*InterceptState); ok {
				state.HasExecutedTool = true
				state.ToolResult = lastMsg.Content

				// 🔥 关键点：返回我们自定义的错误。
				// Eino 的状态机循环一旦在节点执行后收到错误，就会立刻中断，不会再把结果喂给大模型！
				return ctx, ErrToolExecutionFinished
			}
		}
	}
	return ctx, nil
}
