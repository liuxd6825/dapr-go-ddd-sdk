package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/tasks"
)

// RegisterWorkflow
// @Description: 注册工作流
func RegisterWorkflow() {
	// 分析工作流
	ctx := context.Background()
	tasks.RegisterWorkflow(ctx, TaskAnalysisWorkflow)
	tasks.RegisterActivity(ctx, TaskAnalysisActivity)
}
